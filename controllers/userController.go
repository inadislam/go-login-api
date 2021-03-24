package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/badoux/checkmail"
	"github.com/inadislam/go-login-api/auth"
	"github.com/inadislam/go-login-api/database"
	"github.com/inadislam/go-login-api/models"
	"github.com/inadislam/go-login-api/utils"

	"github.com/gorilla/mux"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("request unprocessable"))
		return
	}
	var users models.User
	err = json.Unmarshal(body, &users)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("no json data found in request body"))
		return
	}
	if users.Name == "" || users.Email == "" || users.Username == "" || users.Password == "" {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("fields cannot be empty"))
		return
	}

	if err = checkmail.ValidateFormat(users.Email); err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("invalid email format"))
		return
	}
	userCreated, err := database.SignupHelper(users)
	if err == nil {
		utils.ToJson(w, http.StatusCreated, struct {
			Name     string `json:"name"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Name:     userCreated.Name,
			Username: userCreated.Username,
			Email:    userCreated.Email,
			Password: "Your Password",
		})
	} else {
		utils.ERROR(w, http.StatusConflict, errors.New("username or email aleready exists"))
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("request unprocessable"))
		return
	}
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("no json data found in request body"))
		return
	}
	if user.Username == "" || user.Password == "" {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("fields can not be empty"))
		return
	}
	token, ud, err := database.LoginHelper(user.Username, user.Password)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	} else {

		utils.ToJson(w, http.StatusAccepted, struct {
			Id       uint32 `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Token    string `json:"token"`
		}{
			Id:       ud.ID,
			Name:     ud.Name,
			Username: ud.Username,
			Email:    ud.Email,
			Token:    token,
		})
	}
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
		Status string      `json:"status"`
	}{
		User:   user,
		Status: http.StatusText(http.StatusOK),
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	user, err := database.GetAllUser()
	utils.CheckErr(err)
	utils.ToJson(w, http.StatusCreated, struct {
		User   []models.User `json:"users"`
		Status string        `json:"status"`
	}{
		User:   user,
		Status: http.StatusText(http.StatusCreated),
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

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("request unprocessable"))
		return
	}
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("no json data found in request body"))
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("user unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	if err = checkmail.ValidateFormat(user.Email); err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("invalid email format"))
		return
	}

	userUpdate, err := database.UserUpdate(user, tokenID)
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New("internal server error"))
	}

	utils.ToJson(w, http.StatusAccepted, struct {
		Id       uint32 `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Id:       tokenID,
		Name:     userUpdate.Name,
		Username: userUpdate.Username,
		Email:    userUpdate.Email,
		Password: "Your Password",
	})
}
