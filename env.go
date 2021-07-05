package ditt

import (
	"github.com/gorilla/sessions"
)

var Env = struct {
	DataStore     UserDataStore
	Files         Files
	AdminPassword string
	CookiesStore  sessions.Store
}{
	DataStore:    NewUserDataMemoryStore(),
	Files:        NewMemoryFiles(),
	CookiesStore: sessions.NewCookieStore(),
}
