package config

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/utils"
)

// Config is a parsed configuration document — the conf a service receives in Initialize, and
// every sub-document reachable from it.
//
// # This is the RAW document, and that is the distinction that catches people
//
// The conf handed to Initialize is the rendered YAML as written, so every key in the file is
// readable here whether or not anything declared it. That is the opposite of
// GetStringConfiguration on a configurable object, which reads only the keys declared in the
// object spec's configurations block and answers ("", false) for anything else. Reading a key
// one way when it was meant to be declared the other way is the common failure, and neither
// reports an error.
//
// # Templating has already happened
//
// Every method takes a ctx, but no method resolves context variables through it: the value is
// returned exactly as parsed. A {{var "x"}} expression reaches a Config only if something
// rendered it earlier, and one that was not rendered is returned as the literal template text.
// The ctx is carried for logging and for symmetry with the rest of the SDK, not for lookup.
//
// # Type conversion is per-accessor, not uniform
//
// GetInt, GetFloat and GetBool are permissive: they convert from the neighbouring scalar types
// and from a string, so `retries: 3` and `retries: "3"` both read through GetInt.
//
// GetString is strict and does NOT convert. It requires the value to already be a string, so
// `port: 8080` — which YAML infers as a number — reads as ABSENT through GetString. Quote a
// value in the YAML when the code reads it as text.
//
// The boolean second return means "present AND convertible to this type". It never distinguishes
// a missing key from a present one of the wrong shape, so a config typo and a type mismatch look
// identical at the call site.
type Config interface {
	// GetString returns the value when it is already a string. A value of any other type — a
	// YAML number, a bool — reports not-found rather than being formatted, so quote a value in
	// the YAML when the code reads it as text.
	GetString(ctx ctx.Context, configurationName string) (string, bool)

	// GetBool returns the value as a bool, converting from a string via strconv.ParseBool.
	GetBool(ctx ctx.Context, configurationName string) (bool, bool)

	// GetInt returns the value as an int, converting from a float64 or a string.
	//
	// The float64 conversion TRUNCATES rather than rounding or refusing, and it matters more
	// than it looks: YAML decodes a bare number into a float64 in several paths, so a value
	// written as 1.9 is read as 1 with no indication anything was lost.
	GetInt(ctx ctx.Context, configurationName string) (int, bool)

	// GetFloat returns the value as a float64, converting from an int or a string.
	GetFloat(ctx ctx.Context, configurationName string) (float64, bool)

	// GetStringArray returns a list whose every element is a string.
	//
	// It is all-or-nothing: one non-string element fails the whole call and returns nil, rather
	// than skipping that element or converting it.
	GetStringArray(ctx ctx.Context, configurationName string) ([]string, bool)

	// GetSubConfig returns a nested document as its own Config.
	//
	// The value must actually BE a nested mapping. This is where a quoting mistake in the YAML
	// surfaces: a templated map written in quotes renders to a string rather than a mapping, and
	// a string is not convertible here, so the call reports not-found and the service falls back
	// to its defaults as though the block were absent. Nothing logs, and the block is plainly
	// there in the file — which is what makes it expensive to find.
	GetSubConfig(ctx ctx.Context, configurationName string) (Config, bool)

	// GetStringsMap returns a nested mapping whose values are all strings.
	//
	// All-or-nothing like GetStringArray: one non-string value fails the whole map. Use
	// GetStringMap when the values are of mixed type.
	GetStringsMap(ctx ctx.Context, configurationName string) (utils.StringsMap, bool)

	// GetStringMap returns a nested mapping with values left as interface{}, for a block whose
	// values are not uniformly typed.
	GetStringMap(ctx ctx.Context, configurationName string) (utils.StringMap, bool)

	// GetConfigArray returns a list of nested documents, for a repeated block.
	GetConfigArray(ctx ctx.Context, configurationName string) ([]Config, bool)

	// GetRoot returns the single top-level entry of the document — its name and its contents.
	//
	// It reports not-found unless there is EXACTLY ONE top-level key: it is for unwrapping a
	// document that names its own root, not for reaching the first of several keys.
	GetRoot(ctx ctx.Context) (string, Config, bool)

	// Get returns the raw value with no conversion, for a value none of the typed accessors fit.
	Get(ctx ctx.Context, configurationName string) (interface{}, bool)

	// SetString stores a string value, replacing any existing one.
	SetString(ctx ctx.Context, configurationName string, configurationValue string)

	// Set stores a value of any type, replacing any existing one. It mutates this Config in
	// place rather than returning a new one — see Clone before setting on a config that was
	// handed to you.
	Set(ctx ctx.Context, configurationName string, configurationValue interface{})

	// SetVals stores every entry of vals. A nil map is a no-op rather than an error.
	SetVals(ctx ctx.Context, vals utils.StringMap)

	// Clone returns a copy that can be modified without affecting this one — but only one level
	// deep in practice.
	//
	// Nested documents are copied by value only when they are already Config-typed; a nested
	// mapping still in its parsed map form is shared by reference, and that is the usual state
	// of a document that came from YAML. Setting a key inside a nested block of a clone can
	// therefore change the original. Where that matters, clone the sub-config too.
	Clone() Config

	// ToMap exposes the underlying map. It is NOT a copy: writing through the returned map
	// mutates the Config, and so does anything the caller passes it to.
	ToMap() map[string]interface{}

	// AllConfigurations returns the top-level key names, in unspecified order. Use it to iterate
	// a block whose keys are data — module names, instance names — rather than a fixed schema.
	AllConfigurations(ctx ctx.Context) []string
}
