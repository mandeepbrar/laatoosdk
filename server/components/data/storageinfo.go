package data

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/log"
	serverutils "laatoo.io/sdk/server/utils"
	"laatoo.io/sdk/utils"
)

/**
protobuf numbers

id = 51
deleted=52
isnew=53
createdby=54
updatedby=55
createdat=56
updatedat=57
type=59
name=60
tenant=61
AbstractStorable=62
SoftDeleteStorable=63
Entity=64
AbstractStorableMT=65
SoftDeleteStorableMT=66
HardDeleteAuditable=67
SoftDeleteAuditable=68
HardDeleteAuditableMT=69
SoftDeleteAuditableMT=70
SerializableBase=71
collection=72
tenantname=73
version=74
relationship=75
annotations=76
*/

type StorageInfo struct {
	Id      string      `json:"Id" bson:"Id" protobuf:"bytes,51,opt,name=id,proto3" sql:"type:varchar(50); primary key;" gorm:"primary_key"`
	Version string      `json:"Version" bson:"Version" protobuf:"bytes,74,opt,name=version,proto3" sql:"type:varchar(50);" `
	selfRef interface{} `json:"-" datastore:"-" bson:"-" sql:"-"`
}

func (si *StorageInfo) Constructor(ctx ctx.Context) {
	if si.Id == "" {
		si.Id = ctx.CreateUUID() //uuid.NewV4().String()
	}
}
func (si *StorageInfo) SetSelfReference(ref interface{}) {
	si.selfRef = ref
}

func (si *StorageInfo) GetId() string {
	return si.Id
}
func (si *StorageInfo) SetId(val string) {
	si.Id = val
}

func (si *StorageInfo) GetLabel() string {
	stor := si.selfRef.(core.Storable)
	c := stor.Config()
	if c != nil && c.LabelField != "" {
		v := reflect.ValueOf(stor).Elem()
		f := v.FieldByName(c.LabelField)
		return f.String()
	}
	return ""
}

func (si *StorageInfo) GetVersion() string {
	return si.Version
}

func (si *StorageInfo) PreSave(ctx ctx.Context) error {
	return nil
}
func (si *StorageInfo) PostSave(ctx ctx.Context) error {
	return nil
}
func (si *StorageInfo) PostLoad(ctx ctx.Context) error {
	return nil
}
func (si *StorageInfo) SetValues(ctx ctx.Context, obj interface{}, val utils.StringMap) error {
	delete(val, "Id")
	delete(val, "IsNew")
	delete(val, "CreatedBy")
	delete(val, "UpdatedBy")
	delete(val, "CreatedAt")
	delete(val, "UpdatedAt")
	return serverutils.SetObjectFields(ctx, obj, val, nil, nil)
}

func (si *StorageInfo) IsMultitenant() bool {
	return false
}

func (si *StorageInfo) Join(item core.Storable) {
}
func (si *StorageInfo) Config() *core.StorableConfig {
	return nil
}

func (si *StorageInfo) GetObjectRef() interface{} {
	stor := si.selfRef.(core.Storable)
	c := stor.Config()
	return &StorableRef{Id: si.Id, Type: c.ObjectType, Name: stor.GetLabel(), Version: stor.GetVersion()}
}

func (si *StorageInfo) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error
	if err = rdr.ReadString(c, cdc, "Id", &si.Id); err != nil {
		return err
	}
	return nil
}

func (si *StorageInfo) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error
	if err = wtr.WriteString(c, cdc, "Id", &si.Id); err != nil {
		return err
	}
	return nil
}

