package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
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
		code := database.GetOTp(userCreated.ID)
		fmt.Println(code)
		otp := strconv.FormatInt(code, 10)
		auth.ActiveUser(otp, userCreated.Email, userCreated.Username)
		utils.ToJson(w, http.StatusCreated, struct {
			Id       uint32 `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Active   bool   `json:"active"`
			Message  string `json:"message"`
		}{
			Id:       userCreated.ID,
			Name:     userCreated.Name,
			Username: userCreated.Username,
			Email:    userCreated.Email,
			Password: "Your Password",
			Active:   userCreated.Active,
			Message:  "Check your Email Box for Verification Code",
		})
	} else {
		utils.ERROR(w, http.StatusConflict, errors.New("username or email aleready exists"))
		return
	}
}

func ActiveUser(w http.ResponseWriter, r *http.Request) {
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
	if users.ID == 0 || users.Otp == 0 {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("fields must not be empty"))
		return
	}
	user, err := database.UserById(users.ID)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	fmt.Println(users.Otp)
	fmt.Println(user)
	if users.Otp != user.Otp {
		utils.ERROR(w, http.StatusInternalServerError, errors.New("invalid otp"))
		return
	} else {
		if user.Active == true {
			utils.ERROR(w, http.StatusInternalServerError, errors.New("otp expired"))
		}
		err = database.UserActive(user.ID)
		if err != nil {
			utils.ERROR(w, http.StatusInternalServerError, errors.New("failed to active account, internal server error"))
		}
		utils.ToJson(w, http.StatusOK, struct {
			Id       uint32 `json:"id"`
			Username string `json:"username"`
			Name     string `json:"name"`
			Password string `json:"password"`
			Message  string `json:"message"`
		}{
			Id:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Password: "Your Password",
			Message:  "Your Account Activated.Login Now",
		})
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
		if ud.Active == false {
			utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("Active Your Account First"))
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
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
	}
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

func FindYourAccount(w http.ResponseWriter, r *http.Request) {
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
	if user.Email == "" {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("fields must not be empty"))
		return
	}
	if err = checkmail.ValidateFormat(user.Email); err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("invalid email format"))
		return
	}
	userInfo, err := database.UserByEmail(user.Email)
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, errors.New("user not found, create a account first"))
		return
	} else {
		if userInfo.Active == false {
			utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("your account is not active, you cant change your password before activate your account"))
			return
		}
		otp := database.GetOTp(userInfo.ID)
		fmt.Println(otp)
		codsSt := strconv.FormatInt(otp, 10)
		auth.ForgetMail(codsSt, user.Email, userInfo.Username)
		utils.ToJson(w, http.StatusCreated, struct {
			Email   string `json:"email"`
			Message string `json:"message"`
		}{
			Email:   user.Email,
			Message: "Check your Mail for OTP",
		})
	}
}

func ForgetPassword(w http.ResponseWriter, r *http.Request) {
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
	if user.Otp == 0 || user.ID == 0 || user.Password == "" {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("Fields Cannot be Empty"))
		return
	}
	userInf, err := database.UserById(user.ID)
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}
	if user.Otp != userInf.Otp {
		utils.ERROR(w, http.StatusUnprocessableEntity, errors.New("Invalid OTP"))
		return
	} else {
		err = database.ChangePassword(user.ID, user.Password)
		if err != nil {
			utils.ERROR(w, http.StatusInternalServerError, errors.New("Internal Server Error"))
			return
		}
		utils.ToJson(w, http.StatusCreated, struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Message  string `json:"message"`
		}{
			Username: userInf.Username,
			Password: "******",
			Message:  "Your Password has changed. Go and Login",
		})
	}
}
