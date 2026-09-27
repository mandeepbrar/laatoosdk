package data

import "testing"

// eq builds an equality comparison against a literal, the shape most access rules lower to.
func eq(field string, value interface{}) Predicate {
	return &Comparison{Field: field, Operator: OpEqual, Value: LiteralOperand(value)}
}

// TestEntityAccessFilterZeroValueIsNeverPermission pins that an unset answer refuses.
func TestEntityAccessFilterZeroValueIsNeverPermission(t *testing.T) {
	var unset EntityAccessFilter
	if _, ok := unset.Restrict(eq("A", 1)); ok {
		t.Fatal("a zero-value filter restricted instead of refusing")
	}
	if admitted, _ := unset.Admits(map[string]interface{}{"A": 1}); admitted {
		t.Fatal("a zero-value filter admitted a record")
	}
	var nilFilter *EntityAccessFilter
	if _, ok := nilFilter.Restrict(nil); ok {
		t.Fatal("a nil filter restricted instead of refusing")
	}
	if FilteredEntityAccess(nil).Decision != EntityAccessDenied {
		t.Fatal("a filtered answer with no predicate must be a denial")
	}
}

// TestEntityAccessFilterRestrictDoesNotMutate pins that the caller's filter is left untouched.
func TestEntityAccessFilterRestrictDoesNotMutate(t *testing.T) {
	callerFilter := &Logical{Operator: LogicalAnd, Operands: []Predicate{eq("A", 1), eq("B", 2)}}
	access := FilteredEntityAccess(eq("Owner", "u1"))
	restricted, ok := access.Restrict(callerFilter)
	if !ok {
		t.Fatal("a filtered answer refused")
	}
	if len(callerFilter.Operands) != 2 {
		t.Fatalf("Restrict mutated the caller's filter: %d operands", len(callerFilter.Operands))
	}
	logical, isLogical := restricted.(*Logical)
	if !isLogical || logical.Operator != LogicalAnd || len(logical.Operands) != 2 || logical.Operands[0] != callerFilter {
		t.Fatalf("Restrict did not AND the access predicate onto the caller's filter: %#v", restricted)
	}
	if got, ok := UnrestrictedEntityAccess().Restrict(callerFilter); !ok || got != callerFilter {
		t.Fatal("an unrestricted answer changed the caller's filter")
	}
	if _, ok := DeniedEntityAccess().Restrict(callerFilter); ok {
		t.Fatal("a denied answer restricted instead of refusing")
	}
}

// TestMatchesRecord covers each construct an access rule lowers to, with a negative for each.
func TestMatchesRecord(t *testing.T) {
	record := map[string]interface{}{"Owner": "u1", "Level": 3, "Active": true, "Score": 7.5}
	cases := []struct {
		name string
		pred Predicate
		want bool
	}{
		{"equal", eq("Owner", "u1"), true},
		{"equal miss", eq("Owner", "u2"), false},
		{"int against float literal", eq("Level", 3.0), true},
		{"greater", &Comparison{Field: "Level", Operator: OpGreater, Value: LiteralOperand(2)}, true},
		{"greater miss", &Comparison{Field: "Level", Operator: OpGreater, Value: LiteralOperand(3)}, false},
		{"less equal", &Comparison{Field: "Score", Operator: OpLessEqual, Value: LiteralOperand(7.5)}, true},
		{"bool", eq("Active", true), true},
		{"not equal", &Comparison{Field: "Owner", Operator: OpNotEqual, Value: LiteralOperand("u2")}, true},
		{"and", &Logical{Operator: LogicalAnd, Operands: []Predicate{eq("Owner", "u1"), eq("Active", true)}}, true},
		{"and miss", &Logical{Operator: LogicalAnd, Operands: []Predicate{eq("Owner", "u1"), eq("Active", false)}}, false},
		{"or", &Logical{Operator: LogicalOr, Operands: []Predicate{eq("Owner", "u2"), eq("Active", true)}}, true},
		{"or miss", &Logical{Operator: LogicalOr, Operands: []Predicate{eq("Owner", "u2"), eq("Active", false)}}, false},
		{"not", &Not{Operand: eq("Owner", "u2")}, true},
		{"in", &Membership{Field: "Owner", Values: []Operand{LiteralOperand("u3"), LiteralOperand("u1")}}, true},
		{"in miss", &Membership{Field: "Owner", Values: []Operand{LiteralOperand("u3")}}, false},
		{"not in", &Membership{Field: "Owner", Values: []Operand{LiteralOperand("u3")}, Negated: true}, true},
		{"absent field equals nothing", eq("Missing", "u1"), false},
		{"absent field is null", &NullTest{Field: "Missing"}, true},
		{"present field is not null", &NullTest{Field: "Owner", Negated: true}, true},
		{"absent field is never ordered", &Comparison{Field: "Missing", Operator: OpLess, Value: LiteralOperand(1)}, false},
	}
	for _, tc := range cases {
		got, err := MatchesRecord(tc.pred, record)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestMatchesRecordRefusesWhatItCannotDecide pins that an unsupported shape errors, never guesses.
func TestMatchesRecordRefusesWhatItCannotDecide(t *testing.T) {
	record := map[string]interface{}{"Owner": "u1"}
	refused := []Predicate{
		nil,
		&Comparison{Field: "Owner", Operator: OpEqual, Value: ParameterOperand("p")},
		&Logical{Operator: LogicalAnd},
		&FunctionCall{Function: FuncContains, Field: "Owner"},
	}
	for i, pred := range refused {
		if _, err := MatchesRecord(pred, record); err == nil {
			t.Errorf("case %d: expected a refusal", i)
		}
	}
}
