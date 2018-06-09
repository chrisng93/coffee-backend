package api

import (
	"github.com/gorilla/mux"
)

// Init initializes the router instance and sets handlers for the API.
func Init() *mux.Router {
	r := mux.NewRouter()
	// TODO: Set handlers for endpoints.
	return r
}
