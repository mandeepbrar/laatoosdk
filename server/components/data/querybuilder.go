package data

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// QueryBuilder builds a query and runs it against one entity, in a chain.
//
// It exists because the AST's expressive constructs were otherwise reachable only by assembling
// struct literals and then hand-writing the compile/bind/execute sequence at every call site. That
// is four steps repeated everywhere, and the repetition is where the mistakes live: a caller who
// compiles without binding gets a query with unresolved parameters, and one who reaches for Get
// with a raw *Query gets a type mismatch from the provider rather than an explanation.
//
//	dm := ctx.GetServerElement(core.ServerElementDataManager).(elements.DataManager)
//	items, _, _, _, err := dm.CreateQuery(ctx, "myplugin.Course").
//	        Fields("Id", "Name").
//	        Where(&data.Comparison{Field: "Status", Operator: data.OpEqual,
//	                               Value: data.ParameterOperand("status")}).
//	        Bind(utils.StringsMap{"status": "active"}).
//	        Page(50, 1).
//	        All()
//
// THIS IS AN INTERFACE AND THE PLATFORM'S IMPLEMENTATION IS NOT IN THIS PACKAGE. It lives in
// laatooserver, and that placement is forced rather than chosen: a terminal has to resolve two
// constructs no single component can. An Expand the provider declines to compile falls back to
// reading each referenced record, and a Navigate never travels to a provider at all — both read a
// DIFFERENT entity's component, while a DataComponent is bound to exactly one and every execution
// method on it returns storables of that one object. So whatever owns the read has to sit above
// the component, and the DataManager is what sits there.
//
// Obtain one from elements.DataManager.CreateQuery, or from CreateTextQuery when you hold query
// text rather than predicates. There is deliberately no constructor here and no CreateQuery on
// DataComponent: a builder reachable from a bare component could not resolve a hop, and would
// have to refuse constructs this one accepts — two builders with one name, differing in what they
// silently cannot do.
//
// A builder is a PER-REQUEST value: it captures the request context, so it carries the caller's
// identity, tenant and soft-delete scope through to bind time. Do not store one on a service
// struct or reuse it across requests — it would carry the first caller's scope into every later
// one, which is the failure the compile/bind split exists to prevent. For a fixed-shape query run
// many times, reach the component through GetRegisteredComponent and use CompileQuery once with
// BindQuery per request; this is the convenience path, not a replacement for that one.
//
// Errors are deferred. Any step may fail, and rather than making every link in the chain return a
// pair, the first error is held and returned by whichever terminal call ends the chain. A chain
// that has already failed does no further work.
type QueryBuilder interface {
	// Fail records an error on the chain, so a producer that cannot even build the query returns
	// a builder rather than a nil — the caller still chains, and finds out at the terminal.
	// Returning nil instead would panic at the next link, far from the cause.
	Fail(err error) QueryBuilder

	// Fields limits which properties are read, as Get's props argument does. No call reads them all.
	Fields(fields ...string) QueryBuilder

	// Where adds a predicate, combining with anything already there under and — so calling it
	// more than once narrows the query rather than replacing what came before.
	Where(predicate Predicate) QueryBuilder

	// Through filters by a predicate on a RELATED entity, whose field names resolve on the entity
	// the path reaches. It returns parents filtered by their children; Expanding attaches the
	// children, and NavigatingTo returns them instead.
	Through(path []string, quantifier Quantifier, predicate Predicate) QueryBuilder

	// Expanding attaches related records to each record returned. Expansions accumulate.
	//
	// A provider that declares ExpandingComponent compiles this into its own store query — a SQL
	// join, a mongo $lookup. One that does not, or that declines this particular projection, is
	// not a failure: the implementation resolves the references after the rows arrive, reading
	// each through the TARGET's own component so its tenancy and soft-delete apply.
	Expanding(expansions ...Expansion) QueryBuilder

	// NavigatingTo returns records of a related entity instead of this entity's own. It is always
	// resolved above the component, never sent to a provider, because it changes the result ENTITY
	// and a component is bound to one.
	NavigatingTo(segments ...string) QueryBuilder

	// Bind supplies the parameters the query's ParameterOperands name. Called more than once, the
	// maps merge, with later keys winning.
	Bind(params utils.StringsMap) QueryBuilder

	// Page sets the page size and 1-based page number. Leave it unset for every record.
	Page(size, number int) QueryBuilder

	// OrderBy orders the result as FIELD/DIRECTION PAIRS — OrderBy("Name", "asc") — matching the
	// orderBy the Get family already takes, which providers walk two at a time.
	OrderBy(fieldsAndDirections ...string) QueryBuilder

	// As selects a DTO projection by name, as Get's dao argument does.
	As(dao string) QueryBuilder

	// Mode sets the provider-specific retrieval mode Get accepts. Rarely needed.
	Mode(mode string) QueryBuilder

	// Query returns the query built so far, for a caller that wants to inspect or keep it — to
	// compile it once with CompileQuery, for instance. It is the live query, not a copy.
	Query() *Query

	// Err reports the first error the chain hit, for a caller that wants to check before a terminal.
	Err() error

	// Condition compiles and binds the query, returning the provider-native condition the Get
	// family takes. Use it when the chain's terminals are not the shape you need — Delete, for
	// example, takes a condition and has no builder terminal here.
	//
	// A condition is what a PROVIDER executes, so it cannot carry a hop the provider does not
	// compile. Expanding and NavigatingTo are therefore resolved by the terminals below and not by
	// this method; asking for a condition on a chain carrying an unresolvable hop is refused by
	// name rather than answered with a condition that quietly drops it.
	Condition() (interface{}, error)

	// All runs the query and returns the page, matching Get's result shape: the records, their
	// ids, the total matching the filter, and how many this page carried.
	//
	// With NavigatingTo the result is the records REACHED, not the ones matched, so ids is nil and
	// both counts describe what is returned — the matched set's ids would pair a child with a
	// parent's id, and its counts would report the parents while the caller holds the children.
	All() (records []core.Storable, ids []string, totalrecs int, recsreturned int, err error)

	// One runs the query and returns a single record.
	One() (core.Storable, error)

	// Count returns how many records match the filter, without reading them.
	//
	// It ignores Fields, Page and OrderBy, because none of them changes how many records match — a
	// count that honoured Page would return the page size and be wrong in a way that looks right.
	Count() (int, error)

	// CountGroups counts matching records grouped by a field, which is the only grouping the
	// platform has. It counts and nothing else: there is no sum, average or general aggregation in
	// this AST, deliberately, so a caller needing one reads the records and aggregates in Go.
	CountGroups(groupids []string, group string) (utils.StringMap, error)
}
