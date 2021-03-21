package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	ToJson(w, statusCode, struct {
		Error  error `json:"error"`
		Status int   `json:"status"`
	}{
		Error:  err,
		Status: statusCode,
	})
}

func ToJson(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-type:", "application/json, charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	ToJson(w, 200, struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "It is not implemented by DEV!!",
		Status:  200,
	})
}
