package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"laatoo/sdk/core"
	"strings"
	"time"
)

func init() {
	logFormats["json"] = printJSON
	logFormats["jsonmax"] = printJSONMax
	logFormats["happy"] = printHappy
	logFormats["happymax"] = printHappyMax
}

func printJSON(ctx core.Context, strlevel string, wh SimpleWriteHandler, level int, msg string, args ...interface{}) {
	if len(args)%2 > 0 {
		panic("wrong logging")
	}
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	mapToPrint := map[string]string{"MESSAGE": msg}
	argslen := len(args)
	for i := 0; (i + 1) < argslen; i = i + 2 {
		mapToPrint[args[i].(string)] = fmt.Sprint(args[i+1])
	}
	err := enc.Encode(mapToPrint)
	if err != nil {
		fmt.Println(err)
	}
	wh.Print(ctx, buffer.String(), level)
}
func printJSONMax(ctx core.Context, strlevel string, wh SimpleWriteHandler, level int, msg string, args ...interface{}) {
	if len(args)%2 > 0 {
		panic("wrong logging")
	}
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	mapToPrint := map[string]string{"TIME": time.Now().String(), "LEVEL": strlevel, "MESSAGE": msg}
	if ctx != nil {
		mapToPrint["CONTEXT"] = ctx.GetName()
		mapToPrint["ID"] = ctx.GetId()
	}
	argslen := len(args)
	for i := 0; (i + 1) < argslen; i = i + 2 {
		mapToPrint[args[i].(string)] = fmt.Sprint(args[i+1])
	}
	err := enc.Encode(mapToPrint)
	if err != nil {
		fmt.Println(err)
	}
	wh.Print(ctx, buffer.String(), level)
}
func printHappy(ctx core.Context, strlevel string, wh SimpleWriteHandler, level int, msg string, args ...interface{}) {
	if len(args)%2 > 0 {
		panic("wrong logging")
	}
	var buffer bytes.Buffer
	firstline := fmt.Sprintf("%s: %s", strings.ToUpper(strlevel), msg)
	argslen := len(args)
	if argslen > 0 {
		firstline = fmt.Sprintf("%s    %s:%s", firstline, strings.ToUpper(args[0].(string)), fmt.Sprint(args[1]))
	}
	if argslen > 2 {
		firstline = fmt.Sprintf("%s    %s:%s", firstline, strings.ToUpper(args[2].(string)), fmt.Sprint(args[3]))
	}
	buffer.WriteString(fmt.Sprintln(firstline))
	for i := 4; (i + 1) < argslen; i = i + 2 {
		buffer.WriteString(fmt.Sprintln("		", args[i], ":", args[i+1]))
	}
	if ctx != nil {
		buffer.WriteString(fmt.Sprintln("		", ctx.GetName()))
	}
	wh.Print(ctx, buffer.String(), level)
}
func printHappyMax(ctx core.Context, strlevel string, wh SimpleWriteHandler, level int, msg string, args ...interface{}) {
	if len(args)%2 > 0 {
		panic("wrong logging")
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintln("MESSAGE ", msg))
	buffer.WriteString(fmt.Sprintln("		TIME ", time.Now().String()))
	buffer.WriteString(fmt.Sprintln("		LEVEL ", strlevel))
	if ctx != nil {
		buffer.WriteString(fmt.Sprintln("		CONTEXT ", ctx.GetName()))
		buffer.WriteString(fmt.Sprintln("		ID ", ctx.GetId()))
	}
	argslen := len(args)
	for i := 0; (i + 1) < argslen; i = i + 2 {
		buffer.WriteString(fmt.Sprintln("		", args[i], " ", args[i+1]))
	}
	wh.Print(ctx, buffer.String(), level)
}
