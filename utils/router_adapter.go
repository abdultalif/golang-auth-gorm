package utils

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Adapt(fn func(http.ResponseWriter, *http.Request, httprouter.Params)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, httprouter.Params{})
	}
}