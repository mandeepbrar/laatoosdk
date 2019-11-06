package data

import (
	"laatoo/sdk/common/config"
	"laatoo/sdk/server/core"
	"laatoo/sdk/server/ctx"
	"laatoo/sdk/utils"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/twinj/uuid"
)

type AbstractStorableMT struct {
	Id     string `json:"Id" protobuf:"bytes,51,opt,name=id,proto3" bson:"Id" sql:"type:varchar(100); primary key; unique;index" gorm:"primary_key"`
	Tenant string `json:"Tenant" protobuf:"bytes,61,opt,name=tenant,proto3" bson:"Tenant" sql:"type:varchar(100);"`
}

func NewAbstractStorableMT() *AbstractStorableMT {
	return &AbstractStorableMT{Id: uuid.NewV4().String()}
}
func (as *AbstractStorableMT) Constructor() {
	if as.Id != "" {
		as.Id = uuid.NewV4().String()
	}
}
func (as *AbstractStorableMT) Initialize(ctx ctx.Context, conf config.Config) error {
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

func (as *AbstractStorableMT) String() string {
	return proto.CompactTextString(as)
}
func (as *AbstractStorableMT) ProtoMessage() {}

func (b *AbstractStorableMT) Reset() {
	*b = reflect.New(reflect.TypeOf(b).Elem()).Elem().Interface().(AbstractStorableMT)
	//	b.Id = uuid.NewV4().String()
}

type SoftDeleteStorableMT struct {
	*AbstractStorableMT `json:",inline" initialize:"AbstractStorableMT" protobuf:"group,65,opt,name=AbstractStorableMT,proto3"`
	Deleted             bool `json:"Deleted" bson:"Deleted"`
}

func NewSoftDeleteStorableMT() *SoftDeleteStorableMT {
	return &SoftDeleteStorableMT{NewAbstractStorableMT(), false}
}
func (sds *SoftDeleteStorableMT) IsDeleted() bool {
	return sds.Deleted
}
func (sds *SoftDeleteStorableMT) Delete() {
	sds.Deleted = true
}

type HardDeleteAuditableMT struct {
	*AbstractStorableMT `json:",inline" initialize:"AbstractStorableMT" protobuf:"group,65,opt,name=AbstractStorableMT,proto3"`
	New                 bool      `json:"IsNew" bson:"IsNew" protobuf:"bytes,53,opt,name=isnew,proto3"`
	CreatedBy           string    `json:"CreatedBy" bson:"CreatedBy" protobuf:"bytes,54,opt,name=createdby,proto3" gorm:"column:CreatedBy"`
	UpdatedBy           string    `json:"UpdatedBy" bson:"UpdatedBy" protobuf:"bytes,55,opt,name=updatedby,proto3" gorm:"column:UpdatedBy"`
	CreatedAt           time.Time `json:"CreatedAt" bson:"CreatedAt" protobuf:"bytes,56,opt,name=createdat,proto3" gorm:"column:CreatedAt"`
	UpdatedAt           time.Time `json:"UpdatedAt" bson:"UpdatedAt" protobuf:"bytes,57,opt,name=updatedat,proto3" gorm:"column:UpdatedAt"`
}

func NewHardDeleteAuditableMT() *HardDeleteAuditableMT {
	return &HardDeleteAuditableMT{AbstractStorableMT: NewAbstractStorableMT()}
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
	*SoftDeleteStorableMT `json:",inline" initialize:"SoftDeleteStorableMT" protobuf:"group,66,opt,name=SoftDeleteStorableMT,proto3"`
	New                   bool      `json:"IsNew" bson:"IsNew" protobuf:"bytes,53,opt,name=isnew,proto3"`
	CreatedBy             string    `json:"CreatedBy" bson:"CreatedBy" protobuf:"bytes,54,opt,name=createdby,proto3" gorm:"column:CreatedBy"`
	UpdatedBy             string    `json:"UpdatedBy" bson:"UpdatedBy" protobuf:"bytes,55,opt,name=updatedby,proto3" gorm:"column:UpdatedBy"`
	CreatedAt             time.Time `json:"CreatedAt" bson:"CreatedAt" protobuf:"bytes,56,opt,name=createdat,proto3" gorm:"column:CreatedAt"`
	UpdatedAt             time.Time `json:"UpdatedAt" bson:"UpdatedAt" protobuf:"bytes,57,opt,name=updatedat,proto3" gorm:"column:UpdatedAt"`
}

func NewSoftDeleteAuditableMT() *SoftDeleteAuditableMT {
	return &SoftDeleteAuditableMT{SoftDeleteStorableMT: NewSoftDeleteStorableMT()}
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
