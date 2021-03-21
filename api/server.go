package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/inadislam/go-login-api/database"
	"github.com/inadislam/go-login-api/routes"
)

func StartServer() {
	message := `
				Server starting...
				Server started successfully.
				Database connection eastablished.
				Creating Table in Database
				Please visit http://localhost:8080/home
			`
	database.Connect()
	database.AutoMigrator()
	r := routes.NewRouter()
	fmt.Println(message)
	log.Fatal(http.ListenAndServe(":8080", r))
}
