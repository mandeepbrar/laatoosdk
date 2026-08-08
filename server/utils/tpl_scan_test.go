package utils

import (
	"strings"
	"testing"
)

// The extractor used to be `\[\[(.*?)\]\]`. Non-greedy, so it stopped at the first `]]` — which
// JavaScript produces itself — and Go's `.` excludes newlines, so multi-line never matched and
// was silently left as literal text. Three documented authoring pitfalls were that one pattern.
func TestScanBracketExpressionsHandlesJavaScriptBrackets(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain", "[[ctx.data.name]]", "ctx.data.name"},
		{"optional chain array index", "[[ctx?.data?.list?.[0]]]", "ctx?.data?.list?.[0]"},
		{"array literals in a ternary", "[[a ? ['x','y'] : []]]", "a ? ['x','y'] : []"},
		{"nested indexing", "[[a[b[0]]]]", "a[b[0]]"},
		{"multi-line object", "[[{\n  agent: 'x',\n  show: true\n}]]", "{\n  agent: 'x',\n  show: true\n}"},
		{"closing brackets inside a string", `[[a ? ']]' : b]]`, `a ? ']]' : b`},
		{"escaped quote inside a string", `[[a ? 'it\'s ]]' : b]]`, `a ? 'it\'s ]]' : b`},
	}
	for _, c := range cases {
		var got []string
		out := scanBracketExpressions(c.in, func(inner string) string {
			got = append(got, inner)
			return "<EXPR>"
		})
		if len(got) != 1 {
			t.Errorf("%s: expected exactly one expression, got %d from %q (out=%q)", c.name, len(got), c.in, out)
			continue
		}
		if got[0] != c.want {
			t.Errorf("%s:\n  want %q\n  got  %q", c.name, c.want, got[0])
		}
	}
}

// Text outside expressions must survive untouched, and several expressions in one string must
// each be found.
func TestScanBracketExpressionsPreservesSurroundingText(t *testing.T) {
	out := scanBracketExpressions(`a=[[one]] b=[[two]] c`, func(inner string) string { return "<" + inner + ">" })
	if out != `a=<one> b=<two> c` {
		t.Fatalf("got %q", out)
	}
}

// An unterminated opener must not swallow the rest of the document — the old regex simply did
// not match, and this keeps that behaviour rather than consuming to EOF.
func TestScanBracketExpressionsLeavesUnterminatedOpenerAlone(t *testing.T) {
	in := `before [[ctx.data.x after`
	out := scanBracketExpressions(in, func(inner string) string { return "<EXPR>" })
	if out != in {
		t.Fatalf("unterminated opener should be emitted literally; got %q", out)
	}
}

// End to end through ProcessTemplate: the cases that previously produced truncated JavaScript
// now produce a complete marker.
func TestProcessTemplateEmitsCompleteMarkersForBracketHeavyExpressions(t *testing.T) {
	c := &TestContext{}
	for _, in := range []string{
		`<Div>[[ctx?.data?.list?.[0]]]</Div>`,
		`<Div>[[a ? ['x','y'] : []]]</Div>`,
	} {
		out, err := ProcessTemplate(c, []byte(in), nil)
		if err != nil {
			t.Fatalf("ProcessTemplate(%q) failed: %v", in, err)
		}
		got := string(out)
		if !strings.Contains(got, "javascript###replace@@@") {
			t.Fatalf("expected a marker for %q, got %q", in, got)
		}
		if strings.Contains(got, "[[") || strings.Contains(got, "]]") {
			t.Fatalf("expression was truncated, leaving brackets behind: %q", got)
		}
	}
}
