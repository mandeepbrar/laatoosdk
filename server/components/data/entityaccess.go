package data

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// EntityAction is what a caller is doing to a record, as record-level access rules name it.
//
// The four values are the action column of an entityaccess policy row. A rule is written per
// action, so a caller may read a record it cannot update.
type EntityAction string

const (
	// EntityActionRead covers every read: a list, a by-id fetch, a count, and a record reached
	// through expansion or traversal from another entity.
	EntityActionRead EntityAction = "read"
	// EntityActionCreate is decided against the record being created.
	EntityActionCreate EntityAction = "create"
	// EntityActionUpdate is decided against the stored record AND against the record as it would
	// be after the write, so an update cannot move a record out of the caller's scope.
	EntityActionUpdate EntityAction = "update"
	// EntityActionDelete is decided against the stored record.
	EntityActionDelete EntityAction = "delete"
)

// EntityAccessDecision is the shape of a record-level access answer.
//
// The zero value is EntityAccessInvalid on purpose: a filter nobody set must never read as
// "unrestricted".
type EntityAccessDecision int

const (
	// EntityAccessInvalid is the zero value and is never a real answer. A consumer receiving it
	// must treat it as a failure, not as permission.
	EntityAccessInvalid EntityAccessDecision = iota
	// EntityAccessUnrestricted means no record-level restriction applies to this caller — a
	// system context, or a caller whose applicable rules admit every record (a rule of "true").
	// An opted-in entity with no row for the caller is EntityAccessDenied, never this.
	EntityAccessUnrestricted
	// EntityAccessDenied means no rule admits this caller for this entity and action: a list
	// returns nothing and a single-record operation answers not found.
	EntityAccessDenied
	// EntityAccessFiltered means the caller may reach exactly the records Predicate matches.
	EntityAccessFiltered
)

// EntityAccessFilter is the record-level access a caller holds for one entity and action,
// already reduced to a predicate over the record's own fields.
//
// It is produced by the authorization back end from the entityaccess policy: the rows that apply
// to the caller's roles are selected, the caller's own values are substituted for every r.sub
// reference, and what remains is an expression over r.obj fields — which is this Predicate. Every
// operand is therefore a LITERAL; a filter never carries a parameter, so it can be folded into a
// query at bind time without the caller supplying anything.
//
// A data component applies it beside tenancy: ANDed into a list query so paging and totals
// reflect only admitted records, and as "id == X AND Predicate" for a by-id operation, so a
// record is reachable by id exactly when it would appear in the caller's list.
type EntityAccessFilter struct {
	// Decision says which of the three answers this is.
	Decision EntityAccessDecision
	// Predicate is set only when Decision is EntityAccessFiltered. Field names are the entity's
	// field names as the policy wrote them after "r.obj.".
	Predicate Predicate
}

// UnrestrictedEntityAccess answers that no record-level restriction applies.
func UnrestrictedEntityAccess() *EntityAccessFilter {
	return &EntityAccessFilter{Decision: EntityAccessUnrestricted}
}

// DeniedEntityAccess answers that the caller may reach no record.
func DeniedEntityAccess() *EntityAccessFilter {
	return &EntityAccessFilter{Decision: EntityAccessDenied}
}

// FilteredEntityAccess answers that the caller may reach the records the predicate matches. A nil
// predicate is refused by returning a denial: an empty restriction is never an open one.
func FilteredEntityAccess(predicate Predicate) *EntityAccessFilter {
	if predicate == nil {
		return DeniedEntityAccess()
	}
	return &EntityAccessFilter{Decision: EntityAccessFiltered, Predicate: predicate}
}

// Restrict returns filter ANDed with this access filter's predicate, as a NEW node — the argument
// is never mutated, because a caller's filter is frequently part of a compiled query shared across
// requests. It returns filter unchanged when access is unrestricted, and nil with ok=false when
// access is denied (or the answer is invalid), which the caller must turn into an empty result or
// a not-found rather than an unfiltered read.
func (f *EntityAccessFilter) Restrict(filter Predicate) (restricted Predicate, ok bool) {
	if f == nil {
		return nil, false
	}
	switch f.Decision {
	case EntityAccessUnrestricted:
		return filter, true
	case EntityAccessFiltered:
		if f.Predicate == nil {
			return nil, false
		}
		if filter == nil {
			return f.Predicate, true
		}
		return &Logical{Operator: LogicalAnd, Operands: []Predicate{filter, f.Predicate}}, true
	default:
		return nil, false
	}
}

// Admits reports whether a record, given as field name to value, satisfies this access filter.
//
// It is how a decision is made where there is no stored record to query: the record being
// created, and the record as an update would leave it. It evaluates the same predicate a list
// query compiles, so an in-memory decision and a database decision cannot disagree about a rule.
func (f *EntityAccessFilter) Admits(record map[string]interface{}) (bool, error) {
	if f == nil {
		return false, nil
	}
	switch f.Decision {
	case EntityAccessUnrestricted:
		return true, nil
	case EntityAccessFiltered:
		return MatchesRecord(f.Predicate, record)
	default:
		return false, nil
	}
}

