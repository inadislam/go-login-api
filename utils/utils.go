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
			Status string `json:"status"`
		}{
			Error:  err.Error(),
			Status: http.StatusText(statusCode),
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
