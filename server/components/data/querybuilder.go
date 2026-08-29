package data

import (
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/errors"
	"laatoo.io/sdk/utils"
)

// QueryBuilder builds a query and runs it against one data component, in a chain.
//
// It exists because the AST's expressive constructs were otherwise reachable only by assembling
// struct literals and then hand-writing the compile/bind/execute sequence at every call site. That
// is four steps repeated everywhere, and the repetition is where the mistakes live: a caller who
// compiles without binding gets a query with unresolved parameters, and one who reaches for Get
// with a raw *Query gets a type mismatch from the provider rather than an explanation.
//
//	items, _, _, _, err := dc.CreateQuery(ctx).
//	        Fields("Id", "Name").
//	        Where(&data.Comparison{Field: "Status", Operator: data.OpEqual,
//	                               Value: data.ParameterOperand("status")}).
//	        Bind(utils.StringsMap{"status": "active"}).
//	        Page(50, 1).
//	        All()
//
// The component and the request context are captured at CreateQuery, so a builder is a
// PER-REQUEST value: it carries the caller's identity, tenant and soft-delete scope through to
// bind time. Do not store one on a service struct or reuse it across requests — it would carry the
// first caller's scope into every later one, which is the failure the compile/bind split exists to
// prevent. For a fixed-shape query run many times, keep using CompileQuery once and BindQuery per
// request; this is the convenience path, not a replacement for that one.
//
// Errors are deferred. Any step may fail, and rather than making every link in the chain return a
// pair, the first error is held and returned by whichever terminal call ends the chain. A chain
// that has already failed does no further work.
type QueryBuilder struct {
	component DataComponent
	ctx       core.RequestContext
	query     *Query
	params    utils.StringsMap
	props     []string
	orderBy   []string
	pageSize  int
	pageNum   int
	mode      string
	dao       string
	err       error
}

// NewQueryBuilder starts a chain against a component. Prefer DataComponent.CreateQuery, which
// calls this; the constructor is exported for a caller holding a component it did not get by
// injection, and for the BaseComponent default that implements the interface method.
//
// pageNum starts at 1 rather than 0, matching the Get family: page zero would compute a negative
// skip, which MongoDB rejects outright and other stores answer inconsistently.
func NewQueryBuilder(ctx core.RequestContext, component DataComponent) *QueryBuilder {
	return &QueryBuilder{
		component: component,
		ctx:       ctx,
		query:     NewQuery(),
		pageSize:  -1, // -1/-1 is the Get family's "every record" form
		pageNum:   -1,
	}
}

// Fail records an error on the chain, so an implementor that cannot even build the query returns a
// builder rather than a nil — the caller still chains, and finds out at the terminal. Returning nil
// instead would panic at the next link, far from the cause.
func (b *QueryBuilder) Fail(err error) *QueryBuilder {
	return b.fail(err)
}

// fail records the first error and makes every later step a no-op.
func (b *QueryBuilder) fail(err error) *QueryBuilder {
	if b.err == nil {
		b.err = err
	}
	return b
}

// Fields limits which properties are read, as Get's props argument does. No call reads them all.
func (b *QueryBuilder) Fields(fields ...string) *QueryBuilder {
	b.props = fields
	return b
}

// Where adds a predicate, combining with anything already there under and — so calling it more
// than once narrows the query rather than replacing what came before.
func (b *QueryBuilder) Where(predicate Predicate) *QueryBuilder {
	b.query.Where(predicate)
	return b
}

// Through filters by a predicate on a RELATED entity, whose field names resolve on the entity the
// path reaches. It returns parents filtered by their children; Expanding attaches the children,
// and NavigatingTo returns them instead.
func (b *QueryBuilder) Through(path []string, quantifier Quantifier, predicate Predicate) *QueryBuilder {
	b.query.Through(path, quantifier, predicate)
	return b
}

// Expanding attaches related records to each record returned. Expansions accumulate.
func (b *QueryBuilder) Expanding(expansions ...Expansion) *QueryBuilder {
	b.query.Expanding(expansions...)
	return b
}

