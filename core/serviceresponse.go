package core

import "fmt"

const (
	StatusSuccess         = 200
	StatusServeFile       = 201
	StatusServeBytes      = 202
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

func NewServiceResponse(status int, data map[string]interface{}) *Response {
	return &Response{status, data, nil, true}
}
func NewServiceResponseWithInfo(status int, data interface{}, info map[string]interface{}) *Response {
	var res map[string]interface{}
	if info != nil {
		res = info
	} else {
		res = make(map[string]interface{})
	}
	res["Data"] = data
	return newServiceResponse(status, res, nil, false)
}

func newServiceResponse(status int, data map[string]interface{}, err error, ReturnVal bool) *Response {
	return &Response{status, data, err, ReturnVal}
}

var (
	StatusSuccessResponse       = newServiceResponse(StatusSuccess, nil, nil, true)
	StatusUnauthorizedResponse  = newServiceResponse(StatusUnauthorized, nil, nil, true)
	StatusNotFoundResponse      = newServiceResponse(StatusNotFound, nil, nil, true)
	StatusNotModifiedResponse   = newServiceResponse(StatusNotModified, nil, nil, true)
	StatusInternalErrorResponse = newServiceResponse(StatusInternalError, nil, nil, true)
	StatusBadRequestResponse    = newServiceResponse(StatusBadRequest, nil, nil, true)
)

func SuccessResponse(data interface{}) *Response {
	return newServiceResponse(StatusSuccess, map[string]interface{}{"Data": data}, nil, false)
}

func FunctionalErrorResponse(err error) *Response {
	return newServiceResponse(StatusFunctionalError, nil, err, true)
}

func SuccessResponseWithInfo(data interface{}, info map[string]interface{}) *Response {
	return NewServiceResponseWithInfo(StatusSuccess, data, info)
}

func BadRequestResponse(err string) *Response {
	return newServiceResponse(StatusBadRequest, nil, fmt.Errorf(err), true)
}

func InternalErrorResponse(err string) *Response {
	return newServiceResponse(StatusInternalError, nil, fmt.Errorf(err), true)
}
func UnauthorizedResponse(err string) *Response {
	return newServiceResponse(StatusUnauthorized, nil, fmt.Errorf(err), true)
}
