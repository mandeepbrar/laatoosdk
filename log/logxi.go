// +build !appengine

package log

/*
import (
	//log "github.com/Sirupsen/logrus"
	"laatoosdk/core"
	logxi "logxi/v1"
	"os"
)

func NewLogxiLogger() LoggerInterface {
	return &LogxiLogger{logxi.NewLogger3(os.Stdout, "default", logxi.NewJSONFormatter("default"))}
}

type LogxiLogger struct {
	logger logxi.Logger
}

func (log *LogxiLogger) Trace(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Trace(reqContext, reqContext.GetName(), msg, args...)
}
func (log *LogxiLogger) Debug(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Debug(reqContext, reqContext.GetName(), msg, args...)
}
func (log *LogxiLogger) Info(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Info(reqContext, reqContext.GetName(), msg, args...)
}
func (log *LogxiLogger) Warn(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Warn(reqContext, reqContext.GetName(), msg, args...)
}
func (log *LogxiLogger) Error(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Error(reqContext, reqContext.GetName(), msg, args...)
}
func (log *LogxiLogger) Fatal(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Fatal(reqContext, reqContext.GetName(), msg, args...)
}
func (log *LogxiLogger) SetFormat(format string) {
	switch format {
	case "json":
	case "happy":
		logger := log.logger.(*logxi.DefaultLogger)
		logger.SetFormatter(logxi.NewHappyDevFormatter("default"))
	}
}

func (log *LogxiLogger) SetType(logtype string) {

}

func (log *LogxiLogger) SetLevel(level int) {
	switch level {
	case TRACE:
		log.logger.SetLevel(logxi.LevelTrace)
	case DEBUG:
		log.logger.SetLevel(logxi.LevelDebug)
	case INFO:
		log.logger.SetLevel(logxi.LevelInfo)
	case WARN:
		log.logger.SetLevel(logxi.LevelWarn)
	default:
		log.logger.SetLevel(logxi.LevelError)
	}
}
func (log *LogxiLogger) IsTrace() bool {
	return log.logger.IsTrace()
}
func (log *LogxiLogger) IsDebug() bool {
	return log.logger.IsDebug()
}
func (log *LogxiLogger) IsInfo() bool {
	return log.logger.IsInfo()
}
func (log *LogxiLogger) IsWarn() bool {
	return log.logger.IsWarn()
}
*/
