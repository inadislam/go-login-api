package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inadislam/go-login-api/controllers"
	"github.com/inadislam/go-login-api/middlewares"
	"github.com/inadislam/go-login-api/utils"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/home", utils.NotImplemented).Methods("GET")

	// User Handler
	r.HandleFunc("/user-login", controllers.Login).Methods("POST")
	r.HandleFunc("/user-signup", controllers.Register).Methods("POST")
	r.HandleFunc("/users", middlewares.JwtMiddlewares(http.HandlerFunc(controllers.GetAllUser)))
	r.HandleFunc("/get-user/userid={id}", middlewares.JwtMiddlewares(http.HandlerFunc(controllers.GetUserById)))
	r.HandleFunc("/user-profile/delid={id}", middlewares.JwtMiddlewares(http.HandlerFunc(controllers.GetUserById)))
	r.HandleFunc("/user-update/updateid={id}", middlewares.JwtMiddlewares(http.HandlerFunc(controllers.GetUserById)))
	return r
}
