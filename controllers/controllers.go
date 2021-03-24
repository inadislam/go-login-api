package controllers

import (
	"net/http"

	"github.com/inadislam/go-login-api/utils"
)

func Cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "application/json")
}

func Home(w http.ResponseWriter, r *http.Request) {
	utils.ToJson(w, http.StatusOK, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Welcome to Home Dev!!",
		Status:  http.StatusText(http.StatusOK),
	})
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	utils.ToJson(w, http.StatusOK, struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "It is not implemented by DEV!!",
		Status:  http.StatusText(http.StatusOK),
	})
}
