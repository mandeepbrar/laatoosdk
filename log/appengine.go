// +build appengine

package log

import (
	glog "google.golang.org/appengine/log"
	"laatoo/sdk/core"
	stdlog "log"
)

func NewLogger() LoggerInterface {
	return &LogWrapper{logger: NewSimpleLogger(gaeSimpleLogsHandler()), level: TRACE}
}

func gaeSimpleLogsHandler() SimpleWriteHandler {
	wh := &gaeSimpleWriteHandler{}
	return wh
}

type gaeSimpleWriteHandler struct {
}

func (jh *gaeSimpleWriteHandler) Print(ctx core.Context, app string, msg string, level int, strlevel string) {
	if ctx != nil {
		appengineContext := ctx.GetAppengineContext()
		if appengineContext != nil {
			switch level {
			case TRACE:
				glog.Debugf(appengineContext, msg)
			case DEBUG:
				glog.Debugf(appengineContext, msg)
			case INFO:
				glog.Infof(appengineContext, msg)
			case WARN:
				glog.Warningf(appengineContext, msg)
			default:
				glog.Errorf(appengineContext, msg)
			}
			return
		}
	}
	stdlog.Print(msg)
}
func (jh *gaeSimpleWriteHandler) PrintBytes(ctx core.Context, app string, msg []byte, level int, strlevel string) (int, error) {
	stdlog.Print(string(msg))
	return len(msg), nil
}

/*
import (
	"bytes"
	glog "google.golang.org/appengine/log"
	"laatoosdk/core"
	logxi "logxi/v1"
	"os"
)

func NewLogger() LoggerInterface {
	return &StandaloneLogger{logger: NewJSONLogger(stdJsonLogsHandler()), level: TRACE}
	//return &StandaloneLogger{logxi.NewLoggerWithHandler(&AppEngineHandler{}, "default", logxi.NewJSONFormatter("default"))}
}

type StandaloneLogger struct {
	logger logxi.Logger
}

func (log *StandaloneLogger) Trace(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Trace(reqContext, msg, args...)
}
func (log *StandaloneLogger) Debug(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Debug(reqContext, msg, args...)
}
func (log *StandaloneLogger) Info(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Info(reqContext, msg, args...)
}
func (log *StandaloneLogger) Warn(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Warn(reqContext, msg, args...)
}
func (log *StandaloneLogger) Error(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Error(reqContext, msg, args...)
}
func (log *StandaloneLogger) Fatal(reqContext core.Context, msg string, args ...interface{}) {
	log.logger.Fatal(reqContext, msg, args...)
}

func (log *StandaloneLogger) SetFormat(format string) {
	switch format {
	case "json":
	case "happy":
		logger := log.logger.(*logxi.DefaultLogger)
		logger.SetFormatter(logxi.NewHappyDevFormatter("default"))
	}
}

func (log *StandaloneLogger) SetLevel(level string) {
	switch level {
	case "all":
		log.logger.SetLevel(logxi.LevelAll)
	case "trace":
		log.logger.SetLevel(logxi.LevelTrace)
	case "debug":
		log.logger.SetLevel(logxi.LevelDebug)
	case "info":
		log.logger.SetLevel(logxi.LevelInfo)
	case "warn":
		log.logger.SetLevel(logxi.LevelWarn)
	default:
		log.logger.SetLevel(logxi.LevelError)
	}
}

func (log *StandaloneLogger) SetType(logtype string) {

}

func (log *StandaloneLogger) IsTrace() bool {
	return log.logger.IsTrace()
}
func (log *StandaloneLogger) IsDebug() bool {
	return log.logger.IsDebug()
}
func (log *StandaloneLogger) IsInfo() bool {
	return log.logger.IsInfo()
}
func (log *StandaloneLogger) IsWarn() bool {
	return log.logger.IsWarn()
}

type AppEngineHandler struct {
}

func (ah *AppEngineHandler) WriteLog(ctx interface{}, loggingCtx string, buf *bytes.Buffer, level int, msg string, args []interface{}) {
	if ctx != nil {
		ectx, _ := ctx.(core.Context)
		appengineContext := ectx.GetAppengineContext()
		if appengineContext != nil {
			switch level {
			case logxi.LevelTrace:
				glog.Debugf(appengineContext, buf.String())
			case logxi.LevelDebug:
				glog.Debugf(appengineContext, buf.String())
			case logxi.LevelInfo:
				glog.Infof(appengineContext, buf.String())
			case logxi.LevelWarn:
				glog.Warningf(appengineContext, buf.String())
			default:
				glog.Errorf(appengineContext, buf.String())
			}
		} else {
			buf.WriteTo(os.Stderr)
		}
	}
}

func stdJsonLogsHandler() JsonWriteHandler {
	wh := &StdJsonWriteHandler{}
	return wh
}

type StdJsonWriteHandler struct {
}

func (jh *StdJsonWriteHandler) Print(msg string) {
	glog.Errorf(msg)
}
*/
