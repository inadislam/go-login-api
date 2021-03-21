package middlewares

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/inadislam/go-login-api/models"
)

var secret = "cypher"

func JwtToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorization"] = true
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	return token.SignedString([]byte(secret))
}
