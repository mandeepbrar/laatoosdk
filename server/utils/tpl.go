package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"regexp"
	"strings"
	"text/template"

	genutils "laatoo.io/sdk/utils"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/server/log"
)

// scanBracketExpressions finds each [[ ... ]] expression and replaces it with transform(inner).
//
// It replaces a regexp of the form `\[\[(.*?)\]\]`, which was wrong in two ways that between
// them account for three documented authoring pitfalls:
//
//   - the non-greedy match stops at the FIRST `]]`, and JavaScript produces `]]` of its own.
//     `[[ctx?.data?.list?.[0]]]` captured `ctx?.data?.list?.[0` and `[[a ? ['x','y'] : []]]`
//     captured `a ? ['x','y'] : [` — both emit truncated, broken JavaScript.
//   - Go's `.` excludes newlines, so a multi-line expression never matched at all and was left
//     in the output as literal text, silently.
//
// This scans instead: bracket depth decides the terminator, and string literals are skipped so a
// `]]` inside quotes cannot end the expression. Newlines are ordinary characters.
//
// An unterminated `[[` is emitted literally and scanning continues after it, which keeps the old
// behaviour of leaving malformed input alone rather than consuming the rest of the file.
func scanBracketExpressions(c string, transform func(inner string) string) string {
	var out strings.Builder
	i := 0
	for {
		start := strings.Index(c[i:], "[[")
		if start < 0 {
			out.WriteString(c[i:])
			return out.String()
		}
		start += i
		out.WriteString(c[i:start])

		depth := 0
		var quote byte
		end := -1
		for j := start + 2; j < len(c); j++ {
			ch := c[j]
			if quote != 0 {
				if ch == '\\' {
					j++ // skip the escaped character
					continue
				}
				if ch == quote {
					quote = 0
				}
				continue
			}
			switch ch {
			case '\'', '"', '`':
				quote = ch
			case '[':
				depth++
			case ']':
				if depth > 0 {
					depth--
				} else if j+1 < len(c) && c[j+1] == ']' {
					end = j
				}
			}
			if end >= 0 {
				break
			}
		}

		if end < 0 {
			// Unterminated: emit the opener literally and carry on.
			out.WriteString("[[")
			i = start + 2
			continue
		}
		out.WriteString(transform(c[start+2 : end]))
		i = end + 2
	}
}

func ProcessTemplate(ctx ctx.Context, cont []byte, funcs map[string]interface{}) ([]byte, error) {
	contextVar := func(args ...string) string {
		val, ok := ctx.Get(args[0])
		if ok {
			strval, ok := val.(string)
			if ok {
				if len(args) > 1 && strings.TrimSpace(strval) != "" {
					return fmt.Sprint(args[1], val)
				}
				return strval
			}
			if val == nil {
				return ""
			}
			retval, err := json.Marshal(val)
			if err != nil {
				log.Error(ctx, "Error in conf", slog.Any("Err", err))
			}
			return string(retval)
		}
		return ""
	}

	defaultVar := func(args ...string) string {
		_, ok := ctx.Get(args[0])
		if !ok {
			return contextVar(args[1])
		} else {
			return contextVar(args[0])
		}
	}

	is := func(variable string) bool {
		val, _ := ctx.GetBool(variable)
		return val
	}

	exists := func(variable string) bool {
		_, ok := ctx.Get(variable)
		return ok
	}

	contains := func(variable string, val string) bool {
		vals, ok := ctx.GetStringArray(variable)
		if ok {
			return genutils.StrContains(vals, val) >= 0
		}
		return false
	}

	equals := func(variable string, val string) bool {
		valToCompare, ok := ctx.Get(variable)
		if ok {
			return fmt.Sprintf("%v", valToCompare) == val
		}
		return false
	}

	json := func(variable string) string {
		varval, ok := ctx.Get(variable)
		if ok {
			val, _ := json.Marshal(varval)
			if val != nil {
				return string(val)
			}
		}
		return ""
	}

	jsReplace := func(args ...string) string {
		return fmt.Sprintf("javascript###replace@@@%s###", args[0])
	}

	jsFormat := func(args ...string) string {
		vars := ""
		if len(args) > 1 {
			vars = strings.Join(args[1:], "@@@")
		}
		return fmt.Sprintf("javascript###format@@@%s@@@%s###", args[0], vars)
	}

	funcMap := template.FuncMap{"var": contextVar, "is": is, "jsreplace": jsReplace, "jsformat": jsFormat, "default": defaultVar, "equals": equals, "upper": strings.ToUpper, "lower": strings.ToLower, "title": strings.Title, "exists": exists, "contains": contains, "json": json}
	for k, v := range funcs {
		funcMap[k] = v
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

	var b bytes.Buffer
	wr := io.Writer(&b)

	c := result.String()
	c = scanBracketExpressions(c, func(inner string) string {
		b.Reset()
		xml.EscapeText(wr, []byte(inner))
		return fmt.Sprintf("\"javascript###replace@@@%s###\"", b.String())
	})

	re2 := regexp.MustCompile(`\$\$(.*?)\$\$`)
	c = re2.ReplaceAllStringFunc(c, func(inp string) string {
		b.Reset()
		mval := inp[2 : len(inp)-2]
		xml.EscapeText(wr, []byte(mval))
		return fmt.Sprintf("\"javascript###replace@@@function(values, config, target){return %s}###\"", b.String())
	})

	return []byte(c), nil
}

func GetTemplateFileContent(ctx ctx.Context, name string, funcs map[string]interface{}) ([]byte, error) {
	fileData, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	cont, err := ProcessTemplate(ctx, fileData, funcs)
	return cont, err
}
