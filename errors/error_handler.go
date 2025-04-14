package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/abdultalif/golang-auth-gorm/utils"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if notFoundError(writer, request, err) {
		return
	}

	if badRequestError(writer, request, err) {
		return
	}

	if conflictError(writer, request, err) {
		return
	}

	internalServerError(writer, request, err)
}

func badRequestError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)

        errors := make(map[string]interface{})
        for _, e := range exception {
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
            default:
                message = fmt.Sprintf("Field validation failed on '%s' tag", e.Tag())
            }
            
            errors[fieldName] = append(errors[fieldName].([]string), message)
        }

        errorResponse := utils.WebResponseError{
            Success: false,
            Code:    http.StatusBadRequest,
            Status:  "BAD REQUEST",
            Error:   errors,
        }

		utils.WriteToResponseBody(writer, errorResponse)
		return true
	} else {
		return false
	}
}
func notFoundError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(NotFoundError)
	if ok {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)

		WebResponseError := utils.WebResponseError{
			Success: false,
			Code:    http.StatusNotFound,
			Status: "NOT FOUND",
			Error:   exception.Error,
		}

	utils.WriteToResponseBody(writer, WebResponseError)
		return true
	} else {
		return false
	}
}

func conflictError(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(DuplicateError)
	if ok {
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusConflict)

		WebResponseError := utils.WebResponseError{
			Success: false,
			Code:    http.StatusConflict,
			Status: "CONFLICT",
			Error:   exception.Error,
		}

	utils.WriteToResponseBody(writer, WebResponseError)
		return true
	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)

	WebResponseError := utils.WebResponseError{
		Success: false,
		Code:    http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Error:   err,
	}

	utils.WriteToResponseBody(writer, WebResponseError)
}