// MatchesRecord evaluates a predicate against a record held in memory.
//
// It supports the constructs record-level access rules lower to — comparisons against literals,
// and, or, not, membership and null tests — and refuses anything else with an error rather than
// guessing. A field absent from the record compares as nil: equal only to nil, and never ordered.
func MatchesRecord(predicate Predicate, record map[string]interface{}) (bool, error) {
	switch node := predicate.(type) {
	case nil:
		return false, fmt.Errorf("data: cannot evaluate a nil predicate")
	case *Comparison:
		if node.Value.Kind != OperandLiteral {
			return false, fmt.Errorf("data: in-memory evaluation supports literal operands only, field %s has a %s operand", node.Field, node.Value.Kind)
		}
		return compareValues(record[node.Field], node.Operator, node.Value.Value)
	case *Logical:
		if len(node.Operands) == 0 {
			return false, fmt.Errorf("data: cannot evaluate an empty %s", node.Operator)
		}
		for _, operand := range node.Operands {
			matched, err := MatchesRecord(operand, record)
			if err != nil {
				return false, err
			}
			if node.Operator == LogicalAnd && !matched {
				return false, nil
			}
			if node.Operator == LogicalOr && matched {
				return true, nil
			}
		}
		return node.Operator == LogicalAnd, nil
	case *Not:
		matched, err := MatchesRecord(node.Operand, record)
		if err != nil {
			return false, err
		}
		return !matched, nil
	case *Membership:
		value := record[node.Field]
		found := false
		for _, operand := range node.Values {
			if operand.Kind != OperandLiteral {
				return false, fmt.Errorf("data: in-memory evaluation supports literal operands only, field %s has a %s operand", node.Field, operand.Kind)
			}
			equal, err := compareValues(value, OpEqual, operand.Value)
			if err != nil {
				return false, err
			}
			if equal {
				found = true
				break
			}
		}
		return found != node.Negated, nil
	case *NullTest:
		isNull := isNilValue(record[node.Field])
		return isNull != node.Negated, nil
	default:
		return false, fmt.Errorf("data: in-memory evaluation does not support %s predicates", predicate.Kind())
	}
}

// compareValues applies a comparison operator to two values, normalising numbers to float64 and
// times to their instant so that an int field and a float literal compare by value.
func compareValues(left interface{}, op CompareOperator, right interface{}) (bool, error) {
	if isNilValue(left) || isNilValue(right) {
		bothNil := isNilValue(left) && isNilValue(right)
		switch op {
		case OpEqual:
			return bothNil, nil
		case OpNotEqual:
			return !bothNil, nil
		default:
			return false, nil
		}
	}
	if lnum, lok := toFloat(left); lok {
		if rnum, rok := toFloat(right); rok {
			return orderResult(op, compareFloat(lnum, rnum))
		}
	}
	if ltime, lok := left.(time.Time); lok {
		if rtime, rok := right.(time.Time); rok {
			return orderResult(op, compareTime(ltime, rtime))
		}
	}
	if lbool, lok := left.(bool); lok {
		if rbool, rok := right.(bool); rok {
			switch op {
			case OpEqual:
				return lbool == rbool, nil
			case OpNotEqual:
				return lbool != rbool, nil
			default:
				return false, fmt.Errorf("data: booleans cannot be ordered with %s", op)
			}
		}
	}
	return orderResult(op, strings.Compare(fmt.Sprint(left), fmt.Sprint(right)))
}

// orderResult turns a three-way comparison into the answer for one operator.
func orderResult(op CompareOperator, cmp int) (bool, error) {
	switch op {
	case OpEqual:
		return cmp == 0, nil
	case OpNotEqual:
		return cmp != 0, nil
	case OpGreater:
		return cmp > 0, nil
	case OpGreaterEqual:
		return cmp >= 0, nil
	case OpLess:
		return cmp < 0, nil
	case OpLessEqual:
		return cmp <= 0, nil
	default:
		return false, fmt.Errorf("data: unknown comparison operator %q", op)
	}
}

// compareFloat is a three-way comparison of two numbers.
func compareFloat(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// compareTime is a three-way comparison of two instants.
func compareTime(a, b time.Time) int {
	switch {
	case a.Before(b):
		return -1
	case a.After(b):
		return 1
	default:
		return 0
	}
}

// toFloat converts any Go numeric kind to float64.
func toFloat(value interface{}) (float64, bool) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), true
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	default:
		return 0, false
	}
}

// isNilValue reports whether a value is nil, including a typed nil pointer.
func isNilValue(value interface{}) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
