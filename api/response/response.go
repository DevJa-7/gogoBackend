package response

import (
	"net/http"

	"github.com/labstack/echo"
)

// BasicResponse is a basic response.
type BasicResponse struct {
	Status      int    `json:"status" description:"HTTP Status"`
	MessageType string `json:"messageType" description:"Message type code"`
	Message     string `json:"message" description:"Returned message"`
}

// MessageTypes is an array of message types that using in frontend.
// OK                  : 200
// Created             : 201
// NotModified         : 304
// BadRequest          : 400
// Unauthorized        : 401
// PaymentRequired     : 402
// Forbidden           : 403
// NotFound            : 404
// MethodNotAllowed    : 405
// InternalServerError : 500
type MessageTypes struct {
	OK                  string `description:"200"`
	Created             string `description:"201"`
	NotModified         string `description:"304"`
	BadRequest          string `description:"400"`
	Unauthorized        string `description:"401"`
	PaymentRequired     string `description:"402"`
	Forbidden           string `description:"403"`
	NotFound            string `description:"404"`
	MethodNotAllowed    string `description:"405"`
	InternalServerError string `description:"500"`
}

// Messages is an array of messages
type Messages struct {
	OK      string `description:"200"`
	Created string `description:"201"`
}

// SuccessInterface is a json response of 2xx.
func SuccessInterface(c echo.Context, msg interface{}) error {
	status := http.StatusOK
	return c.JSON(status, msg)
}

// SuccessJSON is a json response of 2xx.
func SuccessJSON(c echo.Context, msg string) error {
	status := http.StatusOK
	messageType := "success"
	return c.JSON(status, Success(status, messageType, msg))
}

// RedirectionJSON is a json response of 3xx.
func RedirectionJSON(c echo.Context, status int, messageType string, msg string) error {
	return c.JSON(status, Redirection(status, messageType, msg))
}

// JSON is a json response.
func JSON(c echo.Context, status int, messageTypes *MessageTypes, messages *Messages, err error) error {
	if err == nil {
		return SuccessJSON(c, messages.OK)
	} else if err != nil {
		return ErrorJSON(c, status, messageTypes, err)
	}
	return nil
}

// ErrorJSON is a json response of errors.
func ErrorJSON(c echo.Context, status int, messageTypes *MessageTypes, err error) error {
	switch status {
	case http.StatusNotModified:
		return KnownErrorJSON(c, status, messageTypes.NotModified, err)
	case http.StatusBadRequest:
		return KnownErrorJSON(c, status, messageTypes.BadRequest, err)
	case http.StatusUnauthorized:
		return KnownErrorJSON(c, status, messageTypes.Unauthorized, err)
	case http.StatusPaymentRequired:
		return KnownErrorJSON(c, status, messageTypes.PaymentRequired, err)
	case http.StatusForbidden:
		return KnownErrorJSON(c, status, messageTypes.Forbidden, err)
	case http.StatusNotFound:
		return KnownErrorJSON(c, status, messageTypes.NotFound, err)
	case http.StatusMethodNotAllowed:
		return KnownErrorJSON(c, status, messageTypes.MethodNotAllowed, err)
	case http.StatusInternalServerError:
		return KnownErrorJSON(c, status, messageTypes.InternalServerError, err)
	default:
		return UnknownErrorJSON(c, status, err)
	}
}

// KnownErrInterface is json response of known error.
func KnownErrInterface(c echo.Context, err interface{}) error {
	status := http.StatusBadRequest
	return c.JSON(status, err)
}

// KnownErrJSON is json response of known error.
func KnownErrJSON(c echo.Context, messageType string, err error) error {
	status := http.StatusBadRequest
	return c.JSON(status, KnownError(status, messageType, err))
}

// KnownErrorJSON is json response of known error.
func KnownErrorJSON(c echo.Context, status int, messageType string, err error) error {
	return c.JSON(status, KnownError(status, messageType, err))
}

// UnknownErrorJSON is json response of unknown error.
func UnknownErrorJSON(c echo.Context, status int, err error) error {
	return c.JSON(status, UnknownError(status, err))
}

// Success is a basic response of 2xx.
func Success(status int, messageType string, msg string) *BasicResponse {
	return &BasicResponse{Status: status, MessageType: messageType, Message: msg}
}

// Redirection is a basic response of 3xx.
func Redirection(status int, messageType string, msg string) *BasicResponse {
	return &BasicResponse{Status: status, MessageType: messageType, Message: msg}
}

// KnownError is a basic response of know errors.
func KnownError(status int, messageType string, err error) *BasicResponse {
	return &BasicResponse{Status: status, MessageType: messageType, Message: err.Error()}
}

// UnknownError is a basic response of unknown errors.
func UnknownError(status int, err error) *BasicResponse {
	return KnownError(status, "error.unknown", err)
}
