package utils

import "regexp"

// [\s\S] rather than [^\n]: the expression may span lines.
//
// While the SDK's extractor was a regex that could not match across newlines either, keeping the
// quotes on a multi-line expression was harmless — it was never converted, so the value stayed
// inert literal text and the document still parsed. Since laatoo.io/sdk v0.1.183 the extractor is
// a scanner and DOES convert it, so surviving quotes now wrap the marker's own quotes and produce
// data=""javascript###...###"", which is not valid XML.
//
// The two changes had to land together: stripping quotes earlier would have left raw brackets in
// attribute position for a multi-line expression the old SDK could not convert, turning a
// parsing-but-inert case into a parse failure.
//
// Bracket depth is handled by backtracking rather than by counting: the trailing quote anchors
// the match, so `"[[ctx?.arr?.[0]]]"` grows the lazy body until the whole pattern fits.
var doubleQuotedTemplateExpressionRegex = regexp.MustCompile(`"\s*(\[\[[\s\S]*?\]\])\s*"`)
var singleQuotedTemplateExpressionRegex = regexp.MustCompile(`'\s*(\[\[[\s\S]*?\]\])\s*'`)

// NormalizeQuotedTemplateExpressions removes wrapping quotes only when the entire
// quoted value is a single [[...]] expression. This keeps registry files valid
// after ProcessTemplate adds its own expression quotes.
//
// It is the first half of the registry-XML pipeline and ProcessTemplate is the second: registry
// XML is not valid XML until both have run, so a consumer that parses the raw file is worse than
// nothing — a bare `&&` in an attribute reads as an error when it is the required form. Both
// halves live here so anything outside the jsui plugin can run the pipeline in the right order.
func NormalizeQuotedTemplateExpressions(content []byte) []byte {
	normalized := doubleQuotedTemplateExpressionRegex.ReplaceAll(content, []byte("$1"))
	return singleQuotedTemplateExpressionRegex.ReplaceAll(normalized, []byte("$1"))
}
