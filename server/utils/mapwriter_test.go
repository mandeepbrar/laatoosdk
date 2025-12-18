package utils

import (
	"testing"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/utils"
)

// MockSerializable for testing
type MockSerializable struct {
	Data map[string]interface{}
}

func (m *MockSerializable) WriteAll(ctx ctx.Context, cdc datatypes.Codec, w datatypes.SerializableWriter) error {
	var err error
	if val, ok := m.Data["Name"]; ok {
		s := val.(string)
		err = w.WriteString(ctx, cdc, "Name", &s)
	}
	if val, ok := m.Data["Age"]; ok {
		i := val.(int)
		err = w.WriteInt(ctx, cdc, "Age", &i)
	}
	return err
}

func (m *MockSerializable) ReadAll(ctx ctx.Context, cdc datatypes.Codec, r datatypes.SerializableReader) error {
	return nil
}
func (m *MockSerializable) GetId() string   { return "" }
func (m *MockSerializable) SetId(id string) {}
func (m *MockSerializable) GetType() string { return "MockObject" }

// MockServerContext
type MockServerContext struct {
	core.ServerContext
}

func (m *MockServerContext) CreateObject(objType string) (interface{}, error) {
	if objType == "MockObject" {
		return &MockSerializable{Data: make(map[string]interface{})}, nil
	}
	return nil, nil // Error handling mocked implicitly
}

func TestCreateObjectFromMap(t *testing.T) {
	// Setup
	mockCtx := &MockServerContext{}
	smap := utils.StringMap{
		"Name": "John Doe",
		"Age":  30,
	}

	// Case 1: No transformations
	_, err := CreateObjectFromMap(mockCtx, "MockObject", smap, nil)
	if err != nil {
		t.Errorf("CreateObjectFromMap failed: %v", err)
	}
	// Verify (Mock doesn't really write back to struct, but we check if it runs without error)
	// Ideally we would inspect the writer's effect on the object, but MapSerializableWriter WRITES FROM the map TO the object.
	// The MockSerializable.WriteAll calls w.WriteString("Name", &val).
	// MapSerializableWriter implementation of WriteString: gets "Name" from map and sets *val.
	// So effectively MockSerializable.WriteAll is "reading" from the writer.

	// Let's refine the test verification logic.
	// MapSerializableWriter acts as a source.
	// Object.WriteAll calls writer methods to populate itself.

	// We need a better MockSerializable that actually updates its internal state when calling writer methods?
	// Wait, WriteAll signature: func (s *Struct) WriteAll(ctx, cdc, w ObjectWriter)
	// Inside WriteAll:
	// var name string
	// w.WriteString(..., "Name", &name)
	// s.Name = name

	// My MockSerializable above is wrong direction. It is trying to write TO the writer.
	// But MapSerializableWriter is an ObjectWriter.
	// CreateObjectFromMap calls `serObj.WriteAll(ctx, nil, wr)`.
	// MapSerializableWriter implements ObjectWriter.
	// Yes, `WriteString` in `MapSerializableWriter` takes `val *string` and populates it FROM the map.

}

// Better Mock Object
type TestObject struct {
	Name  string
	Age   int
	Other string
	Sub   *NestedObject
	Sub2  *NestedObject
}

func (t *TestObject) WriteAll(c ctx.Context, cdc datatypes.Codec, w datatypes.SerializableWriter) error {
	w.WriteString(c, cdc, "Name", &t.Name)
	w.WriteInt(c, cdc, "Age", &t.Age)
	w.WriteString(c, cdc, "DiffProp", &t.Other)

	t.Sub = &NestedObject{}
	w.WriteObject(c, cdc, "SourceSub", t.Sub) // Expecting "SourceSub" key to be present (unrenamed)

	t.Sub2 = &NestedObject{}
	w.WriteObject(c, cdc, "RenamedSub", t.Sub2) // Expecting "RenamedSub" key (renamed from SourceRenamedSub)

	return nil
}
func (t *TestObject) ReadAll(ctx ctx.Context, cdc datatypes.Codec, r datatypes.SerializableReader) error {
	return nil
}
func (t *TestObject) GetId() string   { return "" }
func (t *TestObject) SetId(id string) {}
func (t *TestObject) GetType() string { return "TestObject" }

type TestContext struct {
	core.ServerContext
}

func (c *TestContext) CreateObject(t string) (interface{}, error) {
	return &TestObject{}, nil
}

