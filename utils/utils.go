package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		ToJson(w, statusCode, struct {
			Error  string `json:"error"`
			Status int    `json:"status"`
		}{
			Error:  err.Error(),
			Status: statusCode,
		})
	}
	ToJson(w, http.StatusBadRequest, nil)
}

func ToJson(w http.ResponseWriter, status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-type:", "application/json, charset=utf-8")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
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
