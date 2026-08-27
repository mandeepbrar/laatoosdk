package data

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
)

// FIELD_SOFTDELETE is the stored name of the soft-delete flag, on every provider.
//
// It is a constant rather than configuration because it is not variable: the name is fixed by
// DeletionInfo's own struct tags below, which every store reads. Mongo and Couchbase take it from
// bson/json; Firestore, Datastore and bolthold use the Go field name; and the SQL provider names
// columns through a gorm naming strategy set to the identity function, so the column is "Deleted"
// there too rather than snake-cased.
//
// Exposing it as a config key would let a deployment set a name that exists in no store, and the
// result would be silent: an empty result set on SQL and a predicate matching nothing on the
// schemaless stores, with no error anywhere.
const FIELD_SOFTDELETE = "Deleted"

// SoftDeletable is the read side of soft deletion. Every generated entity satisfies it, because
// the codegen embeds DeletionInfo unconditionally whether or not SoftDelete is configured — so a
// successful type assertion to SoftDeletable says nothing about whether soft delete is ENABLED.
// That decision lives on the data component (StorableConfig.SoftDelete, overridable by the module
// setting), not on the entity.
//
// The interface is asymmetric in practice: reads consult IsDeleted constantly, while the write is
// performed as a field update rather than through SetDeleted. See each method.
type SoftDeletable interface {
	// IsDeleted reports the stored flag. It is the guard used at two distinct layers, and both
	// must hold for a deleted record to stay hidden:
	//
	//   - the query layer, which appends Deleted == false to every read on a soft-delete
	//     component (datastore/dev/plugins/common/basecomponent.go:215-219); and
	//   - the materialisation layer, where CastToStorableCollection and CastToStorableHash
	//     (storageinfo.go:240, :277) drop flagged records after decoding, as do the per-provider
	//     result loops (e.g. boltdatabase kvdataservice_get.go:79).
	//
	// The second layer drops records AFTER the page was cut, so on a provider that filters only
	// there, a page can come back shorter than pageSize while more matching records exist.
	IsDeleted() bool

	// SetDeleted writes the flag on the in-memory object.
	//
	// No production code path calls it. Delete/DeleteMulti/DeleteAll on a soft-delete component
	// do not mutate a loaded entity — they rewrite the request as an Update carrying the value map
	// utils.StringMap{data.FIELD_SOFTDELETE: true} (mongodatabase mongodataservice.go:588-593,
	// :611-615, :630-634, and the same shape in every other provider). Calling SetDeleted yourself
	// and then Save therefore does NOT go through the delete path: no delete event is emitted and
	// no delete permission is consulted. Call Delete.
	SetDeleted(deleted bool)
}

type DeletionInfo struct {
	Deleted bool `json:"Deleted" bson:"Deleted" protobuf:"bytes,52,opt,name=deleted,proto3"`
}

func (di *DeletionInfo) IsDeleted() bool {
	return di.Deleted
}

func (di *DeletionInfo) SetDeleted(deleted bool) {
	di.Deleted = deleted
}
func (di *DeletionInfo) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error
	if err = rdr.ReadBool(c, cdc, FIELD_SOFTDELETE, &di.Deleted); err != nil {
		return err
	}
	return nil
}

func (di *DeletionInfo) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error
	if err = wtr.WriteBool(c, cdc, FIELD_SOFTDELETE, &di.Deleted); err != nil {
		return err
	}
	return nil
}
