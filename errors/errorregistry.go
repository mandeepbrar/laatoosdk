package errors

import (
	"fmt"
	"laatoo/sdk/core"
	"laatoo/sdk/log"
	"runtime/debug"
)

type ErrorLevel int

//levels at which error messages can be logged
const (
	FATAL ErrorLevel = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

//error that is registered for an error code
type Error struct {
	error
	InternalErrorCode string
	Loglevel          ErrorLevel
}

var ShowStack = true

//Error handler for interrupting the error process
//Returns true if the error has been handled
type ErrorHandler func(ctx core.Context, err *Error, info ...interface{}) bool

var (
	//errors register to store all errors in the process
	ErrorsRegister = make(map[string]*Error, 50)
	//registered handlers for errors
	ErrorsHandlersRegister = make(map[string][]ErrorHandler, 20)
)

//register error code
func RegisterCode(internalErrorCode string, loglevel ErrorLevel, err error) {
	ErrorsRegister[internalErrorCode] = &Error{err, internalErrorCode, loglevel}
}

//register error handler for an internal error code
//handler will be called before throwing an error
//nil will be retured if an error is handled
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

func WrapError(ctx core.Context, err error) error {
	if err != nil {
		_, ok := err.(Error)
		if ok {
			return err
		} else {
			return RethrowError(ctx, CORE_ERROR_WRAPPER, err)
		}
	}
	return nil
}

func ThrowError(ctx core.Context, internalErrorCode string, info ...interface{}) error {
	return RethrowError(ctx, internalErrorCode, nil, info...)
}

func RethrowError(ctx core.Context, internalErrorCode string, err error, info ...interface{}) error {
	registeredErr, ok := ErrorsRegister[internalErrorCode]
	if !ok {
		panic(fmt.Errorf("Invalid error code: %s", internalErrorCode))
	}
	return throwError(ctx, registeredErr, err, info...)
}

//throw a registered error code
//rethrow an error with an internal error code
func throwError(ctx core.Context, registeredError *Error, rethrownError error, info ...interface{}) error {
	var errDetails []interface{}
	if rethrownError == nil {
		errDetails = []interface{}{"Err", registeredError.Error(), "Internal Error Code", registeredError.InternalErrorCode}
	} else {
		errDetails = []interface{}{"Err", registeredError.Error(), "Internal Error Code", registeredError.InternalErrorCode, "Root Error", rethrownError}
	}
	var infoArr []interface{}
	if ShowStack {
		infoArr = append(errDetails, info...)
		infoArr = append(infoArr, "Stack", string(debug.Stack()))
	} else {
		infoArr = append(errDetails, info...)
	}
	switch registeredError.Loglevel {
	case FATAL:
		log.Logger.Fatal(ctx, "Encountered error", infoArr...)
	case ERROR:
		log.Logger.Error(ctx, "Encountered error", infoArr...)
	case WARNING:
		log.Logger.Warn(ctx, "Encountered warning", infoArr...)
	case INFO:
		log.Logger.Info(ctx, "Info Error", infoArr...)
	case DEBUG:
		log.Logger.Debug(ctx, "Debug Error", infoArr...)
	}
	//call the handlers while throwing an error
	handlers := ErrorsHandlersRegister[registeredError.InternalErrorCode]
	if handlers != nil {
		handled := false
		for _, val := range handlers {
			handled = val(ctx, registeredError, info...) || handled
		}
		//if an error has been handled, dont throw it
		if handled {
			return nil
		}
	}
	//thwo the error
	return *registeredError
}
