package middlewares

import (
	"net/http"
	"strings"

	"github.com/abdultalif/golang-auth-gorm/errors"
	"github.com/abdultalif/golang-auth-gorm/utils"
)

func JWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			panic(errors.CustomError{
				Code:    http.StatusUnauthorized,
				Status:  "UNAUTHORIZED",
				Message: "Missing or invalid Authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.VerifyToken(tokenString, false)
		if err != nil {
			panic(errors.CustomError{
				Code: http.StatusUnauthorized,
				Status: "UNAUTHORIZED",
				Message: "Invalid or expired token",
			})
		}
		
		r.Header.Set("X-User-ID", claims ["id"].(string))
		next(w, r)

	}
}