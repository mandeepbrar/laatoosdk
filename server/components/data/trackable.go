package data

import (
	"time"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/log"
)

// /auditable entities must have UpdatedBy and UpdatedOn fields to support auditing through update queries
type Trackable interface {
	IsNew() bool
	SetCreatedAt(time.Time)
	GetCreatedAt() time.Time
	SetUpdatedAt(time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedBy(string)
	GetUpdatedBy() string
	SetCreatedBy(string)
	GetCreatedBy() string
}

type TrackingInfo struct {
	New       bool      `json:"IsNew" bson:"IsNew" protobuf:"bytes,53,opt,name=isnew,proto3"`
	CreatedBy string    `json:"CreatedBy" bson:"CreatedBy" protobuf:"bytes,54,opt,name=createdby,proto3" gorm:"column:CreatedBy"`
	UpdatedBy string    `json:"UpdatedBy" bson:"UpdatedBy" protobuf:"bytes,55,opt,name=updatedby,proto3" gorm:"column:UpdatedBy"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt" protobuf:"bytes,56,opt,name=createdat,proto3" gorm:"column:CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt" protobuf:"bytes,57,opt,name=updatedat,proto3" gorm:"column:UpdatedAt"`
}

func (ti *TrackingInfo) SetCreatedAt(val time.Time) {
	ti.CreatedAt = val
}
func (ti *TrackingInfo) GetCreatedAt() time.Time {
	return ti.CreatedAt
}

func (ti *TrackingInfo) SetUpdatedAt(val time.Time) {
	ti.UpdatedAt = val
}
func (ti *TrackingInfo) GetUpdatedAt() time.Time {
	return ti.UpdatedAt
}

func (ti *TrackingInfo) SetUpdatedBy(val string) {
	ti.UpdatedBy = val
}
func (ti *TrackingInfo) GetUpdatedBy() string {
	return ti.UpdatedBy
}

func (ti *TrackingInfo) SetCreatedBy(val string) {
	ti.CreatedBy = val
}
func (ti *TrackingInfo) GetCreatedBy() string {
	return ti.CreatedBy
}

func Track(ctx core.RequestContext, item interface{}) {
	if item != nil {
		auditable, aok := item.(Trackable)
		if aok {
			usr := ctx.GetUser()
			if usr != nil {
				id := usr.GetId()
				if auditable.IsNew() {
					auditable.SetCreatedBy(id)
				}
				auditable.SetUpdatedBy(id)
				tim := time.Now()
				if auditable.IsNew() {
					auditable.SetCreatedAt(tim)
				}
				auditable.SetUpdatedAt(tim)
			} else {
				log.Info(ctx, "Could not audit entity. User nil")
			}
		} else {
			updateMap, mapok := item.(map[string]interface{})
			if mapok {
				usr := ctx.GetUser()
				if usr != nil {
					id := usr.GetId()
					updateMap["UpdatedBy"] = id
					updateMap["UpdatedAt"] = time.Now()
				} else {
					log.Info(ctx, "Could not audit map. User nil")
				}
			}
		}
	}
}

func (ti *TrackingInfo) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error
	if err = rdr.ReadString(c, cdc, "CreatedBy", &ti.CreatedBy); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "UpdatedBy", &ti.UpdatedBy); err != nil {
		return err
	}
	if err = rdr.ReadTime(c, cdc, "CreatedAt", &ti.CreatedAt); err != nil {
		return err
	}
	if err = rdr.ReadTime(c, cdc, "UpdatedAt", &ti.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (ti *TrackingInfo) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error
	if err = wtr.WriteString(c, cdc, "CreatedBy", &ti.CreatedBy); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "UpdatedBy", &ti.UpdatedBy); err != nil {
		return err
	}
	if err = wtr.WriteTime(c, cdc, "CreatedAt", &ti.CreatedAt); err != nil {
		return err
	}
	if err = wtr.WriteTime(c, cdc, "UpdatedAt", &ti.UpdatedAt); err != nil {
		return err
	}
	return nil
}
