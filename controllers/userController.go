package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/inadislam/go-login-api/database"
	"github.com/inadislam/go-login-api/models"
	"github.com/inadislam/go-login-api/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		log.Fatal(err)
		return
	}
	userCreated, err := database.SignupHelper(user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	utils.ToJson(w, http.StatusCreated, userCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		log.Fatal(err)
		return
	}
	token, err := database.LoginHelper(user.Username, user.Password)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		log.Fatal(err)
		return
	}

	utils.ToJson(w, http.StatusOK, struct {
		Token  string `json:"token"`
		Status int    `json:"status"`
	}{
		Token:  token,
		Status: http.StatusOK,
	})
}
