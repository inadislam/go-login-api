package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	utils.CheckErr(err)
	user, err := database.UserById(uint32(key))
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	utils.ToJson(w, http.StatusOK, struct {
		User   models.User `json:"userInfo"`
		Status int         `json:"status"`
	}{
		User:   user,
		Status: http.StatusOK,
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	user, err := database.GetAllUser()
	utils.CheckErr(err)
	utils.ToJson(w, http.StatusCreated, struct {
		User   models.User `json:"users"`
		Status int         `json:"status"`
	}{
		User:   user,
		Status: http.StatusCreated,
	})
}

func DeleteProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, err = database.DeleteUser(uint32(uid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte("User Deleted"))
}
