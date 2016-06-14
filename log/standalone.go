// +build !appengine

package log

import (
	"laatoo/sdk/core"
	"os"
)

func NewLogger() LoggerInterface {
	return &LogWrapper{logger: NewSimpleLogger(stdSimpleLogsHandler()), level: TRACE}
}

func stdSimpleLogsHandler() SimpleWriteHandler {
	wh := &StdSimpleWriteHandler{}
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
