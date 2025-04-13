package routes

import "github.com/julienschmidt/httprouter"

func SetupRouter() *httprouter.Router {
	router := httprouter.New()

	// router.POST("/api/v1/auth")

	return router
}