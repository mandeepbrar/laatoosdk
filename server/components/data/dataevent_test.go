package data

// Specifies DataEvent's own wire contract: it must survive an encoding/json round trip with
// Changes intact at every nesting depth, since that is the whole reason it exists -- the two
// shapes it replaces could not make that guarantee. This mirrors the platform server's own
// TestWireRoundTrip_SurvivesJSON, scoped to the SDK package boundary: this test only proves the
// TYPE round-trips through JSON correctly, not that a specific server's wire-reconstruction path
// (WireToMessage, in laatooserver/src/core/messagingmanager.go) converts a decoded
// map[string]interface{} to utils.StringMap -- that guarantee is the server's own, verified there.

import (
	"encoding/json"
	"testing"

	"laatoo.io/sdk/utils"
)

// TestDataEvent_RoundTripsWithNestedChanges proves a DataEvent whose Changes carries a nested
// StringMap (the "updated" shape: utils.StringMap{"id":..., "type":..., "data": utils.StringMap{...}})
// survives marshal then unmarshal with every field intact, including through the nesting.
func TestDataEvent_RoundTripsWithNestedChanges(t *testing.T) {
	original := DataEvent{
		Id:        "e-1",
		Entity:    "studentmgmt.Enrollment",
		Operation: string(EventDataUpdated),
		Changes: utils.StringMap{
			"id":   "e-1",
			"type": "studentmgmt.Enrollment",
			"data": utils.StringMap{"Status": "enrolled"},
		},
	}

	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	var got DataEvent
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Id != original.Id {
		t.Errorf("Id = %q, want %q", got.Id, original.Id)
	}
	if got.Entity != original.Entity {
		t.Errorf("Entity = %q, want %q", got.Entity, original.Entity)
	}
	if got.Operation != original.Operation {
		t.Errorf("Operation = %q, want %q", got.Operation, original.Operation)
	}

	// encoding/json decodes an interface{} field's object value as map[string]interface{}
	// regardless of what concrete type was marshalled -- this is expected and is exactly why the
	// platform server's own WireToMessage carries a conversion step this test does not: it proves
	// the FIELDS survive, not that they come back as utils.StringMap specifically.
	changes, ok := got.Changes.(map[string]interface{})
	if !ok {
		t.Fatalf("Changes is %T, want map[string]interface{} (encoding/json's decode shape for an object)", got.Changes)
	}
	if changes["id"] != "e-1" || changes["type"] != "studentmgmt.Enrollment" {
		t.Errorf("Changes top level = %v, want id/type preserved", changes)
	}
	nested, ok := changes["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Changes[\"data\"] is %T, want map[string]interface{}", changes["data"])
	}
	if nested["Status"] != "enrolled" {
		t.Errorf("nested Changes[\"data\"] = %v, want Status=enrolled preserved", nested)
	}
}

// TestDataEvent_RoundTripsFlatCreatePayload proves the "created" shape -- Changes as a flat
// entity-field map with no wrapper -- survives identically, distinguishing this from the nested
// "updated" shape the previous test covers. A regression that only handled one of the two shapes
// (e.g. assuming Changes is always wrapped) would fail exactly this case.
func TestDataEvent_RoundTripsFlatCreatePayload(t *testing.T) {
	original := DataEvent{
		Id:        "c-1",
		Entity:    "studentmgmt.Course",
		Operation: string(EventDataCreated),
		Changes: utils.StringMap{
			"Id": "c-1", "Code": "CS104", "Title": "OS", "Credits": 4,
		},
	}

	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got DataEvent
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	changes, ok := got.Changes.(map[string]interface{})
	if !ok {
		t.Fatalf("Changes is %T, want map[string]interface{}", got.Changes)
	}
	if changes["Code"] != "CS104" || changes["Title"] != "OS" {
		t.Errorf("Changes = %v, want the flat create payload preserved with no wrapper", changes)
	}
	// The flat create shape must not be confused with the wrapped update shape: no "data" key was
	// ever set here, and one must not appear after the round trip.
	if _, hasDataKey := changes["data"]; hasDataKey {
		t.Errorf("Changes contains a \"data\" key that was never set -- flat and wrapped shapes are being conflated")
	}
}

// TestDataEvent_NilChangesOnDelete proves the "deleted" shape -- Changes explicitly nil, since
// nothing changed, the entity is simply gone -- round-trips as nil rather than an empty map or a
// JSON null that decodes to something a consumer would mistake for "no changes were sent" versus
// "changes were sent but empty".
func TestDataEvent_NilChangesOnDelete(t *testing.T) {
	original := DataEvent{Id: "d-1", Entity: "studentmgmt.Enrollment", Operation: string(EventDataDeleted), Changes: nil}

	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var got DataEvent
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Changes != nil {
		t.Errorf("Changes = %v, want nil for a delete event", got.Changes)
	}
	if got.Id != "d-1" || got.Entity != "studentmgmt.Enrollment" || got.Operation != string(EventDataDeleted) {
		t.Errorf("got = %+v, want Id/Entity/Operation preserved alongside nil Changes", got)
	}
}
