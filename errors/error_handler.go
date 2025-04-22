package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/abdultalif/golang-auth-gorm/logger"
	"github.com/abdultalif/golang-auth-gorm/utils"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)


type CustomError struct {
	Code    int
	Status  string
	Message string
}

func (e CustomError) Error() string {
	return e.Message
}

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	switch e := err.(type) {
	case CustomError:
		sendErrorResponse(writer, e.Code, e.Status, e.Message)
	case validator.ValidationErrors:
		handleValidationErrors(writer, e)
	case error:
		sendErrorResponse(writer, http.StatusInternalServerError, "INTERNAL SERVER ERROR", e.Error())
	default:
		sendErrorResponse(writer, http.StatusInternalServerError, "INTERNAL SERVER ERROR", "An unknown error occurred")
	}
}

func sendErrorResponse(writer http.ResponseWriter, code int, status string, message interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)

	logger.Log.WithFields(logrus.Fields{
		"code":   code,
		"status": status,
		"error":  message,
	}).Error("Sending error response")
	errorResponse := utils.WebResponseError{
		Success: false,
		Code:    code,
		Status:  status,
		Error:   message,
	}

	utils.WriteToResponseBody(writer, errorResponse)
}

func handleValidationErrors(writer http.ResponseWriter, err validator.ValidationErrors) {
	errors := make(map[string]interface{})
	for _, e := range err {
		field := e.Field()
		fieldName := strings.ToLower(field[strings.LastIndex(field, ".")+1:])
		
		if _, ok := errors[fieldName]; !ok {
			errors[fieldName] = make([]string, 0)
		}
		
		var message string
		switch e.Tag() {
		case "required":
			message = "This field is required"
		case "eqfield":
			message = fmt.Sprintf("Field must be equal to %s", e.Param())
		case "email":
			message = "Invalid email format"
		case "len":
			message = fmt.Sprintf("Length must be %s", e.Param())
		case "min":
			message = fmt.Sprintf("Minimum length is %s", e.Param())
		case "max":
			message = fmt.Sprintf("Maximum length is %s", e.Param())
		case "uuid":
			message = "Invalid UUID format"
		default:
			message = fmt.Sprintf("Field validation failed on '%s' tag", e.Tag())
		}
		
		errors[fieldName] = append(errors[fieldName].([]string), message)
	}

	sendErrorResponse(writer, http.StatusBadRequest, "BAD REQUEST", errors)
}

