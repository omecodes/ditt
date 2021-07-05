package ditt

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	BadInput               = errors.New("bad input")
	AuthenticationRequired = errors.New("authentication required")
	Forbidden              = errors.New("forbidden")
	NotAuthorized          = errors.New("not authorized")
	NotFound               = errors.New("not found")
	Internal               = errors.New("internal")
)

func statusFromError(err error) int {
	switch err {
	case BadInput:
		return http.StatusBadRequest
	case AuthenticationRequired, Forbidden:
		return http.StatusForbidden
	case NotAuthorized:
		return http.StatusUnauthorized
	case NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

/*func writeHttpDataResponse(w http.ResponseWriter, reader io.Reader) {
	w.Header().Add("Content-Type", "application/json")
	_, _ = io.Copy(w, reader)
} */

func writeHttpObjectResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

/*func writeHttpStringResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
} */

func writeHttpErrorResponseWithMessage(w http.ResponseWriter, err error, message string) {
	w.WriteHeader(statusFromError(err))
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
	})
}
