package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"laatoo/sdk/ctx"
	"laatoo/sdk/log"
	"strings"
	"text/template"
)

func ProcessTemplate(ctx ctx.Context, cont []byte, funcs map[string]interface{}) ([]byte, error) {
	contextVar := func(args ...string) string {
		val, ok := ctx.Get(args[0])
		if ok {
			strval, ok := val.(string)
			if ok {
				if len(args) > 1 {
					return fmt.Sprint(args[1], val)
				}
				return strval
			}
			retval, err := json.Marshal(val)
			if err != nil {
				log.Error(ctx, "Error in conf", "Err", err)
			}
			return string(retval)
		}
		return args[0]
	}

	defaultVar := func(args ...string) string {
		_, ok := ctx.Get(args[0])
		if !ok {
			return contextVar(args[1])
		} else {
			return contextVar(args[0])
		}
	}

	eval := func(args ...string) string {
		if len(args) > 1 && args[1] == "insert" {
			return fmt.Sprintf("'+%s+'", args[0])
		}
		return fmt.Sprintf("%s", args[0])
	}

	funcMap := template.FuncMap{"var": contextVar, "eval": eval, "default": defaultVar, "upper": strings.ToUpper, "lower": strings.ToLower, "Title": strings.Title}
	if funcs != nil {
		for k, v := range funcs {
			funcMap[k] = v.(func(variable string) string)
		}
	}
	temp, err := template.New("temp").Funcs(funcMap).Parse(string(cont))
	if err != nil {
		return nil, err
	}
	result := new(bytes.Buffer)
	anon := struct{}{}
	err = temp.Execute(result, anon)
	if err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}

func GetTemplateFileContent(ctx ctx.Context, name string, funcs map[string]interface{}) ([]byte, error) {
	fileData, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	cont, err := ProcessTemplate(ctx, fileData, funcs)
	log.Trace(ctx, "Loaded file", "Name", name, "Conf", string(cont))
	return cont, err
}
