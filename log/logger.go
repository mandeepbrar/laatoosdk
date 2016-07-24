package log

import (
	"io"
	"laatoo/sdk/config"
	"laatoo/sdk/core"
)

const (
	CONF_LOGGING        = "logging"
	CONF_LOGGINGLEVEL   = "level"
	CONF_LOGGER         = "loggertype"
	CONF_APPLICATION    = "application"
	CONF_LOGGING_FORMAT = "format"
	FATAL               = 1
	ERROR               = 2
	WARN                = 3
	INFO                = 4
	DEBUG               = 5
	TRACE               = 6
)

type LoggerInterface interface {
	io.Writer
	Trace(reqContext core.Context, msg string, args ...interface{})
	Debug(reqContext core.Context, msg string, args ...interface{})
	Info(reqContext core.Context, msg string, args ...interface{})
	Warn(reqContext core.Context, msg string, args ...interface{})
	Error(reqContext core.Context, msg string, args ...interface{})
	Fatal(reqContext core.Context, msg string, args ...interface{})

	SetLevel(int)
	SetApplication(string)
	GetApplication() string
	SetType(string)
	SetFormat(string)
	IsTrace() bool
	IsDebug() bool
	IsInfo() bool
	IsWarn() bool
}

var (
	Logger LoggerInterface
)

func init() {
	Logger = NewLogger()
	Logger.SetLevel(TRACE)
}

func GetLevel(level string) int {
	switch level {
	case "all":
		return TRACE
	case "trace":
		return TRACE
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	default:
		return ERROR
	}
}

//configures logging level and returns true if its set to debug
func ConfigLogger(conf config.Config) bool {
	logconf, ok := conf.GetSubConfig(CONF_LOGGING)
	if ok {
		application, _ := logconf.GetString(CONF_APPLICATION)
		if application != "" {
			Logger.SetApplication(application)
		}
		loggerType, _ := logconf.GetString(CONF_LOGGER)
		if loggerType != "" {
			Logger.SetType(loggerType)
		}
		loggingFormat, _ := logconf.GetString(CONF_LOGGING_FORMAT)
		if loggingFormat != "" {
			Logger.SetFormat(loggingFormat)
		}
		loggingLevel, _ := logconf.GetString(CONF_LOGGINGLEVEL)
		if loggingLevel != "" {
			Logger.SetLevel(GetLevel(loggingLevel))
		}
	}
	return (Logger.IsTrace()) || (Logger.IsDebug())
}
