package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/inadislam/go-login-api/controllers"
	"github.com/inadislam/go-login-api/middlewares"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/home", middlewares.BasicMiddleware(http.HandlerFunc(controllers.Home))).Methods("GET")
	r.HandleFunc("/login", middlewares.BasicMiddleware(http.HandlerFunc(controllers.Login))).Methods("POST")
	r.HandleFunc("/signup", middlewares.BasicMiddleware(http.HandlerFunc(controllers.Register))).Methods("POST")
	r.HandleFunc("/users", middlewares.BasicMiddleware(middlewares.IsAuth(http.HandlerFunc(controllers.GetUsers)))).Methods("POST")
	r.HandleFunc("/user/userid={id}", middlewares.BasicMiddleware(middlewares.IsAuth(http.HandlerFunc(controllers.GetUserById)))).Methods("POST")
	r.HandleFunc("/profile/delid={id}", middlewares.BasicMiddleware(middlewares.IsAuth(http.HandlerFunc(controllers.DeleteProfile)))).Methods("POST")
	r.HandleFunc("/profile/updid={id}", middlewares.BasicMiddleware(middlewares.IsAuth(http.HandlerFunc(controllers.UpdateUser)))).Methods("POST")

	r.HandleFunc("/home", controllers.Cors).Methods("OPTIONS")
	r.HandleFunc("/user-login", controllers.Cors).Methods("OPTIONS")
	r.HandleFunc("/user-signup", controllers.Cors).Methods("OPTIONS")
	r.HandleFunc("/users", controllers.Cors).Methods("OPTIONS")
	r.HandleFunc("/user/userid={id}", controllers.Cors).Methods("OPTIONS")
	r.HandleFunc("/profile/delid={id}", controllers.Cors).Methods("OPTIONS")
	r.HandleFunc("/profile/updid={id}", controllers.Cors).Methods("OPTIONS")
	return r
}
