package routes

import (
	"net/http"

	"github.com/abdultalif/golang-auth-gorm/controllers"
	"github.com/abdultalif/golang-auth-gorm/errors"
	"github.com/abdultalif/golang-auth-gorm/middlewares"
	"github.com/julienschmidt/httprouter"
)

func SetupRouter() *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/v1/auth", controllers.Register)
	router.POST("/api/v1/auth/verify", controllers.VerifyUser)
	router.POST("/api/v1/auth/resend-code", controllers.ResendCode)
	router.POST("/api/v1/auth/login", controllers.Login)
	router.POST("/api/v1/auth/refresh-token", controllers.RefreshToken)

	router.Handler("GET", "/api/v1/user/profile", middlewares.JWTAuth(func(w http.ResponseWriter, r *http.Request) {
		controllers.GetProfile(w, r, httprouter.Params{})
	}))


	router.PanicHandler = errors.ErrorHandler
	return router
}
