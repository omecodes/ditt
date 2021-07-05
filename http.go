package ditt

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	endpointVarId    = "id"
	queryParamOffset = "offset"
	queryParamCount  = "count"

	// LoginEndpoint is the HTTP API endpoint to initialise an authenticated session
	LoginEndpoint = "/login"

	// AddUsersEndpoint is the HTTP API endpoint to add users
	AddUsersEndpoint = "/add/users"

	// DeleteUserEndpoint is the HTTP API endpoint to delete data of a user
	DeleteUserEndpoint = "/delete/user/{id}"

	// GetUserEndpoint is the HTTP API endpoint to get user data
	GetUserEndpoint = "/user/{id}"

	// ListUsersEndpoint is the HTTP API endpoint to add users
	ListUsersEndpoint = "/users/list"

	// UpdateUserEndpoint is the HTTP API endpoint to add users
	UpdateUserEndpoint = "/user/{id}"
)

// HandleHttpLoginRequest initializes an APIHandler and calls its APIHandler.Login method
// with user credentials parsed from the request body
func HandleHttpLoginRequest(w http.ResponseWriter, r *http.Request) {

	var credentials = &struct {
		Login    string
		Password string
	}{}
	api := NewAPIHandler()

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			log.Println("credentials parsing:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

	case "x-www-form-urlencoded":
		err := r.ParseForm()
		if err != nil {
			log.Println("credentials parsing:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		credentials.Login = r.Form.Get("login")
		credentials.Password = r.Form.Get("password")

	case "multipart/form-data":
		err := r.ParseMultipartForm(-1)
		if err != nil {
			log.Println("credentials parsing:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		credentials.Login = r.MultipartForm.Value["login"][0]
		credentials.Password = r.MultipartForm.Value["password"][0]
	}

	ok, err := api.Login(r.Context(), credentials.Login, credentials.Password)
	if err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	session, _ := Env.CookiesStore.Get(r, sessionName)
	session.Values[sessionLoggedUserKey] = credentials.Login
	err = session.Save(r, w)
	if err != nil {
		log.Println("session saving:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// HandleHttpAddUsersRequest initializes an APIHandler and calls its APIHandler.AddUsers with the request body content
// The request body content is expected to be an encoded JSON object list
func HandleHttpAddUsersRequest(w http.ResponseWriter, r *http.Request) {

	var (
		err     error
		content io.Reader
	)
	api := NewAPIHandler()

	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "multipart/form-data":
		err := r.ParseMultipartForm(-1)
		if err != nil {
			log.Println("credentials parsing:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fh := r.MultipartForm.File["file"][0]
		file, err := fh.Open()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer func() {
			_ = file.Close()
		}()
		content = file
	default:
		content = r.Body
	}

	err = api.AddUsers(r.Context(), content)
	if err != nil {
		w.WriteHeader(statusFromError(err))
	}
}

// HandleHttpDeleteUserRequest initializes an APIHandler and calls its APIHandler.DeleteUser with userId extracted from the request URI path
func HandleHttpDeleteUserRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars[endpointVarId]

	api := NewAPIHandler()
	err := api.DeleteUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(statusFromError(err))
	}
}

// HandleHttpGetUserRequest initializes an APIHandler and calls its APIHandler.GetUser with userId extracted from the request URI path
// The returned value by GetUser is set as the HTTP response body
func HandleHttpGetUserRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars[endpointVarId]

	api := NewAPIHandler()
	user, err := api.GetUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write([]byte(user))
}

// HandleHttpGetUserListRequest initializes an APIHandler and calls its APIHandler.GetUserList
// "r" is expected to be a POST request with an encoded JSON object list as body
// The returned value by GetUserList is set as the HTTP response body
func HandleHttpGetUserListRequest(w http.ResponseWriter, r *http.Request) {
	var (
		offset, count int
		err           error
	)

	query := r.URL.Query()
	offsetValue := query.Get(queryParamOffset)
	countValue := query.Get(queryParamCount)

	if offsetValue != "" {
		offset, err = strconv.Atoi(offsetValue)
		if err != nil {
			writeHttpErrorResponseWithMessage(w, BadInput, "expected a number as value of 'offset'")
			return
		}
	}
	if countValue != "" {
		count, err = strconv.Atoi(countValue)
		if err != nil {
			writeHttpErrorResponseWithMessage(w, BadInput, "expected a number as value of 'count'")
			return
		}
	}

	api := NewAPIHandler()
	list, err := api.GetUserList(r.Context(), ListOptions{
		Offset: offset,
		Count:  count,
	})
	if err != nil {
		w.WriteHeader(statusFromError(err))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf("{\"offset\": %d,", list.Offset)))
	_, _ = w.Write([]byte("\"data\": ["))
	for ind, userData := range list.UserDataList {
		if ind > 0 {
			_, _ = w.Write([]byte(","))
		}
		_, _ = w.Write([]byte(userData))
	}
	_, _ = w.Write([]byte("]"))
	_, _ = w.Write([]byte("}"))
}

// HandleHttpUpdateUserRequest initializes an APIHandler and calls its APIHandler.UpdateUser with userId extracted from the request URI path
// and the request content body. The request body content is expected to be a JSON encoded UserData object
func HandleHttpUpdateUserRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars[endpointVarId]

	api := NewAPIHandler()
	err := api.DeleteUser(r.Context(), userId)
	if err != nil {
		w.WriteHeader(statusFromError(err))
	}
}
