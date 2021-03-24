package middlewares

import (
	"errors"
	"net/http"

	"github.com/inadislam/go-login-api/auth"
	"github.com/inadislam/go-login-api/utils"
)

func BasicMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json, charset=utf8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func IsAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			utils.ERROR(w, http.StatusUnauthorized, errors.New("request unauthorized"))
			return
		}
		next(w, r)
	})
}
