package data

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// QueryBuilder constructs a query and runs it against one entity.
//
// It exists because the AST's expressive constructs were otherwise reachable only by assembling
// struct literals and then hand-writing the compile/bind/execute sequence at every call site. That
// is four steps repeated everywhere, and the repetition is where the mistakes live: a caller who
// compiles without binding gets a query with unresolved parameters, and one who reaches for Get
// with a raw *Query gets a type mismatch from the provider rather than an explanation.
//
//	dm := ctx.GetServerElement(core.ServerElementDataManager).(elements.DataManager)
//	b, err := dm.CreateQuery(srvCtx, "myplugin.Course")
//	if err != nil { return err }
//	b, err = b.Fields(srvCtx, "Id", "Name")
//	if err != nil { return err }
//	b, err = b.Where(srvCtx, &Comparison{Field: "Status", Operator: OpEqual,
//	                                     Value: ParameterOperand("status")})
//	if err != nil { return err }
//	items, _, _, _, err := b.All(reqCtx, utils.StringsMap{"status": "active"})
//
// IT IS A CONSTRUCTOR, NOT A FLUENT CHAIN, and the doc says so because the name does not.
// Construction returns an error, so the calls are sequential statements with checks between them
// and nothing chains off them. That is the platform's convention applied rather than excepted:
// "anything that could log, call another context method, or fail takes a context and returns an
// error", and once construction validates, it can fail.
//
// The deferred-error channel this replaced -- Fail and Err, removed 2026-09-10 -- put the builder
// in the one position no component may occupy: OFF THE PLATFORM'S ERROR PATH. Fail was a bare
// assignment with no wrap, no log and no code, and it could not have been otherwise, having no
// context to wrap with. A chain that failed and was then abandoned before its terminal discarded
// the error entirely. Applying the convention SUBTRACTED two methods rather than adding machinery,
// which is the evidence it belongs here.
//
// TWO CONTEXTS, AT TWO LIFECYCLES, and this is the load-bearing part of the shape.
//
//   - CONSTRUCTION takes a core.ServerContext. Building a query is not a request: its shape is
//     fixed by the code that writes it. So a query may be built once in Initialize and kept.
//   - EXECUTION takes a core.RequestContext. That is what carries the caller's identity, tenant
//     and soft-delete scope, and it arrives per call.
//
// This is what makes a PREPARED query expressible: hold the constructed builder on a service
// struct, and hand each request's own context to the terminal. The older warning -- do not store a
// builder, it captures the request context -- described the shape before this split and is the
// opposite of the rule now.
//
// THIS IS AN INTERFACE AND THE PLATFORM'S IMPLEMENTATION IS NOT IN THIS PACKAGE. It lives in
// laatooserver, and that placement is forced rather than chosen: a terminal has to resolve
// constructs no single component can. An Expand the provider declines to compile falls back to
// reading each referenced record, and navigation returns a DIFFERENT entity's records, while a
// DataComponent is bound to exactly one. So whatever owns the read sits above the component, and
// the DataManager is what sits there.
//
// Obtain one from elements.DataManager.CreateQuery, or from CreateTextQuery when you hold query
// text rather than predicates. There is deliberately no constructor here and no CreateQuery on
// DataComponent: a builder reachable from a bare component could not resolve a hop, and would have
// to refuse constructs this one accepts -- two builders with one name, differing in what they
// silently cannot do.
type QueryBuilder interface {
	// Fields limits which properties are read, as Get's props argument does. No call reads them all.
	Fields(ctx core.ServerContext, fields ...string) (QueryBuilder, error)

	// Where adds a predicate, combining with anything already there under and -- so calling it
	// more than once narrows the query rather than replacing what came before.
	Where(ctx core.ServerContext, predicate Predicate) (QueryBuilder, error)

	// Through filters by a predicate on a RELATED entity, whose field names resolve on the entity
	// the path reaches. It returns parents filtered by their children; Expanding attaches the
	// children, and NavigatingTo returns them instead.
	Through(ctx core.ServerContext, path []string, quantifier Quantifier, predicate Predicate) (QueryBuilder, error)

	// Expanding attaches related records to each record returned. Expansions accumulate.
	//
	// A provider that declares ExpandingComponent compiles this into its own store query -- a SQL
	// join, a mongo $lookup. One that does not, or that declines this particular projection, is
	// not a failure: the implementation resolves the references after the rows arrive, reading
	// each through the TARGET's own component so its tenancy and soft-delete apply.
	Expanding(ctx core.ServerContext, expansions ...Expansion) (QueryBuilder, error)

	// NavigatingTo returns records of a related entity instead of this entity's own, following the
	// named reference fields. It replaces rather than appends: a query has one result type.
	//
	// The segments are REFERENCE FIELD NAMES -- "Course", not the relationshipname stamped on the
	// reference and not a collection. What each reaches is a per-row fact on StorableRef.Type, so
	// nothing here names a target entity.
	//
	// This is the DECLARATION. NavigatingComponent.NavigatingTo is the execution, and the two share
	// a name deliberately: one vocabulary word for one concept. The difference is that this one is
	// lazy and that one performs reads.
	NavigatingTo(ctx core.ServerContext, segments ...string) (QueryBuilder, error)

	// Page sets the page size and 1-based page number. Leave it unset for every record.
	//
	// Under NavigatingTo it pages the records REACHED rather than the ones matched, when the
	// provider serves the navigation natively. See NavigatingComponent.
	Page(ctx core.ServerContext, size, number int) (QueryBuilder, error)

	// OrderBy orders the result as FIELD/DIRECTION PAIRS -- OrderBy(ctx, "Name", "asc") -- matching
	// the orderBy the Get family already takes, which providers walk two at a time.
	OrderBy(ctx core.ServerContext, fieldsAndDirections ...string) (QueryBuilder, error)

	// As selects a DTO projection by name, as Get's dao argument does.
	As(ctx core.ServerContext, dao string) (QueryBuilder, error)

	// Mode sets the provider-specific retrieval mode Get accepts. Rarely needed.
	Mode(ctx core.ServerContext, mode string) (QueryBuilder, error)

	// Query returns the query built so far, for a caller that wants to inspect or keep it. It is
	// the live query, not a copy.
	//
	// It does NOT carry navigation: navigation is a projection and never lived on the AST. Navigate
	// reports what the builder holds.
	Query() *Query

	// Navigate reports the navigation segments this builder was given, or nil.
	Navigate() []string

	// Prepare compiles the query against its provider ONCE, returning a builder that will bind and
	// execute the compiled form per request instead of compiling at every terminal.
	//
	// This is the same compile/bind split DataComponent already exposes -- CompileQuery at a
	// ServerContext, BindQuery at a RequestContext -- surfaced where a caller can reach it. A
	// dataset IS a prepared query: its shape is fixed, so planning is paid once. A Go chain whose
	// shape varies per request simply does not call this, and the terminals compile as before.
	//
	// Call it last: a construction call after Prepare invalidates the compiled form and is refused
	// rather than silently recompiling, because a caller who prepared expects the cost to be paid.
	Prepare(ctx core.ServerContext) (QueryBuilder, error)

	// Condition compiles and binds the query, returning the provider-native condition the Get
	// family takes. Use it when the terminals are not the shape you need -- Delete, for example,
	// takes a condition and has no terminal here.
	//
	// A condition is what a PROVIDER executes, so it cannot carry a hop the provider does not
	// compile, and it owns no read to resolve one afterwards. Expanding is therefore resolved by
	// the terminals below; a chain carrying navigation is refused by name, because a condition
	// executes against one component and navigation changes which entity the result holds.
	Condition(ctx core.RequestContext, params utils.StringsMap) (interface{}, error)

	// All runs the query and returns the page, matching Get's result shape: the records, their
	// ids, the total matching the filter, and how many this page carried.
	//
	// params supplies what the query's ParameterOperands name; parameters belong to the caller and
	// not to the query's shape, which is why they arrive here and not at construction.
	//
	// With NavigatingTo the result is the records REACHED, not the ones matched, so ids is nil and
	// both counts describe what is returned -- the matched set's ids would pair a child with a
	// parent's id, and its counts would report the parents while the caller holds the children.
	All(ctx core.RequestContext, params utils.StringsMap) (records []core.Storable, ids []string, totalrecs int, recsreturned int, err error)

	// One runs the query and returns a single record. Under NavigatingTo it returns the first
	// record REACHED, which is the only reading that does not contradict the query as written.
	One(ctx core.RequestContext, params utils.StringsMap) (core.Storable, error)

	// Count returns how many records match the filter, without reading them.
	//
	// It ignores Fields, Page and OrderBy, because none of them changes how many records match -- a
	// count that honoured Page would return the page size and be wrong in a way that looks right.
	// It refuses a chain carrying navigation: the count of what a navigation reaches is a different
	// question from the count of what the filter matched, and answering either silently would be
	// wrong half the time.
	Count(ctx core.RequestContext, params utils.StringsMap) (int, error)

	// CountGroups counts matching records grouped by a field, which is the only grouping the
	// platform has. It counts and nothing else: there is no sum, average or general aggregation in
	// this AST, deliberately, so a caller needing one reads the records and aggregates in Go.
	CountGroups(ctx core.RequestContext, params utils.StringsMap, groupids []string, group string) (utils.StringMap, error)
}
