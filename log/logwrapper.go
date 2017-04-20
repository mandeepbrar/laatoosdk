package log

import (
	"laatoo/sdk/core"
)

type LogWrapper struct {
	logger LoggerInterface
	level  int
}

func (log *LogWrapper) Trace(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Trace(reqContext, msg, args...)
}
func (log *LogWrapper) Debug(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Debug(reqContext, msg, args...)
}
func (log *LogWrapper) Info(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Info(reqContext, msg, args...)
}
func (log *LogWrapper) Warn(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Warn(reqContext, msg, args...)
}
func (log *LogWrapper) Error(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Error(reqContext, msg, args...)
}
func (log *LogWrapper) Fatal(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Fatal(reqContext, msg, args...)
}
func (log *LogWrapper) SetFormat(format string) {
	log.logger.SetFormat(format)
}

func (log *LogWrapper) Write(p []byte) (n int, err error) {
	return log.logger.Write(p)
}

func (log *LogWrapper) SetType(loggertype string) {
	if loggertype == "syslog" {
		log.logger = NewSysLogger(log.logger.GetApplication())
	}
	/*if loggertype == "logrus" {
		log.logger = NewLogrus()
	}
	if loggertype == "logxi" {
		log.logger = NewLogxiLogger()
	}*/
}

func (log *LogWrapper) SetApplication(app string) {
	log.logger.SetApplication(app)
}
func (log *LogWrapper) GetApplication() string {
	return log.logger.GetApplication()
}
func (log *LogWrapper) SetLevel(level int) {
	log.level = level
	log.logger.SetLevel(level)
}
func (log *LogWrapper) IsTrace() bool {
	return log.logger.IsTrace()
}
func (log *LogWrapper) IsDebug() bool {
	return log.logger.IsDebug()
}
func (log *LogWrapper) IsInfo() bool {
	return log.logger.IsInfo()
}
func (log *LogWrapper) IsWarn() bool {
	return log.logger.IsWarn()
}
