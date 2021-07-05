package ditt

import (
	"sync"

	"github.com/tidwall/sjson"
	"golang.org/x/crypto/bcrypt"
)

// UserDataCallback is function that handle a UserData
type UserDataCallback func(UserData) error

type UserDataProvider func(callback UserDataCallback) error

// UserDataProcessor is a convenience for user data processor
type UserDataProcessor interface {
	ProcessData(data UserData) (UserData, error)
}

// UserDataProcessorFunc is a function that implements  of UserDataProcessor
type UserDataProcessorFunc func(data UserData) (UserData, error)

func (f UserDataProcessorFunc) ProcessData(data UserData) (UserData, error) {
	return f(data)
}

// UserDataProcessingResult holds info about UserData processing result
type UserDataProcessingResult struct {
	Err    error
	UserId string
	Data   UserData
}

// ConcurrentUserDataProcessingRunner starts a routine. And for each provided UserData calls the callback and returns
// the result to the caller through the "results" channel
type ConcurrentUserDataProcessingRunner struct {
	Provider            UserDataProvider
	Processor           UserDataProcessor
	TasksResultsSignals chan<- chan UserDataProcessingResult
	ResultSignal        chan<- chan error
}

func (r ConcurrentUserDataProcessingRunner) dispatch() error {
	results := make(chan UserDataProcessingResult, 1)
	defer close(results)

	wg := &sync.WaitGroup{}
	r.TasksResultsSignals <- results

	err := r.Provider(func(data UserData) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := UserDataProcessingResult{UserId: data.Id()}
			result.Data, result.Err = r.Processor.ProcessData(data)
			results <- result
		}()
		return nil
	})
	wg.Wait()
	return err
}

func (r ConcurrentUserDataProcessingRunner) Run() {
	go func() {
		err := r.dispatch()
		result := make(chan error)
		defer close(result)
		r.ResultSignal <- result
		result <- err
	}()
}

var writeProcessors = []UserDataProcessor{
	UserDataProcessorFunc(saveDataIntoFile),
	UserDataProcessorFunc(hashPassword),
	// UserDataProcessorFunc(transformUserId),
}

var readProcessors = []UserDataProcessor{
	// UserDataProcessorFunc(removeUserId),
	UserDataProcessorFunc(mergeWithDataFromFile),
}

func processData(processors []UserDataProcessor, data UserData) (UserData, error) {
	var err error
	for _, processor := range processors {
		data, err = processor.ProcessData(data)
		if err != nil {
			return "", err
		}
	}
	return data, nil
}

func hashPassword(data UserData) (UserData, error) {
	if data == "" {
		return data, nil
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(data.Password()), 12)
	if err != nil {
		return "", err
	}
	updateData, err := sjson.Set(string(data), "password", string(hashBytes))
	return UserData(updateData), err
}

func saveDataIntoFile(data UserData) (UserData, error) {
	if data == "" {
		return data, nil
	}

	err := Env.Files.Save(data.Id(), data.Data())
	if err != nil {
		return "", err
	}

	updateData, err := sjson.Delete(string(data), "data")
	return UserData(updateData), err
}

/* func transformUserId(data UserData) (UserData, error) {
	if data == "" {
		return data, nil
	}

	id := data.Id()
	updateData, err := sjson.Set(string(data), "userId", id)
	return UserData(updateData), err
}*/

func mergeWithDataFromFile(data UserData) (UserData, error) {
	if data == "" {
		return data, nil
	}

	content, err := Env.Files.Get(data.Id())
	if err != nil {
		return data, err
	}

	if content != "" {
		updateData, err := sjson.Set(string(data), "data", content)
		return UserData(updateData), err
	}
	return data, nil
}

/* func removeUserId(data UserData) (UserData, error) {

	if data == "" {
		return data, nil
	}

	updateData, err := sjson.Delete(string(data), "userId")
	return UserData(updateData), err
} */
