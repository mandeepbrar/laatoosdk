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

func NewSysLogger() LoggerInterface {
	logWriter, err := syslog.Dial("", "", syslog.LOG_ERR, "Laatoo")
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

func (jh *StdSimpleWriteHandler) Print(ctx core.Context, msg string, level int) {
	os.Stderr.WriteString(msg)
}
func (jh *StdSimpleWriteHandler) PrintBytes(ctx core.Context, msg []byte) (int, error) {
	return os.Stderr.Write(msg)
}

type SyslogWriteHandler struct {
	writer io.Writer
}

func (jh *SyslogWriteHandler) Print(ctx core.Context, msg string, level int) {
	jh.writer.Write([]byte(msg))
}
func (jh *SyslogWriteHandler) PrintBytes(ctx core.Context, msg []byte) (int, error) {
	return jh.writer.Write(msg)
}
