package utils

type WebResponseSuccess struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type WebResponseError struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Error   interface{} `json:"error"`
}