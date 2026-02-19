package core

import (
	"context"
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


// StreamChunk represents a single chunk in a streaming response.
type StreamChunk struct {
	Status   int
	Data     interface{}
	MetaInfo map[string]interface{}
	Error    error
	Final    bool // true when this is the last chunk
}

// ResponseStream is the producer-side interface a service uses to push data.
type ResponseStream struct {
	chunks chan *StreamChunk
	done   chan struct{}
	ctx    RequestContext
	cancel context.CancelFunc
}

func NewResponseStream(ctx RequestContext, bufferSize int) *ResponseStream {
	streamCtx, cancel := ctx.WithCancel()
	return &ResponseStream{
		chunks: make(chan *StreamChunk, bufferSize),
		done:   make(chan struct{}),
		ctx:    streamCtx.(RequestContext),
		cancel: cancel,
	}
}

// Send sends a chunk to the stream. Blocks if buffer is full.
func (rs *ResponseStream) Send(chunk *StreamChunk) error {
	select {
	case <-rs.ctx.Done():
		return rs.ctx.Err()
	case rs.chunks <- chunk:
		return nil
	}
}

// Close signals the end of the stream.
func (rs *ResponseStream) Close() {
	close(rs.chunks)
}

// Cancel aborts the stream from the consumer side.
func (rs *ResponseStream) Cancel() {
	rs.cancel()
}

// Chunks returns the read-only channel for consumers.
func (rs *ResponseStream) Chunks() <-chan *StreamChunk {
	return rs.chunks
}