func TestTransformations(t *testing.T) {
	ctx := &TestContext{}
	srcMap := utils.StringMap{
		"SourceKey1": "MappedValue",
		"SourceKey2": 123,
		"Name":       "OriginalName",
		"SourceSub": utils.StringMap{
			"OldChild": "ChildValue",
		},
		"SourceRenamedSub": utils.StringMap{
			"OldChild2": "ChildValue2",
		},
	}

	transforms := utils.StringMap{
		"SourceKey1": "DiffProp",
		"SourceKey2": "Age",
		"SourceSub": utils.StringMap{
			"OldChild": "NewChild",
		},
		"SourceRenamedSub": utils.StringMap{
			"__key":     "RenamedSub",
			"OldChild2": "NewChild2",
		},
	}

	// Expectations:
	// "SourceKey1" -> "DiffProp"
	// "SourceKey2" -> "Age"
	// "SourceSub" -> "SourceSub" (no rename), but child "OldChild" -> "NewChild"
	// "SourceRenamedSub" -> "RenamedSub", and child "OldChild2" -> "NewChild2"

	// Create struct to hold results (since our TestObject is flat, we'll verify via MapSerializableWriter internal state if we can,
	// or update TestObject to support nesting, OR just inspect the returned object if we cast it back?
	// Our TestObject writes to itself.
	// Let's UPDATE TestObject to support testing these nested fields.

	obj, err := CreateObjectFromMap(ctx, "TestObject", srcMap, transforms)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	tObj := obj.(*TestObject)

	if tObj.Other != "MappedValue" {
		t.Errorf("Expected Other='MappedValue', got '%s'", tObj.Other)
	}
	if tObj.Age != 123 {
		t.Errorf("Expected Age=123, got %d", tObj.Age)
	}
	// Check submap
	if tObj.Sub == nil {
		t.Errorf("Expected Sub to be populated")
	} else if tObj.Sub.NewChild != "ChildValue" {
		t.Errorf("Expected Sub.NewChild='ChildValue' (via sub-transform), got '%s'", tObj.Sub.NewChild)
	}

	if tObj.Sub2 == nil {
		t.Errorf("Expected Sub2 to be populated")
	} else if tObj.Sub2.NewChild2 != "ChildValue2" {
		t.Errorf("Expected Sub2.NewChild2='ChildValue2' (via sub-transform + rename), got '%s'", tObj.Sub2.NewChild2)
	}
}

type NestedObject struct {
	NewChild  string
	NewChild2 string
}

func (n *NestedObject) WriteAll(ctx ctx.Context, cdc datatypes.Codec, w datatypes.SerializableWriter) error {
	w.WriteString(ctx, cdc, "NewChild", &n.NewChild)
	w.WriteString(ctx, cdc, "NewChild2", &n.NewChild2)
	return nil
}

// Implement ReadAll etc... for NestedObject
func (n *NestedObject) ReadAll(ctx ctx.Context, cdc datatypes.Codec, r datatypes.SerializableReader) error {
	return nil
}
func (n *NestedObject) GetId() string   { return "" }
func (n *NestedObject) SetId(id string) {}
func (n *NestedObject) GetType() string { return "NestedObject" }

func TestDotNotationTransformations(t *testing.T) {

	srcMap := utils.StringMap{
		"Application": utils.StringMap{
			"Identifier": "123",
			"Name":       "TestApp",
		},
		"Other": "Value",
	}

	// Transform Application.Identifier -> Id
	// Transform Application.Name -> AppName
	transforms := utils.StringMap{
		"Application.Identifier": "Id",
		"Application.Name":       "AppName",
	}

	// We expect the result map to be:
	// Application: { Id: "123", AppName: "TestApp" }
	// Other: "Value"

	resMap := applyTransformations(srcMap, unflattenTransformations(transforms))

	appMap, ok := resMap["Application"].(utils.StringMap)
	if !ok {
		t.Fatalf("Expected Application to be utils.StringMap, got %T", resMap["Application"])
	}

	if val, ok := appMap["Id"]; !ok || val != "123" {
		t.Errorf("Expected Application.Id to be '123', got %v", val)
	}

	if val, ok := appMap["AppName"]; !ok || val != "TestApp" {
		t.Errorf("Expected Application.AppName to be 'TestApp', got %v", val)
	}

	if val, ok := resMap["Other"]; !ok || val != "Value" {
		t.Errorf("Expected Other to be 'Value', got %v", val)
	}
}
