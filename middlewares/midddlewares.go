package middlewares

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/inadislam/go-login-api/utils"
)

var secret = "cypher"

func JwtMiddlewares(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenstring := r.Header.Get("authorization")
		if tokenstring == "" {
			utils.ToJson(w, http.StatusUnauthorized, struct{
				Message string `json:"message"`
				Status int `json:"status"`
			}{
				Message: "Invalid token User Unauthorized",
				Status: http.StatusUnauthorized,
			})
		} else {
			result, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err == nil && result.Valid {
				next.ServeHTTP(w, r)
			} else {
				utils.ToJson(w, http.StatusUnauthorized, struct{
					Message string `json:"message"`
					Status int `json:"status"`
				}{
					Message: "Invalid token User Unauthorized",
					Status: http.StatusUnauthorized,
				}
			}
		}
	})
}