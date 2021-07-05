package ditt

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Config struct {
	Port      int `json:"port"`
	TlsConfig *tls.Config
}

func Serve(config *Config) error {
	var handler http.Handler
	router := mux.NewRouter()

	router.Name("Login").Path(LoginEndpoint).Methods(http.MethodPost).HandlerFunc(HandleHttpLoginRequest)
	router.Name("Create").Path(AddUsersEndpoint).Methods(http.MethodPost).HandlerFunc(HandleHttpAddUsersRequest)
	router.Name("Delete").Path(DeleteUserEndpoint).Methods(http.MethodDelete).HandlerFunc(HandleHttpDeleteUserRequest)
	router.Name("Read").Path(GetUserEndpoint).Methods(http.MethodGet).HandlerFunc(HandleHttpGetUserRequest)
	router.Name("List").Path(ListUsersEndpoint).Methods(http.MethodGet).HandlerFunc(HandleHttpGetUserListRequest)
	router.Name("Update").Path(UpdateUserEndpoint).Methods(http.MethodPatch).HandlerFunc(HandleHttpUpdateUserRequest)

	handler = router
	handler = sessionHttpMiddleware(handler)
	handler = loggerHttpMiddleware(handler)

	srv := http.Server{
		Addr:      fmt.Sprintf(":%d", config.Port),
		Handler:   handler,
		TLSConfig: config.TlsConfig,
	}
	defer func() {
		_ = srv.Shutdown(context.Background())
	}()

	log.Println("Listen ", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
	return err
}
