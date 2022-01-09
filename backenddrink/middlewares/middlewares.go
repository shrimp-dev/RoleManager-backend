package middlewares

import (
	"context"
	"drinkBack/models"
	"drinkBack/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func Cors() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			w.Header().Set("Access-Control-Expose-Headers", "Authorization")
			if req.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func Authenticate() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			b := req.Header.Get("Authorization")
			token := strings.Replace(b, "Bearer ", "", 1)
			var usr models.AccessTokenClaims
			ok, err := utils.VerifyAuthenticationToken(token, &usr)
			if err != nil || !ok {
				http.Error(w, "Error trying to authenticate you", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(req.Context(), "usrToken", usr)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
