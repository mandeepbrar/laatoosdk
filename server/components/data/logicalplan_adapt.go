package data

// Adapters between the legacy Query fields and the canonical plan.
//
// These are what make the plan additive. Query keeps Filter, Expand and Navigate as its live
// fields and every provider keeps compiling them; the plan is built FROM those on demand. Nothing
// is required to migrate, and a caller that never mentions a plan is unaffected.
//
// The round trip is asymmetric on purpose, and the asymmetry is the honest part:
//
//   - PlanFromQuery is total. Every Query expressible today lowers to a plan.
//   - QueryFromPlan is PARTIAL, and says so by returning false. A plan the fluent builder or a
//     future front-end can express — a mixed-direction sort, an expansion nested under a filter —
//     has no Query encoding, and inventing one would silently change what the query asks.
//
// A provider on the legacy path therefore keeps receiving exactly what it receives today, and a
// plan that cannot be degraded to that path is refused rather than approximated.

// PlanFromQuery builds the canonical plan for a query, without modifying it.
//
// Node order is scan → expand* → filter → project → sort → page, which is the order the platform
// already evaluates in: expansions resolve before the filter that may reference them, Navigate
// retargets after filtering, and paging is last. An optimizer is free to reorder; this is the
// unoptimized starting shape.
func PlanFromQuery(q *Query) PlanNode {
	if q == nil {
		return nil
	}

	// the entity is not carried on Query — it comes from the component that holds it — so the
	// scan is left unnamed here and filled in by the caller that knows. See PlanForEntity.
	var node PlanNode = &PlanScan{}

	// expansions first: a filter may reference what they reach
	for i := range q.Expand {
		node = expansionToPlan(node, q.Expand[i])
	}

	if q.Filter != nil {
		// a Traversal on the filter is a QUANTIFIED expansion, not an ordinary predicate — it
		// decides which parent rows match by following a reference. Lifting it into a PlanExpand
		// is what makes it the same node as an $expand, which is the whole point of the IR.
		if trav, ok := q.Filter.(*Traversal); ok {
			node = traversalToPlan(node, trav)
		} else {
			node = &PlanFilter{Input: node, Predicate: q.Filter}
		}
	}

	if len(q.Navigate) > 0 {
		node = &PlanProject{Input: node, Navigate: q.Navigate}
	}

	return node
}

// PlanForEntity is PlanFromQuery with the scan's entity filled in.
//
// Query does not carry the entity — the DataComponent holding it does — so a caller building a
// plan for a provider supplies it here rather than leaving a scan that names nothing.
func PlanForEntity(q *Query, entity string) PlanNode {
	root := PlanFromQuery(q)
	if scan, ok := FindScan(root); ok {
		scan.Entity = entity
	}
	return root
}

// expansionToPlan lowers one Expansion, and its nested expansions, into a projection-shaped
// PlanExpand — no Quantifier, because $expand attaches records rather than deciding which parent
// rows match.
func expansionToPlan(input PlanNode, e Expansion) PlanNode {
	return &PlanExpand{
		Input:     input,
		Reference: e.Field,
		// Levels is $levels: a depth bound on the same reference, so it lowers to MaxDepth.
		// LevelsMax (-1) means unbounded and passes through as-is.
		MaxDepth: e.Levels,
		Filter:   e.Filter,
		Project:  e.Select,
		// Required is $expand's "drop parents that reach nothing"; Optional is its inverse, which
		// is the sense the plan states because the permissive case is the default.
		Optional: !e.Required,
		Nested:   nestedExpansions(e.Expand),
		Top:      e.Top,
		Skip:     e.Skip,
		Count:    e.Count,
	}
}

// nestedExpansions lowers a nested $expand list into plan expansions.
//
// Nested expansions carry no Input of their own: they hang off their parent's Nested field, since
// they are evaluated against the entity the parent reached rather than against the plan's rows.
func nestedExpansions(list []Expansion) []PlanExpand {
	if len(list) == 0 {
		return nil
	}
	out := make([]PlanExpand, 0, len(list))
	for _, e := range list {
		out = append(out, PlanExpand{
			Reference: e.Field,
			MaxDepth:  e.Levels,
			Filter:    e.Filter,
			Project:   e.Select,
			Optional:  !e.Required,
			Nested:    nestedExpansions(e.Expand),
			Top:       e.Top,
			Skip:      e.Skip,
			Count:     e.Count,
		})
	}
	return out
}

