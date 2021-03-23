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
	r.HandleFunc("/home", middlewares.BasicMiddleware(http.HandlerFunc(utils.NotImplemented))).Methods("GET")

	// User Handler
	r.HandleFunc("/user-login", middlewares.BasicMiddleware(http.HandlerFunc(controllers.Login))).Methods("POST")
	r.HandleFunc("/user-signup", middlewares.BasicMiddleware(controllers.Register)).Methods("POST")
	r.HandleFunc("/users", middlewares.BasicMiddleware(middlewares.JwtMiddleware(controllers.GetUsers))).Methods("POST")
	r.HandleFunc("/user/userid={id}", middlewares.BasicMiddleware(middlewares.JwtMiddleware(controllers.GetUserById))).Methods("POST")
	r.HandleFunc("/profile/delid={id}", middlewares.BasicMiddleware(middlewares.JwtMiddleware(controllers.DeleteProfile))).Methods("POST")
	r.HandleFunc("/profile/updid={id}", middlewares.BasicMiddleware(middlewares.JwtMiddleware(utils.NotImplemented))).Methods("POST")
	return r
}
