// +build !appengine

package log

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"laatoo/sdk/core"
	"strings"
	"time"
)

func init() {
	logFormats["happycolor"] = printHappyColor
	logFormats["happymaxcolor"] = printHappyMaxColor
}

func printHappyMaxColor(ctx core.Context, app string, strlevel string, wh SimpleWriteHandler, level int, msg string, args ...interface{}) {
	if len(args)%2 > 0 {
		panic("wrong logging")
	}
	var buffer bytes.Buffer
	firstline := ""
	if level <= ERROR {
		firstline = color.RedString("%s: %s", strings.ToUpper(strlevel), msg)
	} else if level == INFO {
		firstline = color.BlueString("%s: %s", strings.ToUpper(strlevel), msg)
	} else {
		firstline = color.GreenString("%s: %s", strings.ToUpper(strlevel), msg)
	}
	argslen := len(args)
	if argslen > 0 {
		firstline = fmt.Sprintf("%s    %s", firstline, color.MagentaString("%s:%s", strings.ToUpper(args[0].(string)), fmt.Sprint(args[1])))
	}
	if argslen > 2 {
		firstline = fmt.Sprintf("%s    %s", firstline, color.CyanString("%s:%s", strings.ToUpper(args[2].(string)), fmt.Sprint(args[3])))
	}
	buffer.WriteString(fmt.Sprintln(firstline))
	for i := 4; (i + 1) < argslen; i = i + 2 {
		buffer.WriteString(fmt.Sprintln("		", args[i], ":", args[i+1]))
	}
	buffer.WriteString(fmt.Sprintln("		TIME ", time.Now().String()))
	buffer.WriteString(fmt.Sprintln("		LEVEL ", strlevel))
	buffer.WriteString(fmt.Sprintln("		CONTEXT ", ctx.GetName()))
	if ctx != nil {
		buffer.WriteString(fmt.Sprintln("		PATH ", ctx.GetPath()))
		buffer.WriteString(fmt.Sprintln("		ID ", ctx.GetId()))
	}
	wh.Print(ctx, app, buffer.String(), level, strlevel)
}
func printHappyColor(ctx core.Context, app string, strlevel string, wh SimpleWriteHandler, level int, msg string, args ...interface{}) {
	if len(args)%2 > 0 {
		panic("wrong logging")
	}
	var buffer bytes.Buffer
	firstline := ""
	if level <= ERROR {
		firstline = color.RedString("%s: %s", strings.ToUpper(strlevel), msg)
	} else if level == INFO {
		firstline = color.BlueString("%s: %s", strings.ToUpper(strlevel), msg)
	} else {
		firstline = color.GreenString("%s: %s", strings.ToUpper(strlevel), msg)
	}
	argslen := len(args)
	if argslen > 0 {
		firstline = fmt.Sprintf("%s    %s", firstline, color.MagentaString("%s:%s", strings.ToUpper(args[0].(string)), fmt.Sprint(args[1])))
	}
	if argslen > 2 {
		firstline = fmt.Sprintf("%s    %s", firstline, color.CyanString("%s:%s", strings.ToUpper(args[2].(string)), fmt.Sprint(args[3])))
	}
	buffer.WriteString(fmt.Sprintln(firstline))
	for i := 4; (i + 1) < argslen; i = i + 2 {
		buffer.WriteString(fmt.Sprintln("		", args[i], ":", args[i+1]))
	}
	if ctx != nil {
		buffer.WriteString(fmt.Sprintln("		", ctx.GetName()))
	}
	wh.Print(ctx, app, buffer.String(), level, strlevel)
}