// NavigatingTo returns records of a related entity instead of this component's own. It is resolved
// above the component, by the server, since a component is bound to a single entity.
func (b *QueryBuilder) NavigatingTo(segments ...string) *QueryBuilder {
	b.query.NavigatingTo(segments...)
	return b
}

// Bind supplies the parameters the query's ParameterOperands name. Called more than once, the maps
// merge, with later keys winning.
func (b *QueryBuilder) Bind(params utils.StringsMap) *QueryBuilder {
	if b.params == nil {
		b.params = utils.StringsMap{}
	}
	for name, value := range params {
		b.params[name] = value
	}
	return b
}

// Page sets the page size and 1-based page number. Leave it unset for every record.
func (b *QueryBuilder) Page(size, number int) *QueryBuilder {
	b.pageSize, b.pageNum = size, number
	return b
}

// OrderBy orders the result as FIELD/DIRECTION PAIRS — OrderBy("Name", "asc") — matching the
// orderBy the Get family already takes, which providers walk two at a time.
func (b *QueryBuilder) OrderBy(fieldsAndDirections ...string) *QueryBuilder {
	b.orderBy = fieldsAndDirections
	return b
}

// As selects a DTO projection by name, as Get's dao argument does.
func (b *QueryBuilder) As(dao string) *QueryBuilder {
	b.dao = dao
	return b
}

// Mode sets the provider-specific retrieval mode Get accepts. Rarely needed.
func (b *QueryBuilder) Mode(mode string) *QueryBuilder {
	b.mode = mode
	return b
}

// Query returns the query built so far, for a caller that wants to inspect or keep it — to compile
// it once with CompileQuery, for instance. It is the live query, not a copy.
func (b *QueryBuilder) Query() *Query {
	return b.query
}

// Err reports the first error the chain hit, for a caller that wants to check before a terminal.
func (b *QueryBuilder) Err() error {
	return b.err
}

// Condition compiles and binds the query, returning the provider-native condition the Get family
// takes. Use it when the chain's terminals are not the shape you need — Delete, for example, takes
// a condition and has no builder terminal here.
func (b *QueryBuilder) Condition() (interface{}, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.component == nil {
		// a builder with no component came from NewQueryBuilder(ctx, nil); say so here
		// rather than panicking on the nil interface at the first terminal
		return nil, errors.BadArg(b.ctx, "component")
	}
	// one call rather than CompileQuery + BindQuery: a builder's query is per-request by
	// construction, so there is no compiled form to reuse
	return b.component.CreateQueryCondition(b.ctx, b.query, b.params)
}

// All runs the query and returns the page, matching Get's result shape: the records, their ids,
// the total matching the filter, and how many this page carried.
func (b *QueryBuilder) All() (records []core.Storable, ids []string, totalrecs int, recsreturned int, err error) {
	condition, err := b.Condition()
	if err != nil {
		return nil, nil, 0, 0, err
	}
	return b.component.Get(b.ctx, b.props, condition, b.pageSize, b.pageNum, b.mode, b.orderBy, b.dao)
}

// One runs the query and returns a single record.
func (b *QueryBuilder) One() (core.Storable, error) {
	condition, err := b.Condition()
	if err != nil {
		return nil, err
	}
	return b.component.GetOne(b.ctx, b.props, condition, b.dao)
}

// Count returns how many records match the filter, without reading them.
//
// It ignores Fields, Page and OrderBy, because none of them changes how many records match — a
// count that honoured Page would return the page size and be wrong in a way that looks right.
func (b *QueryBuilder) Count() (int, error) {
	condition, err := b.Condition()
	if err != nil {
		return 0, err
	}
	return b.component.Count(b.ctx, condition)
}

// CountGroups counts matching records grouped by a field, which is the only grouping the platform
// has. It counts and nothing else: there is no sum, average or general aggregation in this AST,
// deliberately, so a caller needing one reads the records and aggregates in Go.
func (b *QueryBuilder) CountGroups(groupids []string, group string) (utils.StringMap, error) {
	condition, err := b.Condition()
	if err != nil {
		return nil, err
	}
	return b.component.CountGroups(b.ctx, condition, groupids, group)
}
