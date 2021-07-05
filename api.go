package ditt

import (
	"context"
	"io"

	"github.com/tidwall/gjson"
)

const (
	DefaultUserListCount = 1
)

// UserData wraps string for readability
type UserData string

func (d UserData) Id() string {
	result := gjson.Get(string(d), "id")
	return result.String()
}

// Data retrieves the value of the "data" field
func (d UserData) Data() string {
	result := gjson.Get(string(d), "data")
	return result.String()
}

// Password retrieves the value of the "password" field
func (d UserData) Password() string {
	result := gjson.Get(string(d), "password")
	return result.String()
}

// UserDataList holds info about a range of UserData objects
type UserDataList struct {
	Offset       int        `json:"offset"`
	UserDataList []UserData `json:"user_data_list"`
}

// ListOptions defines a range limit
type ListOptions struct {
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

type APIHandler interface {

	// Login creates an authenticated session if credentials match a registered user
	Login(ctx context.Context, username string, password string) (bool, error)

	// AddUsers parses user data list from the "reader" stream and store them in the database
	AddUsers(ctx context.Context, reader io.Reader) error

	// DeleteUser deletes the user identified by "userId' from the database
	DeleteUser(ctx context.Context, userId string) error

	// GetUser retrieves the data of the user identified by "userId"
	GetUser(ctx context.Context, userId string) (UserData, error)

	// GetUserList retrieves a set of users from the database
	GetUserList(ctx context.Context, opts ListOptions) (*UserDataList, error)

	// UpdateUser updates the data of the user identified by "userId"
	UpdateUser(ctx context.Context, userId string, userData UserData) error
}
