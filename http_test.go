package ditt

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	_httpTestsCookies []*http.Cookie
)

func setupHttpTests() {
	Env.AdminPassword = "password"
	Env.CookiesStore = sessions.NewCookieStore()
}

type _httpTestUser struct {
	Id       string `json:"id"`
	Password string `json:"password"`
	Data     string `json:"data"`
}

var (
	_httpTestAddUsersBodyContent = `
[
{
	"id": "user-1",	
	"password": "pass-1",
	"data": "data-1"
},
{
	"id": "user-2",	
	"password": "pass-2",
	"data": "data-2"
},
{
	"id": "user-3",	
	"password": "pass-3",
	"data": "data-3"
},
{
	"id": "user-4",	
	"password": "pass-4",
	"data": "data-4"
},
{
	"id": "user-5",	
	"password": "pass-5",
	"data": "data-5"
}
]
`
)

func TestHandleHttpLoginRequest(t *testing.T) {
	Convey("Login", t, func() {
		setupHttpTests()

		bodyContent := `{"login": "admin", "password": "password"}`
		r := httptest.NewRequest(http.MethodPost, LoginEndpoint, bytes.NewBufferString(bodyContent))

		w := httptest.NewRecorder()
		HandleHttpAddUsersRequest(w, r)

		res := w.Result()
		defer func() {
			_ = res.Body.Close()
		}()

		So(res.StatusCode, ShouldEqual, http.StatusOK)
		_httpTestsCookies = res.Cookies()
	})
}

func TestHandleHttpAddUsersRequest(t *testing.T) {
	Convey("Add Users", t, func() {
		setupHttpTests()

		r := httptest.NewRequest(http.MethodPost, AddUsersEndpoint, bytes.NewBufferString(_httpTestAddUsersBodyContent))
		for _, cookie := range _httpTestsCookies {
			r.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		HandleHttpAddUsersRequest(w, r)

		res := w.Result()
		defer func() {
			_ = res.Body.Close()
		}()
		So(res.StatusCode, ShouldEqual, http.StatusOK)
	})
}

func TestHandleHttpGetUserRequest(t *testing.T) {
	Convey("Get User", t, func() {
		setupHttpTests()

		endpoint := strings.Replace(GetUserEndpoint, endpointVarId, "user-1", 1)

		r := httptest.NewRequest(http.MethodGet, endpoint, nil)
		for _, cookie := range _httpTestsCookies {
			r.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		HandleHttpAddUsersRequest(w, r)

		res := w.Result()
		defer func() {
			_ = res.Body.Close()
		}()
		So(res.StatusCode, ShouldEqual, http.StatusOK)

		var user *_httpTestUser
		err := json.NewDecoder(r.Body).Decode(&user)
		So(err, ShouldBeNil)

		So(user.Id, ShouldEqual, "user-1")
		So(user.Password, ShouldEqual, "password-1")
		So(user.Data, ShouldEqual, "data-1")
	})
}

func TestHandleHttpGetUserListRequest(t *testing.T) {
	Convey("List Users", t, func() {
		setupHttpTests()

		r := httptest.NewRequest(http.MethodGet, ListUsersEndpoint, nil)
		for _, cookie := range _httpTestsCookies {
			r.AddCookie(cookie)
		}

		w := httptest.NewRecorder()
		HandleHttpAddUsersRequest(w, r)

		res := w.Result()
		defer func() {
			_ = res.Body.Close()
		}()

		So(res.StatusCode, ShouldEqual, http.StatusOK)

	})
}

func TestHandleHttpUpdateUserRequest(t *testing.T) {
	Convey("Update User", t, func() {
		setupHttpTests()

		endpoint := strings.Replace(UpdateUserEndpoint, endpointVarId, "user-1", 1)

		r := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(_httpTestAddUsersBodyContent))
		w := httptest.NewRecorder()
		HandleHttpAddUsersRequest(w, r)

		res := w.Result()
		defer func() {
			_ = res.Body.Close()
		}()
		So(res.StatusCode, ShouldEqual, http.StatusOK)

		endpoint = strings.Replace(GetUserEndpoint, endpointVarId, "user-2", 1)

		r = httptest.NewRequest(http.MethodGet, endpoint, nil)
		w = httptest.NewRecorder()
		HandleHttpAddUsersRequest(w, r)

		res2 := w.Result()
		defer func() {
			_ = res2.Body.Close()
		}()
		So(res2.StatusCode, ShouldEqual, http.StatusNotFound)
	})
}

func TestHandleHttpDeleteUserRequest(t *testing.T) {
	Convey("", t, func() {

		bodyContent := `{
			"id": "user-2",	
			"password": "pass-2-updated",
			"data": "data-2-updated"}`

		endpoint := strings.Replace(DeleteUserEndpoint, endpointVarId, "user-2", 1)

		r := httptest.NewRequest(http.MethodPatch, endpoint, bytes.NewBufferString(bodyContent))
		w := httptest.NewRecorder()
		HandleHttpAddUsersRequest(w, r)

		res := w.Result()
		defer func() {
			_ = res.Body.Close()
		}()
		So(res.StatusCode, ShouldEqual, http.StatusOK)

		// Loading update user data

		endpoint = strings.Replace(GetUserEndpoint, endpointVarId, "user-2", 1)

		r = httptest.NewRequest(http.MethodGet, endpoint, nil)
		w = httptest.NewRecorder()
		HandleHttpAddUsersRequest(w, r)

		res2 := w.Result()
		defer func() {
			_ = res2.Body.Close()
		}()
		So(res2.StatusCode, ShouldEqual, http.StatusOK)

		var user *_httpTestUser
		err := json.NewDecoder(r.Body).Decode(&user)
		So(err, ShouldBeNil)

		So(user.Id, ShouldEqual, "user-2")
		So(user.Password, ShouldEqual, "pass-2-updated")
		So(user.Data, ShouldEqual, "data-2-updated")
	})
}
