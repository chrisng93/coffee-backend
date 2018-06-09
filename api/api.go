package api

import (
	"github.com/gorilla/mux"
)

func Init() *mux.Router {
	r := mux.NewRouter()
	return r
}
