package data

// The canonical logical plan.
//
// One algebraic tree that every query front-end lowers into — declarative dataset filters, OData
// text, openCypher text, and the fluent builder in querybuilder_fluent.go — and that every
// provider compiles from. It exists so the optimizer and the capability decision are written
// ONCE rather than per provider, and so a graph traversal and a relational expansion stop being
// two different things (see PlanExpand).
//
// NAMING. The obvious name for the node interface is taken: LogicalOperator is already this
// package's `and`/`or` string type (query.go:89-90). Nodes are PlanNode and the field on Query is
// Plan, because a plan that shadows an existing type does not compile and renaming the operator
// type would break every provider that switches on it.
//
// WHAT THIS IS NOT, YET. Adding the tree does not move any execution onto it. Query keeps Filter,
// Expand and Navigate as its live fields, every provider keeps compiling them, and Plan is
// populated alongside rather than instead. That is what lets the ~130 plugins pinning this module
// keep working with no edit, and it is why the adapters below are the load-bearing part of this
// file rather than a convenience.

// PlanNode is one node of the canonical logical plan.
//
// Children is the whole traversal contract: an optimizer pass, a capability check and a provider
// lowering all walk the tree through it, so a node that returns nil children is a leaf and a node
// that forgets to return one hides a subtree from every pass at once.
type PlanNode interface {
	// PlanKind identifies the node without a type switch, mirroring Predicate.Kind.
	PlanKind() PlanNodeKind
	// Children returns the node's inputs, in evaluation order. Never nil for a non-leaf.
	Children() []PlanNode
}

// PlanNodeKind names a node type. String-valued, like PredicateKind, so an unknown value read
// from a future version is inspectable rather than an opaque integer.
type PlanNodeKind string

const (
	// PlanKindScan reads an entity's records. The only leaf.
	PlanKindScan PlanNodeKind = "scan"
	// PlanKindFilter restricts its input to rows satisfying a predicate.
	PlanKindFilter PlanNodeKind = "filter"
	// PlanKindProject narrows which fields its input carries, and can retarget which entity the
	// plan returns — see PlanProject.Navigate.
	PlanKindProject PlanNodeKind = "project"
	// PlanKindExpand follows a reference. Covers BOTH a graph traversal and a relational
	// expansion; see PlanExpand.
	PlanKindExpand PlanNodeKind = "expand"
	// PlanKindSort orders its input.
	PlanKindSort PlanNodeKind = "sort"
	// PlanKindPage windows its input.
	PlanKindPage PlanNodeKind = "page"
)

// PlanScan reads the records of one entity. It is the only node with no input.
type PlanScan struct {
	// Entity is the qualified object name, "plugin.Entity", as DataComponent.GetObject reports it.
	Entity string
	// Connection optionally pins which dataconnection to read from, for the providers that span
	// more than one. Empty means the component's own.
	Connection string
}

// PlanKind identifies this node as a scan.
func (s *PlanScan) PlanKind() PlanNodeKind { return PlanKindScan }

// Children returns nil: a scan is the tree's leaf.
func (s *PlanScan) Children() []PlanNode { return nil }

// PlanFilter restricts its input to the rows its predicate admits.
//
// Predicate is the EXISTING Predicate tree, deliberately. The predicate algebra is already
// complete, already lowered from three front-ends, and already compiled by seven providers;
// rebuilding it as plan nodes would invalidate all of that to gain nothing. The plan layer adds
// structure ABOVE the predicate, not a replacement for it.
type PlanFilter struct {
	Input     PlanNode
	Predicate Predicate
}

// PlanKind identifies this node as a filter.
func (f *PlanFilter) PlanKind() PlanNodeKind { return PlanKindFilter }

// Children returns the filtered input.
func (f *PlanFilter) Children() []PlanNode { return []PlanNode{f.Input} }

// PlanProject narrows the fields its input carries, and optionally retargets the plan.
type PlanProject struct {
	Input PlanNode
	// Fields is the property set to keep. Empty means every field.
	Fields []string
	// Navigate retargets what the plan RETURNS: instead of the scanned entity, the entities
	// reached by following these segments. This is where Query.Navigate lowers to, and where
	// Cypher's `RETURN <var>` lands.
	//
	// It is a projection rather than its own node because it changes which rows come back and
	// nothing else — the same thing a projection does, one level out.
	Navigate []string
}

// PlanKind identifies this node as a projection.
func (p *PlanProject) PlanKind() PlanNodeKind { return PlanKindProject }

// Children returns the projected input.
func (p *PlanProject) Children() []PlanNode { return []PlanNode{p.Input} }