// StorableRef is a typed directed edge: a reference from the record holding it to the storable
// named by Type and Id. Id, Type, Name and Version are all carried per instance, so a ref is
// resolvable without consulting the schema — which is what lets expansion resolve a target at
// runtime with no schema lookup, and what Relationship is carried the same way to preserve.
type StorableRef struct {
	Id      string `json:"Id" bson:"Id" protobuf:"bytes,51,opt,name=id,proto3" sql:"type:varchar(100);"`
	Type    string `json:"Type" bson:"Type" protobuf:"bytes,59,opt,name=type,proto3" sql:"type:varchar(100);"`
	Name    string `json:"Name" bson:"Name" protobuf:"bytes,60,opt,name=name,proto3" sql:"type:varchar(300);"`
	Version string `json:"Version" bson:"Version" protobuf:"bytes,74,opt,name=version,proto3" sql:"type:varchar(50);" `
	// Relationship is the edge TYPE — what a traversal matches on — declared per reference
	// field in the entity YAML under the key `relationshipname` and stamped into the entity's
	// generated WriteAll by codegen, exactly as Type already is. The generated guard is
	// `if ref.Relationship == "" { ref.Relationship = "<declared>" }`, so an empty value means
	// "use the declaration" and an explicit one set by a caller survives.
	//
	// It is a stored value rather than a struct tag on purpose. A tag needs no SDK change at
	// all, but it is invisible to every storage engine: shadow columns and key-prefix indexes
	// index stored VALUES, so only a stored relationship can answer "every KNOWS edge" from an
	// index, and only a stored one survives into a native graph provider.
	//
	// It is scalar and non-list so that the entity `index` flag emits its tags for it, which
	// is the whole reason it is not simply a key in Annotations.
	Relationship string `json:"Relationship" bson:"Relationship" protobuf:"bytes,75,opt,name=relationship,proto3" sql:"type:varchar(100);"`
	// DataConnection names the dataconnection holding the referenced record, forming the other
	// half of the (dataconnection, Type) pair the registry is keyed on. Empty means "the target is
	// registered on exactly one connection, resolve it there", so every reference stored before
	// this field existed resolves exactly as it did -- there is no migration and no backfill.
	//
	// It is needed because Type alone stops identifying a component once the same entity name is
	// registered on two connections, which is precisely what a per-connection plugin creates. The
	// alternative -- resolving to whichever connection happens to hold the name -- would route a
	// read to a store the reference never named, and that fallback is rejected platform-wide.
	//
	// SET IT ONLY WHEN THE TARGET IS AMBIGUOUS, and expect it to be VALIDATED against the registry
	// on save: a value naming a connection the target is not registered on is refused, the way a
	// contradicting Type already is. A per-row copy of a fact the registry owns is a fact that can
	// be wrong, so the registry stays authoritative and this stays a disambiguator.
	//
	// A REFERENCE ACROSS CONNECTIONS IS HOP-EXECUTOR-ONLY, PERMANENTLY. A dataconnection is a
	// store boundary, no provider can compile a join across two of them, so such a reference can
	// never be resolved natively however capable the provider. That is the same constraint
	// traversal already carries for cross-store paths, now reachable through an ordinary
	// storableref rather than only through a graph query -- worth knowing before modelling one,
	// because it is a permanent property of the reference and not a temporary provider gap.
	//
	// Cross-store references themselves are not new: the hop executor has always resolved each
	// entity through its own component, so a mongo record referencing a postgres one already
	// worked. This names WHICH component when the name alone is ambiguous.
	DataConnection string `json:"DataConnection,omitempty" bson:"DataConnection,omitempty" protobuf:"bytes,77,opt,name=dataconnection,proto3" sql:"type:varchar(100);"`
	// Annotations is free-form metadata that TRAVELS with the reference — provenance, the
	// actor who created the link, a display label, a correlation id. It is persisted, and it
	// is deliberately NOT queryable.
	//
	// The line between this and Relationship is the same line the platform draws between an
	// annotation and an edge entity, and it is worth stating because the temptation to filter
	// on an annotation is obvious: anything you want to FILTER on is an edge entity, anything
	// you merely want to CARRY is an annotation. A map is unindexable on every provider — the
	// index flag emits tags for scalar, non-list fields only — and its values compare as
	// strings, so a range query over one would be wrong as well as slow. A query naming an
	// annotation key is refused rather than answered.
	//
	// A nil map costs nothing: WriteAll skips it entirely, so an unused Annotations adds no
	// bytes to any stored reference.
	Annotations utils.StringsMap `json:"Annotations,omitempty" bson:"Annotations,omitempty" protobuf:"bytes,76,rep,name=annotations,proto3" sql:"-"`
	Entity      core.Storable    `json:"-" datastore:"-" bson:"-" sql:"-" firestore:"-" protobuf:"group,64,opt,name=Entity,proto3"`
}

func (si *StorableRef) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error
	if err = rdr.ReadString(c, cdc, "Id", &si.Id); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "Type", &si.Type); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "Name", &si.Name); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "Version", &si.Version); err != nil {
		return err
	}
	// both new fields read clean from data written before they existed: ReadString and ReadMap
	// each leave the target untouched when the property is absent, so an old reference decodes
	// to an empty Relationship and a nil Annotations rather than to an error. That is what
	// makes this addition need no migration.
	if err = rdr.ReadString(c, cdc, "Relationship", &si.Relationship); err != nil {
		return err
	}
	// same tolerance as Relationship above: a reference written before this field existed reads
	// back empty, which resolves through the registry when the target is on one connection
	if err = rdr.ReadString(c, cdc, "DataConnection", &si.DataConnection); err != nil {
		return err
	}
	if err = rdr.ReadMap(c, cdc, "Annotations", &si.Annotations); err != nil {
		return err
	}
	return nil
}

