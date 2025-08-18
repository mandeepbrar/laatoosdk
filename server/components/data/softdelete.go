package data

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
)

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
	if err = rdr.ReadBool(c, cdc, "Deleted", &di.Deleted); err != nil {
		return err
	}
	return nil
}

func (di *DeletionInfo) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error
	if err = wtr.WriteBool(c, cdc, "Deleted", &di.Deleted); err != nil {
		return err
	}
	return nil
}