// PlanExpand follows a reference from each row of its input.
//
// THIS IS THE UNIFICATION, and it is the reason the plan layer is worth having. Today the platform
// expresses the same idea twice: a Traversal is a PREDICATE (it decides which parent rows match,
// and lives on Query.Filter), while an Expansion is a PROJECTION (it attaches related records to
// rows that already matched, and lives on Query.Expand). Cypher lowers to the first, OData
// $expand to the second, and every provider implements both separately.
//
// One node covers both because the difference is not in the traversal, it is in what the caller
// does with the result — which Quantifier records:
//
//   - Quantifier empty: a projection. The reached records are attached. This is $expand.
//   - QuantifierAny/QuantifierAll: a predicate. The parent row is kept or dropped according to
//     whether the reached records satisfy Filter. This is a Cypher traversal.
//
// The distinction still matters for feasibility and must not be flattened away: a projection can
// be executed after the rows return, by the hop executor, while a quantified expansion decides
// which parents match and so has nothing to attach afterwards. That is why a traversal is
// native-or-nothing and an expansion is not.
type PlanExpand struct {
	Input PlanNode
	// Reference is the field being followed, as declared on the entity.
	Reference string
	// Relationship optionally names the edge type, for providers backing traversal with an edge
	// collection rather than an embedded reference.
	Relationship string
	// MinDepth and MaxDepth bound a variable-length walk. Both zero means a single hop.
	MinDepth int
	MaxDepth int
	// Filter constrains the REACHED entity, not the parent.
	Filter Predicate
	// Project narrows the fields the reached entity carries.
	Project []string
	// Quantifier decides whether this is a predicate or a projection. See the type comment.
	Quantifier Quantifier
	// Optional keeps parent rows whose reference reaches nothing, rather than dropping them.
	Optional bool
	// TargetScope carries the tenancy requirement the reached entity must satisfy. Non-nil when
	// the target is multitenant — dropping it is a cross-tenant read, so it travels with the
	// expansion rather than being reapplied by each provider.
	TargetScope *ScopeRequirement
	// Nested are expansions beneath this one.
	Nested []PlanExpand
	// Top, Skip and Count carry $expand's per-expansion options.
	Top   int
	Skip  int
	Count bool
}

// PlanKind identifies this node as an expansion.
func (e *PlanExpand) PlanKind() PlanNodeKind { return PlanKindExpand }

// Children returns the expanded input.
func (e *PlanExpand) Children() []PlanNode { return []PlanNode{e.Input} }

// IsPredicate reports whether this expansion decides which parent rows match, rather than
// attaching records to rows that already matched.
//
// The two are executed differently and only one can fall back to the hop executor, so callers ask
// this rather than re-deriving it from Quantifier and getting it inconsistently right.
func (e *PlanExpand) IsPredicate() bool {
	return e.Quantifier == QuantifierAny || e.Quantifier == QuantifierAll
}

// PlanSort orders its input.
type PlanSort struct {
	Input PlanNode
	// Fields and Descending are parallel: Descending[i] applies to Fields[i]. A short or absent
	// Descending means ascending throughout.
	//
	// Per-field direction is representable HERE even though some providers cannot execute it —
	// bolthold and badgerhold apply Reverse() to a whole query, so a mixed-direction sort is
	// refused by those. Representing it and refusing it is the honest shape; flattening it into a
	// single flag at the IR would silently answer a different question than the caller asked.
	Fields     []string
	Descending []bool
}

// PlanKind identifies this node as a sort.
func (s *PlanSort) PlanKind() PlanNodeKind { return PlanKindSort }

// Children returns the sorted input.
func (s *PlanSort) Children() []PlanNode { return []PlanNode{s.Input} }

// DescendingFor reports whether field index i sorts descending, tolerating a short Descending.
func (s *PlanSort) DescendingFor(i int) bool {
	if i < 0 || i >= len(s.Descending) {
		return false
	}
	return s.Descending[i]
}

// PlanPage windows its input.
//
// Skip and Limit are ZERO-BASED and absolute, not the platform's 1-based pageNum — the conversion
// belongs at the call site that has the page number, and doing it here would make the node's
// meaning depend on who built it.
type PlanPage struct {
	Input PlanNode
	Skip  int
	// Limit of zero or less means unbounded.
	Limit int
}

// PlanKind identifies this node as a page.
func (p *PlanPage) PlanKind() PlanNodeKind { return PlanKindPage }

// Children returns the paged input.
func (p *PlanPage) Children() []PlanNode { return []PlanNode{p.Input} }

// WalkPlan visits every node depth-first, root first, stopping early if visit returns false.
//
// Every pass over the tree goes through this rather than writing its own recursion, so a new node
// type is reachable by all of them as soon as it returns its children.
func WalkPlan(root PlanNode, visit func(PlanNode) bool) {
	if root == nil || visit == nil {
		return
	}
	if !visit(root) {
		return
	}
	for _, child := range root.Children() {
		if child != nil {
			WalkPlan(child, visit)
		}
	}
}

// FindScan returns the plan's leaf scan, and whether there was exactly one.
//
// A plan with no scan is malformed; a plan with several is a join, which this version does not
// build — so a caller that needs "the entity this plan reads" gets a definite answer or an
// explicit false, rather than the first scan it happens to encounter.
func FindScan(root PlanNode) (*PlanScan, bool) {
	var found *PlanScan
	count := 0
	WalkPlan(root, func(n PlanNode) bool {
		if s, ok := n.(*PlanScan); ok {
			found = s
			count++
		}
		return true
	})
	return found, count == 1
}
