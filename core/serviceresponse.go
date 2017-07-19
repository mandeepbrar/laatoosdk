package core

const (
	StatusSuccess      = 200
	StatusServeFile    = 201
	StatusServeBytes   = 202
	StatusUnauthorized = 401
	StatusNotFound     = 404
	StatusRedirect     = 301
	StatusNotModified  = 305
	StatusBadRequest   = 400
)

/***Header****/
const (
	ContentType     = "Content-Type"
	ContentEncoding = "Content-Encoding"
	LastModified    = "Last-Modified"
)

func NewServiceResponse(status int, data interface{}, info ServiceParamsMap) *ServiceResponse {
	return newServiceResponse(status, data, info, false)
}
func newServiceResponse(status int, data interface{}, info ServiceParamsMap, ReturnVal bool) *ServiceResponse {
	return &ServiceResponse{status, data, info, ReturnVal}
}

var (
	StatusSuccessResponse      = newServiceResponse(StatusSuccess, nil, nil, true)
	StatusUnauthorizedResponse = newServiceResponse(StatusUnauthorized, nil, nil, true)
	StatusNotFoundResponse     = newServiceResponse(StatusNotFound, nil, nil, true)
	StatusBadRequestResponse   = newServiceResponse(StatusBadRequest, nil, nil, true)
	StatusNotModifiedResponse  = newServiceResponse(StatusNotModified, nil, nil, true)
)

func SuccessResponse(data interface{}) *ServiceResponse {
	return newServiceResponse(StatusSuccess, data, nil, true)
}

func BadRequestResponse(data string) *ServiceResponse {
	return newServiceResponse(StatusBadRequest, data, nil, true)
}

func UnauthorizedResponse(data string) *ServiceResponse {
	return newServiceResponse(StatusUnauthorized, data, nil, true)
}
