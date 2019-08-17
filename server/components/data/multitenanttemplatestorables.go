package data

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/core"
	"laatoo/sdk/server/ctx"
	"laatoo/sdk/utils"
	"reflect"
	"time"

	"github.com/twinj/uuid"
)

type AbstractStorableMT struct {
	Id     string `json:"Id" bson:"Id" sql:"type:varchar(100); primary key; unique;index" gorm:"primary_key"`
	Tenant string `json:"Tenant" bson:"Tenant" sql:"type:varchar(100);"`
	Empty  string `json:"-" bson:"-" sql:"-"`
}

func NewAbstractStorableMT() AbstractStorableMT {
	return AbstractStorableMT{Id: uuid.NewV4().String()}
}
func (as *AbstractStorableMT) Initialize(ctx ctx.Context, conf config.Config) error {
	as.Id = uuid.NewV4().String()
	return nil
}
func (as *AbstractStorableMT) GetId() string {
	return as.Id
}
func (as *AbstractStorableMT) SetId(val string) {
	as.Id = val
}
func (as *AbstractStorableMT) GetLabel() string {
	c := as.Config()
	if c != nil && c.LabelField != "" {
		v := reflect.ValueOf(c)
		f := v.FieldByName(c.LabelField)
		if !f.IsNil() {
			return f.String()
		}
	}
	return ""
}

func (as *AbstractStorableMT) PreSave(ctx core.RequestContext) error {
	return nil
}
func (as *AbstractStorableMT) PostSave(ctx core.RequestContext) error {
	return nil
}
func (as *AbstractStorableMT) PostLoad(ctx core.RequestContext) error {
	return nil
}
func (as *AbstractStorableMT) SetValues(obj interface{}, val map[string]interface{}) {
	utils.SetObjectFields(obj, val)
}
func (as *AbstractStorableMT) IsDeleted() bool {
	return false
}
func (as *AbstractStorableMT) Delete() {
}

func (as *AbstractStorableMT) IsMultitenant() bool {
	return true
}
func (as *AbstractStorableMT) GetTenant() string {
	return as.Tenant
}
func (as *AbstractStorableMT) SetTenant(tenant string) {
	as.Tenant = tenant
}

func (as *AbstractStorableMT) Join(item Storable) {
}
func (as *AbstractStorableMT) Config() *StorableConfig {
	return nil
}

type SoftDeleteStorableMT struct {
	AbstractStorableMT `bson:",inline"`
	Deleted            bool `json:"Deleted" bson:"Deleted"`
}

func NewSoftDeleteStorableMT() SoftDeleteStorableMT {
	return SoftDeleteStorableMT{NewAbstractStorableMT(), false}
}
func (sds *SoftDeleteStorableMT) IsDeleted() bool {
	return sds.Deleted
}
func (sds *SoftDeleteStorableMT) Delete() {
	sds.Deleted = true
}

type HardDeleteAuditableMT struct {
	AbstractStorableMT `bson:",inline"`
	New                bool      `json:"IsNew" bson:"IsNew"`
	CreatedBy          string    `json:"CreatedBy" bson:"CreatedBy" gorm:"column:CreatedBy"`
	UpdatedBy          string    `json:"UpdatedBy" bson:"UpdatedBy" gorm:"column:UpdatedBy"`
	CreatedAt          time.Time `json:"CreatedAt" bson:"CreatedAt" gorm:"column:CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt" bson:"UpdatedAt" gorm:"column:UpdatedAt"`
}

func NewHardDeleteAuditableMT() HardDeleteAuditableMT {
	return HardDeleteAuditableMT{AbstractStorableMT: NewAbstractStorableMT()}
}
func (hda *HardDeleteAuditableMT) IsNew() bool {
	return hda.New
}
func (hda *HardDeleteAuditableMT) PreSave(ctx core.RequestContext) error {
	hda.New = (hda.CreatedBy == "")
	return nil
}
func (hda *HardDeleteAuditableMT) SetCreatedAt(val time.Time) {
	hda.CreatedAt = val
}
func (hda *HardDeleteAuditableMT) SetUpdatedAt(val time.Time) {
	hda.UpdatedAt = val
}

func (hda *HardDeleteAuditableMT) SetUpdatedBy(val string) {
	hda.UpdatedBy = val
}
func (hda *HardDeleteAuditableMT) SetCreatedBy(val string) {
	hda.CreatedBy = val
}
func (hda *HardDeleteAuditableMT) GetCreatedBy() string {
	return hda.CreatedBy
}

type SoftDeleteAuditableMT struct {
	SoftDeleteStorableMT `bson:",inline"`
	New                  bool      `json:"IsNew" bson:"IsNew"`
	CreatedBy            string    `json:"CreatedBy" bson:"CreatedBy" gorm:"column:CreatedBy"`
	UpdatedBy            string    `json:"UpdatedBy" bson:"UpdatedBy" gorm:"column:UpdatedBy"`
	CreatedAt            time.Time `json:"CreatedAt" bson:"CreatedAt" gorm:"column:CreatedAt"`
	UpdatedAt            time.Time `json:"UpdatedAt" bson:"UpdatedAt" gorm:"column:UpdatedAt"`
}

func NewSoftDeleteAuditableMT() SoftDeleteAuditableMT {
	return SoftDeleteAuditableMT{SoftDeleteStorableMT: NewSoftDeleteStorableMT()}
}
func (hda *SoftDeleteAuditableMT) IsNew() bool {
	return hda.New
}
func (hda *SoftDeleteAuditableMT) PreSave(ctx core.RequestContext) error {
	hda.New = (hda.CreatedBy == "")
	return nil
}
func (hda *SoftDeleteAuditableMT) SetCreatedAt(val time.Time) {
	hda.CreatedAt = val
}
func (hda *SoftDeleteAuditableMT) SetUpdatedAt(val time.Time) {
	hda.UpdatedAt = val
}

func (hda *SoftDeleteAuditableMT) SetUpdatedBy(val string) {
	hda.UpdatedBy = val
}
func (hda *SoftDeleteAuditableMT) SetCreatedBy(val string) {
	hda.CreatedBy = val
}
func (hda *SoftDeleteAuditableMT) GetCreatedBy() string {
	return hda.CreatedBy
}
