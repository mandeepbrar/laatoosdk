package errors

import (
	"fmt"
	"log/slog"
	"runtime/debug"

	"laatoo.io/sdk/ctx"
)

type ErrorLevel int

// levels at which error messages can be logged
const (
	FATAL ErrorLevel = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

// error that is registered for an error code
type Error struct {
	error
	InternalErrorCode string
	info              []slog.Attr
}

func (err *Error) UnderlyingError() error {
	return err.error
}
func (err *Error) Error() string {
	return fmt.Sprint("Root Error: ", err.error.Error(), ", Error code: ", err.InternalErrorCode, err.info)
}

var ShowStack = true

// Error handler for interrupting the error process
// Returns true if the error has been handled
type ErrorHandler func(ctx ctx.Context, err *Error, info ...slog.Attr) bool

var (
	//errors register to store all errors in the process
	ErrorsRegister = make(map[string]string, 50)
	//registered handlers for errors
	ErrorsHandlersRegister = make(map[string][]ErrorHandler, 20)
)

// register error code
func RegisterCode(internalErrorCode string, errMessage string) {
	ErrorsRegister[internalErrorCode] = errMessage
}

// register error handler for an internal error code
// handler will be called before throwing an error
// nil will be retured if an error is handled
func RegisterErrorHandler(internalErrorCode string, eh ErrorHandler) {
	val := ErrorsHandlersRegister[internalErrorCode]
	//add a new array of handlers if it doesnt exist already
	if val == nil {
		val = []ErrorHandler{}
	}
	//append the handler to the existing list and add to the map
	val = append(val, eh)
	ErrorsHandlersRegister[internalErrorCode] = val
}
func throwStandardError(ctx ctx.Context, internalErrorCode string, info ...slog.Attr) error {
	errMsg, ok := ErrorsRegister[internalErrorCode]
	if !ok {
		panic(fmt.Errorf("Invalid error code: %s", internalErrorCode))
	}
	return throwError(ctx, errMsg, internalErrorCode, nil, info...)
}
func ThrowError(ctx ctx.Context, message string, internalErrorCode string, info ...slog.Attr) error {
	return RethrowError(ctx, message, internalErrorCode, nil, info...)
}

func RethrowError(ctx ctx.Context, message string, internalErrorCode string, err error, info ...slog.Attr) error {
	return throwError(ctx, message, internalErrorCode, err, info...)
}

// throw a registered error code
// rethrow an error with an internal error code
func throwError(ctx ctx.Context, message, internalErrorCode string, thrownErr error, info ...slog.Attr) error {
	infoArr := []slog.Attr{slog.String("Context", ctx.GetName()), slog.String("Path", ctx.GetPath()), slog.String("ContextID", ctx.GetId())}
	if thrownErr != nil {
		rethrownError, ok := thrownErr.(*Error)
		if ok {
			infoArr = append(infoArr, slog.String("Root Error", rethrownError.error.Error()), slog.String("Root Error code", rethrownError.InternalErrorCode))
			infoArr = append(infoArr, rethrownError.info...)
		} else {
			infoArr = append(infoArr, slog.String("Stack", string(debug.Stack())), slog.String("Root Error", thrownErr.Error()))
		}
	} else {
		infoArr = append(infoArr, slog.String("Stack", string(debug.Stack())))
	}
	infoArr = append(infoArr, info...)
	/*	switch registeredError.Loglevel {
		case FATAL:
			log.Fatal(ctx, "Encountered error", infoArr...)
		case ERROR:
			log.Error(ctx, "Encountered error", infoArr...)
		case WARNING:
			log.Warn(ctx, "Encountered warning", infoArr...)
		case INFO:
			log.Info(ctx, "Info Error", infoArr...)
		case DEBUG:
			log.Debug(ctx, "Debug Error", infoArr...)
		}*/
	err := &Error{
		error:             fmt.Errorf(message),
		info:              infoArr,
		InternalErrorCode: internalErrorCode,
	}
	//call the handlers while throwing an error
	handlers := ErrorsHandlersRegister[internalErrorCode]
	if handlers != nil {
		handled := false
		for _, val := range handlers {
			handled = val(ctx, err, info...) || handled
		}
		//if an error has been handled, dont throw it
		if handled {
			return nil
		}
	}
	//thwo the error
	return err
}
