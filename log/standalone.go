// +build !appengine

package log

import (
	"io"
	"laatoo/sdk/core"
	"log/syslog"
	"os"
)

func NewLogger() LoggerInterface {
	return &LogWrapper{logger: NewSimpleLogger(stdSimpleLogsHandler()), level: TRACE}
}

func NewSysLogger(appname string) LoggerInterface {
	logWriter, err := syslog.Dial("", "", syslog.LOG_ERR, appname)
	if err != nil {
		return NewLogger()
	}
	return &LogWrapper{logger: NewSimpleLogger(sysLogsHandler(logWriter)), level: TRACE}
}

func stdSimpleLogsHandler() SimpleWriteHandler {
	wh := &StdSimpleWriteHandler{}
	return wh
}

func sysLogsHandler(writer io.Writer) SimpleWriteHandler {
	wh := &SyslogWriteHandler{writer}
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

type SyslogWriteHandler struct {
	writer io.Writer
}

func (jh *SyslogWriteHandler) Print(ctx core.Context, appname string, msg string, level int, strlevel string) {
	jh.writer.Write([]byte(msg))
}
func (jh *SyslogWriteHandler) PrintBytes(ctx core.Context, appname string, msg []byte, level int, strlevel string) (int, error) {
	return jh.writer.Write(msg)
}
