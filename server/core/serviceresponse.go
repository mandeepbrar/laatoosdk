package core

import (
	"fmt"

	"laatoo.io/sdk/utils"
)

const (
	StatusSuccess         = 200
	StatusMoreInfo        = 250
	StatusServeFile       = 251
	StatusServeBytes      = 252
	StatusServeStream     = 253
	StatusUnauthorized    = 401
	StatusNotFound        = 404
	StatusRedirect        = 301
	StatusNotModified     = 305
	StatusBadRequest      = 400
	StatusInternalError   = 500
	StatusFunctionalError = 550
)

/***Header****/
const (
	ContentType     = "Content-Type"
	ContentEncoding = "Content-Encoding"
	LastModified    = "Last-Modified"
)

func NewServiceResponse(status int, data interface{}) *Response {
	return newServiceResponse(status, data, nil, nil, true)
}
func NewServiceResponseWithInfo(status int, data interface{}, info utils.StringMap) *Response {
	return newServiceResponse(status, data, info, nil, false)
}

func newServiceResponse(status int, data interface{}, info utils.StringMap, err error, ReturnVal bool) *Response {
	return &Response{status, data, info, err, ReturnVal}
}

var (
	StatusSuccessResponse      = newServiceResponse(StatusSuccess, nil, nil, nil, true)
	StatusUnauthorizedResponse = newServiceResponse(StatusUnauthorized, nil, nil, nil, true)
	StatusNotFoundResponse     = newServiceResponse(StatusNotFound, nil, nil, nil, true)
	StatusNotModifiedResponse  = newServiceResponse(StatusNotModified, nil, nil, nil, true)
)

func SuccessResponse(data interface{}) *Response {
	return newServiceResponse(StatusSuccess, data, nil, nil, false)
}

func RedirectResponse(data interface{}) *Response {
	return newServiceResponse(StatusRedirect, data, nil, nil, false)
}

func FunctionalErrorResponse(err error) *Response {
	return newServiceResponse(StatusFunctionalError, nil, nil, err, true)
}

func SuccessResponseWithInfo(data interface{}, info utils.StringMap) *Response {
	return NewServiceResponseWithInfo(StatusSuccess, data, info)
}
func SuccessServeBytes(data []byte) *Response {
	return newServiceResponse(StatusSuccess, data, nil, nil, false)
}
func BadRequestResponse(err string) *Response {
	return newServiceResponse(StatusBadRequest, nil, nil, fmt.Errorf(err), true)
}

func InternalErrorResponse(err string) *Response {
	return newServiceResponse(StatusInternalError, nil, nil, fmt.Errorf(err), true)
}
func UnauthorizedResponse(err string) *Response {
	return newServiceResponse(StatusUnauthorized, nil, nil, fmt.Errorf(err), true)
}
