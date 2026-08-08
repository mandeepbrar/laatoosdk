package utils

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log/slog"
	"strings"
	"text/template"

	genutils "laatoo.io/sdk/utils"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/server/log"
)

// scanExpressions finds each open ... close expression and replaces it with transform(inner).
//
// It replaces regexps of the form `\[\[(.*?)\]\]` and `\$\$(.*?)\$\$`, both of which were wrong
// in the same two ways:
//
//   - Non-greedy, so they stopped at the first closing delimiter. For [[...]] that delimiter is
//     `]]`, which JavaScript produces itself: `[[ctx?.data?.list?.[0]]]` captured
//     `ctx?.data?.list?.[0` and `[[a ? ['x','y'] : []]]` captured `a ? ['x','y'] : [`. Both emit
//     truncated JavaScript and leave stray brackets in attribute position.
//   - Go's `.` excludes newlines, so a multi-line expression never matched at all and was left
//     in the output as literal text. Silently — no error anywhere.
//
// This scans instead. String literals are skipped, so a closing delimiter inside quotes cannot
// end an expression, and newlines are ordinary characters.
//
// nestOpen and nestClose track balanced pairs the payload may contain. For [[...]] they are '['
// and ']', so the terminating `]]` is the one at depth zero — which matters because the same
// characters serve both roles: in `[[a?.[0]]]` the inner index and the terminator share brackets.
// For $$...$$ there is no nesting concept, the delimiters being identical, so both are 0 and the
// first unquoted `$$` terminates.
//
// An unterminated opener is emitted literally and scanning continues after it, which keeps the
// old behaviour of leaving malformed input alone rather than consuming to EOF.
func scanExpressions(c, open, close string, nestOpen, nestClose byte, transform func(inner string) string) string {
	var out strings.Builder
	i := 0
	for {
		rel := strings.Index(c[i:], open)
		if rel < 0 {
			out.WriteString(c[i:])
			return out.String()
		}
		start := i + rel
		out.WriteString(c[i:start])

		depth := 0
		var quote byte
		end := -1
		for j := start + len(open); j < len(c); j++ {
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
				continue
			}
			if nestOpen != 0 && ch == nestOpen {
				depth++
				continue
			}
			// The terminator only counts at depth zero. Checked before the nesting close is
			// consumed, because the two overlap: in `[[a?.[0]]]` the same `]` characters both
			// close the inner index and form the terminator.
			if depth == 0 && strings.HasPrefix(c[j:], close) {
				end = j
				break
			}
			if nestClose != 0 && ch == nestClose && depth > 0 {
				depth--
				continue
			}
		}

		if end < 0 {
			out.WriteString(open)
			i = start + len(open)
			continue
		}
		out.WriteString(transform(c[start+len(open) : end]))
		i = end + len(close)
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
	c = scanExpressions(c, "[[", "]]", '[', ']', func(inner string) string {
		b.Reset()
		xml.EscapeText(wr, []byte(inner))
		return fmt.Sprintf("\"javascript###replace@@@%s###\"", b.String())
	})

	c = scanExpressions(c, "$$", "$$", 0, 0, func(inner string) string {
		b.Reset()
		xml.EscapeText(wr, []byte(inner))
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
