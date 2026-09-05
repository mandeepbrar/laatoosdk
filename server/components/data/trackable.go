package data

import (
	"time"

	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/core"
	"laatoo.io/sdk/server/log"
)

// Trackable is the audit-stamp contract, satisfied by the embedded TrackingInfo struct below.
// Auditable entities must have UpdatedBy and UpdatedAt fields to support auditing through update
// queries.
//
// The SETTERS ARE NOT FOR CALLERS. Track (below) is the only production caller of any of them,
// across every datastore provider — see mongodatabase mongodataservice.go:174, sqldatabase
// sqldataservice.go:250, boltdatabase kvdataservice.go:257, couchbasedatabase
// couchbasedataservice.go:160, gaedatastore gaedataservice.go:152 and their sibling write paths.
// Stamping a field yourself before Save is at best redundant and at worst overwritten.
//
// Two conditions gate the whole mechanism, and neither reports anything when it fails:
//
//   - The component must be trackable. The flag comes from StorableConfig.Trackable, overridable
//     by the module's auditable setting (basecomponent.go:72-77). If it is off, Track is never
//     reached and all four fields stay at their zero values.
//   - The request must have a user. Track reads ctx.GetUser().GetId(); a nil user skips every
//     stamp and logs at INFO ("Could not audit entity. User nil"), which is not an error and does
//     not fail the save.
//
// The codegen embeds TrackingInfo in every generated entity regardless of the Trackable flag, so
// satisfying this interface does not mean auditing is enabled for the entity.
type Trackable interface {
	// IsNew reports whether this record has never been persisted, which decides whether Track
	// stamps the created fields as well as the updated ones. The embedded implementation derives
	// it from CreatedAt rather than from a stored flag; see the comment on TrackingInfo.IsNew for
	// why.
	IsNew() bool

	// SetCreatedAt records the creation instant. Track calls it only on the first tracked save
	// (when IsNew reports true), using the same time.Now() value it gives SetUpdatedAt.
	SetCreatedAt(time.Time)

	// GetCreatedAt returns the creation instant, zero until the first tracked save. It is also
	// what IsNew tests, so it doubles as the persisted/not-persisted signal.
	GetCreatedAt() time.Time

	// SetUpdatedAt records the modification instant. Track calls it on every tracked write.
	SetUpdatedAt(time.Time)

	// GetUpdatedAt returns the modification instant, zero until the first tracked save.
	GetUpdatedAt() time.Time

	// SetUpdatedBy records the modifying user's id. Track calls it on every tracked write.
	SetUpdatedBy(string)

	// GetUpdatedBy returns the id of the user who last wrote the record — an id, never a name or
	// an email.
	GetUpdatedBy() string

	// SetCreatedBy records the creating user's id. Track calls it only on the first tracked save.
	SetCreatedBy(string)

	// GetCreatedBy returns the id of the user who created the record.
	//
	// It stays empty on any entity whose first save happened while the component was not trackable
	// or while the request had no user, and nothing later fills it in: the update-by-map path
	// takes the non-Trackable branch of Track, which stamps only UpdatedBy and UpdatedAt into the
	// value map and never the created pair (trackable.go:96-106).
	GetCreatedBy() string
}

type TrackingInfo struct {
	New       bool      `json:"IsNew" bson:"IsNew" protobuf:"bytes,53,opt,name=isnew,proto3"`
	CreatedBy string    `json:"CreatedBy" bson:"CreatedBy" protobuf:"bytes,54,opt,name=createdby,proto3" gorm:"column:CreatedBy"`
	UpdatedBy string    `json:"UpdatedBy" bson:"UpdatedBy" protobuf:"bytes,55,opt,name=updatedby,proto3" gorm:"column:UpdatedBy"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt" protobuf:"bytes,56,opt,name=createdat,proto3" gorm:"column:CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt" bson:"UpdatedAt" protobuf:"bytes,57,opt,name=updatedat,proto3" gorm:"column:UpdatedAt"`
}

// IsNew reports whether this record has never been persisted, which is what decides whether Track
// stamps the created fields as well as the updated ones.
//
// It derives that from CreatedAt rather than from the New field. Nothing in the platform ever sets
// New -- not the SDK, not the server, not any module -- and TrackingInfo.WriteAll does not persist
// it, so `return ti.New` would satisfy the interface and still leave CreatedBy and CreatedAt
// permanently unwritten. CreatedAt is stamped by Track itself and does round-trip through storage,
// so it is zero exactly until the first tracked save and non-zero forever after.
//
// The New field is left in place because it is part of the stored column layout gorm derives from
// this struct; removing it would be a storage change rather than a code one.
func (ti *TrackingInfo) IsNew() bool {
	return ti.CreatedAt.IsZero()
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

// WriteAll emits the four tracking fields and deliberately NOT New (json tag "IsNew"). Nothing
// in the platform sets New — IsNew() derives newness from CreatedAt, as documented above — so
// the field would be false on every record ever encoded. encoding/json over the same struct does
// emit it from the tag; the two wire forms differ here by decision, recorded 2026-09-05 (use case
// codec-encoding, "wire-format divergences"). Emitting it is a wire-format change, not a fix.
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
