package routes

import (
	"net/http"

	"github.com/abdultalif/golang-auth-gorm/controllers"
	"github.com/abdultalif/golang-auth-gorm/errors"
	"github.com/abdultalif/golang-auth-gorm/middlewares"
	"github.com/abdultalif/golang-auth-gorm/utils"

	"github.com/julienschmidt/httprouter"
)

func ProtectedRoute(fn func(http.ResponseWriter, *http.Request, httprouter.Params)) http.HandlerFunc {
	return middlewares.JWTAuth(utils.Adapt(fn))
}

func SetupRouter() *httprouter.Router {	
	router := httprouter.New()

	router.POST("/api/v1/auth", controllers.Register)
	router.POST("/api/v1/auth/verify-otp", controllers.VerifyOTP)
	router.POST("/api/v1/auth/resend-code", controllers.ResendCode)
	router.POST("/api/v1/auth/login", controllers.Login)
	router.POST("/api/v1/auth/refresh-token", controllers.RefreshToken)
	router.POST("/api/v1/auth/forgot-password", controllers.ForgotPassword)
	router.POST("/api/v1/auth/check-token", controllers.CheckToken)
	router.POST("/api/v1/auth/reset-password", controllers.ResetPassword)

	router.Handler("GET", "/api/v1/user/profile", ProtectedRoute(controllers.GetProfile))



	router.PanicHandler = errors.ErrorHandler
	return router
}
