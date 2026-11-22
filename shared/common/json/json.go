package json

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct{}

type APIResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) Ok(ctx echo.Context, message string, data interface{}) error {
	return ctx.JSON(http.StatusOK, APIResponse{
		Code:    http.StatusOK,
		Status:  "SUCCESS",
		Message: message,
		Data:    data,
	})
}

func (r *Response) Created(ctx echo.Context, message string, data interface{}) error {
	return ctx.JSON(http.StatusCreated, APIResponse{
		Code:    http.StatusCreated,
		Status:  "SUCCESS",
		Message: message,
		Data:    data,
	})
}

func (r *Response) BadRequest(ctx echo.Context, message string, err interface{}) error {
	var errDetails interface{}
	if e, ok := err.(error); ok {
		errDetails = e.Error()
	} else {
		errDetails = err
	}
	return ctx.JSON(http.StatusBadRequest, APIResponse{
		Code:    http.StatusBadRequest,
		Status:  "BAD_REQUEST",
		Message: message,
		Errors:  errDetails,
	})
}

func (r *Response) Unauthorized(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusUnauthorized, APIResponse{
		Code:    http.StatusUnauthorized,
		Status:  "UNAUTHORIZED",
		Message: message,
	})
}

func (r *Response) NotFound(ctx echo.Context, message string) error {
	return ctx.JSON(http.StatusNotFound, APIResponse{
		Code:    http.StatusNotFound,
		Status:  "NOT_FOUND",
		Message: message,
	})
}

func (r *Response) InternalServerError(ctx echo.Context, message string, err error) error {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}
	return ctx.JSON(http.StatusInternalServerError, APIResponse{
		Code:    http.StatusInternalServerError,
		Status:  "INTERNAL_SERVER_ERROR",
		Message: message,
		Errors:  errStr,
	})
}
