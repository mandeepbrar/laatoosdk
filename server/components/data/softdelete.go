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

// Object stored by data service
type SoftDeletable interface {
	IsDeleted() bool
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
