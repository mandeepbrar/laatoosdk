package log

import (
	"laatoo/sdk/core"
)

const (
	STR_TRACE = "Trace"
	STR_DEBUG = "Debug"
	STR_INFO  = "Info"
	STR_WARN  = "Warn"
	STR_ERROR = "Error"
	STR_FATAL = "Fatal"
)

type logPrinter func(ctx core.Context, app string, strlevel string, wh SimpleWriteHandler, level int, msg string, args ...interface{})

var (
	logFormats = make(map[string]logPrinter, 5)
)

type SimpleWriteHandler interface {
	Print(ctx core.Context, app string, msg string, level int, strlevel string)
	PrintBytes(ctx core.Context, app string, msg []byte, level int, strlevel string) (int, error)
}

func NewSimpleLogger(wh SimpleWriteHandler) LoggerInterface {
	logger := &SimpleLogger{format: "json", level: INFO, wh: wh, app: "Laatoo"}
	logger.printer = printJSON
	return logger
}

type SimpleLogger struct {
	wh  SimpleWriteHandler
	app string
	//buffer bytes.Buffer
	format  string
	printer logPrinter
	level   int
}

func (log *SimpleLogger) Trace(ctx core.Context, msg string, args ...interface{}) {
	if log.level > DEBUG {
		log.printer(ctx, log.app, STR_TRACE, log.wh, TRACE, msg, args...)
	}
}
func (log *SimpleLogger) Debug(ctx core.Context, msg string, args ...interface{}) {
	if log.level > INFO {
		log.printer(ctx, log.app, STR_DEBUG, log.wh, DEBUG, msg, args...)
	}
}
func (log *SimpleLogger) Info(ctx core.Context, msg string, args ...interface{}) {
	if log.level > WARN {
		log.printer(ctx, log.app, STR_INFO, log.wh, INFO, msg, args...)
	}
}
func (log *SimpleLogger) Warn(ctx core.Context, msg string, args ...interface{}) {
	if log.level > ERROR {
		log.printer(ctx, log.app, STR_WARN, log.wh, WARN, msg, args...)
	}
}
func (log *SimpleLogger) Error(ctx core.Context, msg string, args ...interface{}) {
	if log.level > FATAL {
		log.printer(ctx, log.app, STR_ERROR, log.wh, ERROR, msg, args...)
	}
}
func (log *SimpleLogger) Fatal(ctx core.Context, msg string, args ...interface{}) {
	log.printer(ctx, log.app, STR_FATAL, log.wh, FATAL, msg, args...)
}

func (log *SimpleLogger) SetFormat(format string) {
	log.format = format
	printer, ok := logFormats[format]
	if ok {
		log.printer = printer
	}
}

func (log *SimpleLogger) SetType(loggertype string) {
}

func (log *SimpleLogger) SetApplication(app string) {
	log.app = app
}
func (log *SimpleLogger) GetApplication() string {
	return log.app
}
func (log *SimpleLogger) SetLevel(level int) {
	log.level = level
}
func (log *SimpleLogger) IsTrace() bool {
	return log.level == TRACE
}
func (log *SimpleLogger) IsDebug() bool {
	return log.level == DEBUG
}
func (log *SimpleLogger) IsInfo() bool {
	return log.level == INFO
}
func (log *SimpleLogger) IsWarn() bool {
	return log.level == WARN
}

func (log *SimpleLogger) Write(p []byte) (int, error) {
	return log.wh.PrintBytes(nil, log.app, p, INFO, STR_INFO)
}
