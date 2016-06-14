// +build !appengine

package log

/*
import (
	"github.com/Sirupsen/logrus"
	"laatoo/sdk/core"
)

func NewLogrus() LoggerInterface {
	return &LogrusLogger{logrus.New()}
}

type LogrusLogger struct {
	logger *logrus.Logger
}

func (log *LogrusLogger) Trace(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Debug([]interface{}{reqContext.GetName(), msg, args})
}
func (log *LogrusLogger) Debug(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Debug([]interface{}{reqContext.GetName(), msg, args})
}
func (log *LogrusLogger) Info(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Info([]interface{}{reqContext.GetName(), msg, args})
}
func (log *LogrusLogger) Warn(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Warn([]interface{}{reqContext.GetName(), msg, args})
}
func (log *LogrusLogger) Error(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Error([]interface{}{reqContext.GetName(), msg, args})
}
func (log *LogrusLogger) Fatal(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Fatal([]interface{}{reqContext.GetName(), msg, args})
}
func (log *LogrusLogger) SetFormat(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	case "happy":
	}
}
func (log *LogrusLogger) SetType(logtype string) {

}

func (log *LogrusLogger) SetLevel(level int) {
	switch level {
	case TRACE:
		log.logger.Level = logrus.DebugLevel
	case DEBUG:
		log.logger.Level = logrus.DebugLevel
	case INFO:
		log.logger.Level = logrus.InfoLevel
	case WARN:
		log.logger.Level = logrus.WarnLevel
	default:
		log.logger.Level = logrus.ErrorLevel
	}
}
func (log *LogrusLogger) IsTrace() bool {
	return log.logger.Level == logrus.DebugLevel
}
func (log *LogrusLogger) IsDebug() bool {
	return log.logger.Level == logrus.DebugLevel
}
func (log *LogrusLogger) IsInfo() bool {
	return log.logger.Level == logrus.InfoLevel
}
func (log *LogrusLogger) IsWarn() bool {
	return log.logger.Level == logrus.WarnLevel
}
*/
