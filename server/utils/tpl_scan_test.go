package utils

import (
	"strings"
	"testing"
)

// The extractor used to be `\[\[(.*?)\]\]`. Non-greedy, so it stopped at the first `]]` — which
// JavaScript produces itself — and Go's `.` excludes newlines, so multi-line never matched and
// was silently left as literal text. Three documented authoring pitfalls were that one pattern.
func TestScanExpressionsHandlesJavaScriptBrackets(t *testing.T) {
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
		out := scanExpressions(c.in, "[[", "]]", '[', ']', func(inner string) string {
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
func TestScanExpressionsPreservesSurroundingText(t *testing.T) {
	out := scanExpressions(`a=[[one]] b=[[two]] c`, "[[", "]]", '[', ']', func(inner string) string { return "<" + inner + ">" })
	if out != `a=<one> b=<two> c` {
		t.Fatalf("got %q", out)
	}
}

// An unterminated opener must not swallow the rest of the document — the old regex simply did
// not match, and this keeps that behaviour rather than consuming to EOF.
func TestScanExpressionsLeavesUnterminatedOpenerAlone(t *testing.T) {
	in := `before [[ctx.data.x after`
	out := scanExpressions(in, "[[", "]]", '[', ']', func(inner string) string { return "<EXPR>" })
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

// $$...$$ is sugar for a form-value predicate: it wraps the expression in
// function(values, config, target){return …}. That signature is fixed and non-obvious, and the
// sugar is what makes forgetting it unrepresentable — an author who writes the bare expression
// in a [[...]] instead gets a truthy string and a permanently visible element, silently.
//
// It shared the regex extractor's defects: `\$\$(.*?)\$\$` is non-greedy and `.` excludes
// newlines. Now on the same scanner, with no nesting concept since the delimiters are identical.
func TestScanExpressionsHandlesDollarForm(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"plain", `$$values.otpsent$$`, `values.otpsent`},
		{"parenthesised comparison", `$$(values.otpsent!="true")$$`, `(values.otpsent!="true")`},
		{"template literal with ${}", "$$`${values.name} x`$$", "`${values.name} x`"},
		{"multi-line predicate", "$$values.a &&\n  values.b$$", "values.a &&\n  values.b"},
		{"dollars inside a string", `$$values.k === '$$' ? 1 : 2$$`, `values.k === '$$' ? 1 : 2`},
	}
	for _, c := range cases {
		var got []string
		scanExpressions(c.in, "$$", "$$", 0, 0, func(inner string) string {
			got = append(got, inner)
			return "<EXPR>"
		})
		if len(got) != 1 {
			t.Errorf("%s: expected one expression, got %d from %q", c.name, len(got), c.in)
			continue
		}
		if got[0] != c.want {
			t.Errorf("%s:\n  want %q\n  got  %q", c.name, c.want, got[0])
		}
	}
}

// The wrapper is the whole point of the form, so pin it end to end.
func TestProcessTemplateWrapsDollarFormInAFormPredicate(t *testing.T) {
	c := &TestContext{}
	out, err := ProcessTemplate(c, []byte("<Field visible=$$values.a &&\n  values.b$$/>"), nil)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}
	got := string(out)
	if !strings.Contains(got, "function(values, config, target){return") {
		t.Fatalf("expected the form-predicate wrapper, got: %s", got)
	}
	if strings.Contains(got, "$$") {
		t.Fatalf("multi-line $$ expression was not extracted: %s", got)
	}
}

// Both forms must land on the same marker, since jsui replaces one thing.
func TestBothExpressionFormsProduceTheSameMarker(t *testing.T) {
	c := &TestContext{}
	bracket, err := ProcessTemplate(c, []byte(`[[function(values, config, target){return values.x}]]`), nil)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}
	dollar, err := ProcessTemplate(c, []byte(`$$values.x$$`), nil)
	if err != nil {
		t.Fatalf("ProcessTemplate failed: %v", err)
	}
	if string(bracket) != string(dollar) {
		t.Fatalf("the two forms should be interchangeable:\n  [[..]] %s\n  $$..$$ %s", bracket, dollar)
	}
}