// traversalToPlan lowers a Traversal into a quantified PlanExpand.
//
// A Traversal's Path may cross several references; each segment becomes its own expansion so the
// plan states the hops explicitly. The predicate and the quantifier attach to the LAST segment,
// because that is the entity they constrain.
func traversalToPlan(input PlanNode, t *Traversal) PlanNode {
	if t == nil {
		return input
	}
	path := t.Path
	if len(path) == 0 {
		// a traversal naming no path still carries a relationship for edge-backed providers
		path = []string{""}
	}
	node := input
	for i, segment := range path {
		last := i == len(path)-1
		expand := &PlanExpand{
			Input:        node,
			Reference:    segment,
			Relationship: t.Relationship,
			Optional:     t.MatchOptional,
		}
		if last {
			// depth, predicate, quantifier and scope constrain the entity the path REACHES, so
			// they belong on the final hop and nowhere else
			expand.MinDepth = t.MinDepth
			expand.MaxDepth = t.MaxDepth
			expand.Filter = t.Predicate
			expand.Quantifier = t.Quantifier
			expand.TargetScope = t.TargetScope
		}
		node = expand
	}
	return node
}

// QueryFromPlan degrades a plan back to a Query, reporting whether it could be expressed.
//
// False means the plan carries something the legacy fields cannot hold, and the caller must use
// the plan directly or refuse. It is never an approximation: a partial encoding would be a query
// asking a different question than the plan, which is worse than a refusal because it returns
// rows.
func QueryFromPlan(root PlanNode) (*Query, bool) {
	if root == nil {
		return nil, false
	}
	q := NewQuery()
	var expansions []Expansion
	ok := true
	// walk from the root down; each node contributes to the query it can, and any node that
	// cannot be encoded clears ok rather than being skipped silently
	node := root
	for node != nil {
		switch n := node.(type) {
		case *PlanScan:
			node = nil
			continue
		case *PlanFilter:
			if q.Filter != nil {
				// two filters would have to be conjoined, and doing that silently changes an
				// optimizer's deliberate placement into a single flattened predicate
				ok = false
			}
			q.Filter = n.Predicate
			node = n.Input
		case *PlanProject:
			if len(n.Navigate) > 0 {
				q.Navigate = n.Navigate
			}
			// Fields has no home on Query — projection travels as the separate props argument on
			// Get — so a plan carrying one cannot be fully encoded
			if len(n.Fields) > 0 {
				ok = false
			}
			node = n.Input
		case *PlanExpand:
			if n.IsPredicate() {
				// a quantified expansion is a Traversal on the filter, and only when nothing else
				// already occupies it
				if q.Filter != nil {
					ok = false
				} else {
					q.Filter = planExpandToTraversal(n)
				}
			} else {
				expansions = append(expansions, planExpandToExpansion(n))
			}
			node = n.Input
		case *PlanSort, *PlanPage:
			// sorting and paging travel as arguments to Get, not as query fields, so a plan
			// carrying them has no Query encoding for that part
			ok = false
			node = firstChild(n)
		default:
			ok = false
			node = firstChild(n)
		}
	}
	// expansions were collected walking down, so the outermost was seen first; Query lists them
	// in application order
	for i, j := 0, len(expansions)-1; i < j; i, j = i+1, j-1 {
		expansions[i], expansions[j] = expansions[j], expansions[i]
	}
	q.Expand = expansions
	return q, ok
}

// firstChild returns a node's first input, or nil for a leaf.
func firstChild(n PlanNode) PlanNode {
	children := n.Children()
	if len(children) == 0 {
		return nil
	}
	return children[0]
}

// planExpandToTraversal encodes a quantified expansion back as a Traversal.
func planExpandToTraversal(e *PlanExpand) *Traversal {
	return &Traversal{
		Path:          []string{e.Reference},
		Quantifier:    e.Quantifier,
		Predicate:     e.Filter,
		MinDepth:      e.MinDepth,
		MaxDepth:      e.MaxDepth,
		Relationship:  e.Relationship,
		TargetScope:   e.TargetScope,
		MatchOptional: e.Optional,
	}
}

// planExpandToExpansion encodes a projection expansion back as an Expansion.
func planExpandToExpansion(e *PlanExpand) Expansion {
	exp := Expansion{
		Field:    e.Reference,
		Levels:   e.MaxDepth,
		Required: !e.Optional,
		Filter:   e.Filter,
		Select:   e.Project,
		Top:      e.Top,
		Skip:     e.Skip,
		Count:    e.Count,
	}
	for _, nested := range e.Nested {
		n := nested
		exp.Expand = append(exp.Expand, planExpandToExpansion(&n))
	}
	return exp
}