func (si *StorableRef) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error
	if err = wtr.WriteString(c, cdc, "Id", &si.Id); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "Type", &si.Type); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "Name", &si.Name); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "Version", &si.Version); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "Relationship", &si.Relationship); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "DataConnection", &si.DataConnection); err != nil {
		return err
	}
	// guarded, unlike every other field here, because a StorableRef is embedded in a large
	// share of the platform's stored records: WriteMap is handed a non-nil POINTER to a nil
	// map, so it would emit the key regardless and add an empty Annotations to every reference
	// ever written. Skipping it keeps an unused annotation map genuinely free, and ReadMap
	// decodes the absence back to the same nil map.
	if len(si.Annotations) > 0 {
		if err = wtr.WriteMap(c, cdc, "Annotations", &si.Annotations); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalJSON accepts both a plain string ID ("id123") and the full object form ({"Id":"id123",...}).
// UI forms typically send FK fields as plain string values, so this eliminates the need for
// transform: config entries that map e.g. Owner → Owner.Id.
func (si *StorableRef) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	// Plain string: "id123" — treat as Id
	if data[0] == '"' {
		var id string
		if err := json.Unmarshal(data, &id); err != nil {
			return err
		}
		si.Id = id
		return nil
	}
	// Object form: {"Id":"...","Type":"...",...} — unmarshal via alias to avoid infinite recursion.
	type storableRefAlias StorableRef
	return json.Unmarshal(data, (*storableRefAlias)(si))
}

func StorableArrayToMap(items []core.Storable) map[string]core.Storable {
	res := make(map[string]core.Storable, len(items))
	for _, item := range items {
		res[item.GetId()] = item
	}
	return res
}

//Factory function for creating storable
//type StorableCreator func() interface{}

// CastToStorableCollection converts a slice of records into storables and their ids, dropping any
// entry that is soft-deleted or nil.
//
// Both returned slices are sized to the INPUT length and filled from index 0 up to the number of
// entries that survived, so both must be truncated to that count. Returning ids at its full length
// left one empty string per dropped entry, which callers see as a record identity that does not
// exist — and providers derive recsreturned from len(ids), so a page containing a deleted record
// also reported more records than it returned and could inflate totalrecs above the true count.
func CastToStorableCollection(cx ctx.Context, items interface{}) ([]core.Storable, []string, error) {
	arr := reflect.ValueOf(items)
	if arr.Kind() == reflect.Ptr {
		arr = arr.Elem()
	}
	if arr.Kind() != reflect.Slice {
		return nil, nil, fmt.Errorf("Invalid cast to Storable. Type of Item: %s", arr.Kind())
	}
	length := arr.Len()
	retVal := make([]core.Storable, length)
	ids := make([]string, length)
	j := 0
	for i := 0; i < length; i++ {
		itemKind := arr.Index(i).Kind()
		var valPtr interface{}
		if itemKind == reflect.Ptr {
			valPtr = arr.Index(i).Interface()
		} else {
			valPtr = arr.Index(i).Addr().Interface()
		}
		if valPtr != nil {
			stor, ok := valPtr.(core.Storable)
			if !ok {
				return nil, nil, fmt.Errorf("Invalid cast to Storable. Item: %s", valPtr)
			}
			softDeletable, ok := stor.(SoftDeletable)
			if ok && softDeletable.IsDeleted() {
				continue
			}
			ids[j] = stor.GetId()
			retVal[j] = stor
			j++
		} else {
			log.Warn(cx, "Nil object received", slog.Int("index", i))
		}
	}
	// both slices carry j survivors; truncating only one of them is what produced the phantom ids
	return retVal[0:j], ids[0:j], nil
}

func CastToStorableHash(items interface{}) (map[string]core.Storable, error) {
	arr := reflect.ValueOf(items)
	if arr.Kind() == reflect.Ptr {
		arr = arr.Elem()
	}
	if arr.Kind() != reflect.Slice {
		return nil, fmt.Errorf("Invalid cast to Storable. Type of Item: %s", arr.Kind())
	}
	length := arr.Len()
	retVal := make(map[string]core.Storable, length)
	for i := 0; i < length; i++ {
		itemKind := arr.Index(i).Kind()
		var valPtr interface{}
		if itemKind == reflect.Ptr {
			valPtr = arr.Index(i).Interface()
		} else {
			valPtr = arr.Index(i).Addr().Interface()
		}
		stor, ok := valPtr.(core.Storable)
		if !ok {
			return nil, fmt.Errorf("Invalid cast to Storable. Item: %s %s %t", valPtr, arr.Index(i).Kind(), arr.Index(i).IsNil())
		}
		softDeletable, ok := stor.(SoftDeletable)
		if ok && softDeletable.IsDeleted() {
			continue
		}
		retVal[stor.GetId()] = stor
	}
	return retVal, nil
}
