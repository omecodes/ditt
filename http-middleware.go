package ditt

import (
	"log"
	"net/http"
	"time"
)

type statusCatcher struct {
	status int
	http.ResponseWriter
}

func (catcher *statusCatcher) WriteHeader(statusCode int) {
	catcher.status = statusCode
	catcher.ResponseWriter.WriteHeader(statusCode)
}

func loggerHttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		c := &statusCatcher{
			status:         0,
			ResponseWriter: w,
		}

		next.ServeHTTP(c, r)

		duration := time.Since(start)
		log.Printf("%s %s %s %s\n", r.Method, r.RequestURI, http.StatusText(c.status), duration)
	})
}

const (
	sessionName          = "auth-session"
	sessionLoggedUserKey = "logged-user"
)

func sessionHttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Env.CookiesStore.Get(r, sessionName)
		value, exists := session.Values[sessionLoggedUserKey]
		if exists {
			r = r.WithContext(ContextWithLoggedUser(r.Context(), value.(string)))
		}
		next.ServeHTTP(w, r)
	})
}
