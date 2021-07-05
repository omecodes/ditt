package ditt

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"log"
	"sync"
)

// UserDataStore is a convenience for UserData persistence management
type UserDataStore interface {

	// Save saves user data
	Save(data UserData) error

	// Delete deletes the userData matching the given id
	Delete(id string) error

	// Get retrieves userData matching the given id
	Get(id string) (UserData, error)

	// ListForUser fetches a range of UserData that matches the userId
	// and pass each pass each parsed userdata to the callback
	ListForUser(userId string, offset, count int, callback UserDataCallback) error

	// List fetches a range of UserData and pass each pass each
	// parsed userdata to the callback
	List(offset, count int, callback UserDataCallback) error
}

type memoryDataStore struct {
	sync.Mutex
	records map[string]UserData
}

func (m *memoryDataStore) Save(data UserData) error {
	m.Lock()
	defer m.Unlock()
	m.records[data.Id()] = data
	return nil
}

func (m *memoryDataStore) Delete(id string) error {
	m.Lock()
	defer m.Unlock()
	_, found := m.records[id]
	if !found {
		return NotFound
	}
	delete(m.records, id)
	return nil
}

func (m *memoryDataStore) Get(id string) (UserData, error) {
	m.Lock()
	defer m.Unlock()
	data, found := m.records[id]
	if !found {
		return "", NotFound
	}
	return data, nil
}

func (m *memoryDataStore) ListForUser(userId string, offset, count int, callback UserDataCallback) error {
	m.Lock()
	defer m.Unlock()

	loadedCount := 0

	for id, data := range m.records {
		if userId == id {
			if offset == 0 {
				err := callback(data)
				if err != nil {
					return err
				}
				loadedCount++
				if loadedCount == count {
					return nil
				}
			} else {
				offset--
			}
		}
	}

	return nil
}

func (m *memoryDataStore) List(offset, count int, callback UserDataCallback) error {
	m.Lock()
	defer m.Unlock()
	loadedCount := 0

	for _, data := range m.records {
		if offset == 0 {
			err := callback(data)
			if err != nil {
				return err
			}
			loadedCount++
			if loadedCount == count {
				return nil
			}
		} else {
			offset--
		}
	}
	return nil
}

// NewUserDataMemoryStore constructs a memory based UserDataStore
func NewUserDataMemoryStore() UserDataStore {
	return &memoryDataStore{records: make(map[string]UserData)}
}

const (
	databaseName   = "ditt"
	collectionName = "users"
)

type mongoDataStore struct {
	usersCollection *mgo.Collection
	db              *mgo.Database
	session         *mgo.Session
}

func (m *mongoDataStore) Save(data UserData) error {
	var doc interface{}
	err := bson.UnmarshalJSON([]byte(data), &doc)
	if err != nil {
		return BadInput
	}

	err = m.usersCollection.Insert(doc)
	if err != nil {
		log.Println("mongo save:", err)
		return Internal
	}
	return nil
}

func (m *mongoDataStore) Delete(id string) error {
	err := m.usersCollection.Remove(bson.M{"id": id})
	if err != nil {
		if err == mgo.ErrNotFound {
			return NotFound
		}
		return Internal
	}
	return nil
}

func (m *mongoDataStore) Get(id string) (UserData, error) {
	var result interface{}
	err := m.usersCollection.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			return "", NotFound
		}
		return "", Internal
	}

	data, err := bson.MarshalJSON(result)
	return UserData(data), err
}

func (m *mongoDataStore) ListForUser(userId string, offset, count int, callback UserDataCallback) error {
	iter := m.usersCollection.Find(bson.M{"id": userId}).Limit(count).Skip(offset).Iter()
	defer func() {
		_ = iter.Close()
	}()

	var result interface{}
	for iter.Next(&result) {
		data, err := bson.MarshalJSON(result)
		if err != nil {
			return err
		}
		userData := UserData(data)
		err = callback(userData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mongoDataStore) List(offset, count int, callback UserDataCallback) error {
	iter := m.usersCollection.Find(bson.M{}).Limit(count).Skip(offset).Iter()
	defer func() {
		_ = iter.Close()
	}()

	var result interface{}
	for iter.Next(&result) {
		data, err := bson.MarshalJSON(result)
		if err != nil {
			return err
		}

		userData := UserData(data)
		err = callback(userData)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMongoUserDataStore(uri string) (UserDataStore, error) {
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}

	db := session.DB(databaseName)
	col := db.C(collectionName)
	err = col.EnsureIndex(mgo.Index{
		Key:      []string{"id"},
		Unique:   true,
		DropDups: true,
	})
	if err != nil && !mgo.IsDup(err) {
		return nil, err
	}

	return &mongoDataStore{
		session:         session,
		usersCollection: col,
		db:              db,
	}, nil
}
