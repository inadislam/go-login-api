package routes

import (
	"github.com/gorilla/mux"
	"github.com/inadislam/go-login-api/controllers"
	"github.com/inadislam/go-login-api/utils"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/home", utils.NotImplemented).Methods("GET")
	r.HandleFunc("/user-login", controllers.Login).Methods("POST")
	r.HandleFunc("/user-signup", controllers.Register).Methods("POST")
	return r
}
