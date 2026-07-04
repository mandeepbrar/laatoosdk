package data

import (
	"encoding/json"
	"testing"
)

// TestStorableRefUnmarshalJSON verifies that StorableRef.UnmarshalJSON accepts both
// the plain string form and the nested object form.
func TestStorableRefUnmarshalJSON(t *testing.T) {
	t.Run("plain string ID", func(t *testing.T) {
		var ref StorableRef
		if err := json.Unmarshal([]byte(`"id123"`), &ref); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Id != "id123" {
			t.Errorf("expected Id=%q, got %q", "id123", ref.Id)
		}
	})

	t.Run("nested object form", func(t *testing.T) {
		var ref StorableRef
		if err := json.Unmarshal([]byte(`{"Id":"id456","Type":"myplugin.Customer","Name":"Acme"}`), &ref); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ref.Id != "id456" {
			t.Errorf("expected Id=%q, got %q", "id456", ref.Id)
		}
		if ref.Type != "myplugin.Customer" {
			t.Errorf("expected Type=%q, got %q", "myplugin.Customer", ref.Type)
		}
		if ref.Name != "Acme" {
			t.Errorf("expected Name=%q, got %q", "Acme", ref.Name)
		}
	})

	t.Run("null input", func(t *testing.T) {
		ref := StorableRef{Id: "existing"}
		if err := json.Unmarshal([]byte(`null`), &ref); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// null must not overwrite existing fields
		if ref.Id != "existing" {
			t.Errorf("null should leave Id intact, got %q", ref.Id)
		}
	})

	t.Run("round-trip plain string", func(t *testing.T) {
		var ref StorableRef
		if err := json.Unmarshal([]byte(`"id789"`), &ref); err != nil {
			t.Fatalf("unmarshal error: %v", err)
		}
		marshaled, err := json.Marshal(ref)
		if err != nil {
			t.Fatalf("marshal error: %v", err)
		}
		var ref2 StorableRef
		if err := json.Unmarshal(marshaled, &ref2); err != nil {
			t.Fatalf("second unmarshal error: %v", err)
		}
		if ref2.Id != "id789" {
			t.Errorf("round-trip: expected Id=%q, got %q", "id789", ref2.Id)
		}
	})

	t.Run("empty bytes no panic", func(t *testing.T) {
		var ref StorableRef
		// json.Unmarshal never passes empty bytes, but guard against it anyway
		if err := ref.UnmarshalJSON([]byte{}); err != nil {
			t.Fatalf("empty bytes should not error: %v", err)
		}
	})
}
