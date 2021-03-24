package auth

import (
	"encoding/json"
	"errors"
	"html"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/inadislam/go-login-api/models"
)

func UserPrep(r *http.Request) (models.User, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return models.User{}, errors.New("request unprocessable")
	}
	var user models.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return models.User{}, errors.New("no json data found in request body")
	}
	user.ID = 0
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return user, nil
}
