package utils

import (
	"encoding/json"
	"net/http"

	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/sirupsen/logrus"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to read request body")
		panic(err)
	}
}


func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	enncoder := json.NewEncoder(writer)
	err := enncoder.Encode(response)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to write response body")
		panic(err)
	}
}