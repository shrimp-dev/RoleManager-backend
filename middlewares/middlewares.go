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
			if ok, err := utils.VerifyAuthenticationToken(token, utils.AUTH, &usr); err != nil || !ok {
				http.Error(w, "Error trying to authenticate you", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(req.Context(), "usrToken", usr)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func CreateUserAuthenticate() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)
			token := vars["token"]
			var creator models.AccessTokenClaims
			if ok, err := utils.VerifyAuthenticationToken(token, utils.INVITE, &creator); err != nil || !ok {
				http.Error(w, "Error trying to authenticate your invite", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(req.Context(), "creator", creator)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
