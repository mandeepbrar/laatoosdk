// +build !appengine

package log

import (
	"laatoo/sdk/core"
	"os"
)

func NewLogger() LoggerInterface {
	return &LogWrapper{logger: NewSimpleLogger(stdSimpleLogsHandler()), level: TRACE}
}

/*
func NewSysLogger(appname string) LoggerInterface {
	logWriter, err := syslog.Dial("", "", syslog.LOG_ERR, appname)
	if err != nil {
		return NewLogger()
	}
	return &LogWrapper{logger: NewSimpleLogger(sysLogsHandler(logWriter)), level: TRACE}
}


func sysLogsHandler(writer io.Writer) SimpleWriteHandler {
	wh := &SyslogWriteHandler{writer}
	return wh
}*/

func stdSimpleLogsHandler() SimpleWriteHandler {
	wh := &StdSimpleWriteHandler{}
	return wh
}

type StdSimpleWriteHandler struct {
}

func (jh *StdSimpleWriteHandler) Print(ctx core.Context, appname string, msg string, level int, strlevel string) {
	os.Stderr.WriteString(msg)
}
func (jh *StdSimpleWriteHandler) PrintBytes(ctx core.Context, appname string, msg []byte, level int, strlevel string) (int, error) {
	return os.Stderr.Write(msg)
}
