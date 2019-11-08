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
*/

type SerializableBase struct {
}

func (b *SerializableBase) Reset() {
	*b = reflect.New(reflect.TypeOf(b).Elem()).Elem().Interface().(SerializableBase)
}

func (m *SerializableBase) String() string { return proto.CompactTextString(m) }

func (*SerializableBase) ProtoMessage() {}

type AbstractStorable struct {
	Id string `json:"Id" bson:"Id" protobuf:"bytes,51,opt,name=id,proto3" sql:"type:varchar(100); primary key; unique;index" gorm:"primary_key"`
}

func NewAbstractStorable() *AbstractStorable {
	return &AbstractStorable{Id: uuid.NewV4().String()}
}

func (as *AbstractStorable) Constructor() {
	if as.Id != "" {
		as.Id = uuid.NewV4().String()
	}
}

func (as *AbstractStorable) Initialize(ctx ctx.Context, conf config.Config) error {
	return nil
}

func (as *AbstractStorable) GetId() string {
	return as.Id
}
func (as *AbstractStorable) SetId(val string) {
	as.Id = val
}
func (as *AbstractStorable) GetLabel() string {
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

func (as *AbstractStorable) PreSave(ctx core.RequestContext) error {
	return nil
}
func (as *AbstractStorable) PostSave(ctx core.RequestContext) error {
	return nil
}
func (as *AbstractStorable) PostLoad(ctx core.RequestContext) error {
	return nil
}
func (as *AbstractStorable) SetValues(obj interface{}, val map[string]interface{}) {
	delete(val, "Id")
	delete(val, "IsNew")
	delete(val, "CreatedBy")
	delete(val, "UpdatedBy")
	delete(val, "CreatedAt")
	delete(val, "UpdatedAt")
	utils.SetObjectFields(obj, val)
}

func (as *AbstractStorable) IsMultitenant() bool {
	return false
}

func (as *AbstractStorable) Join(item Storable) {
}
func (as *AbstractStorable) Config() *StorableConfig {
	return nil
}

func (b *AbstractStorable) Reset() {
	*b = reflect.New(reflect.TypeOf(b).Elem()).Elem().Interface().(AbstractStorable)
	//	b.Id = uuid.NewV4().String()
}

func (m *AbstractStorable) String() string { return proto.CompactTextString(m) }

func (*AbstractStorable) ProtoMessage() {}

type SoftDeleteStorable struct {
	*AbstractStorable `json:",inline" initialize:"AbstractStorable" protobuf:"group,62,opt,name=AbstractStorable,proto3"`
	Deleted           bool `json:"Deleted" bson:"Deleted" protobuf:"bytes,52,opt,name=deleted,proto3"`
}

func NewSoftDeleteStorable() *SoftDeleteStorable {
	return &SoftDeleteStorable{NewAbstractStorable(), false}
}
func (sds *SoftDeleteStorable) IsDeleted() bool {
	return sds.Deleted
}
func (sds *SoftDeleteStorable) SoftDeleteField() string {
	return "Deleted"
}

type HardDeleteAuditable struct {
	*AbstractStorable `json:",inline" initialize:"AbstractStorable" protobuf:"group,62,opt,name=AbstractStorable,proto3"`
	New               bool      `json:"IsNew" bson:"IsNew" protobuf:"bytes,53,opt,name=isnew,proto3"`
	CreatedBy         string    `json:"CreatedBy" bson:"CreatedBy" protobuf:"bytes,54,opt,name=createdby,proto3" gorm:"column:CreatedBy"`
	UpdatedBy         string    `json:"UpdatedBy" bson:"UpdatedBy" protobuf:"bytes,55,opt,name=updatedby,proto3" gorm:"column:UpdatedBy"`
	CreatedAt         time.Time `json:"CreatedAt" bson:"CreatedAt" protobuf:"bytes,56,opt,name=createdat,proto3" gorm:"column:CreatedAt"`
	UpdatedAt         time.Time `json:"UpdatedAt" bson:"UpdatedAt" protobuf:"bytes,57,opt,name=updatedat,proto3" gorm:"column:UpdatedAt"`
}

func NewHardDeleteAuditable() *HardDeleteAuditable {
	return &HardDeleteAuditable{AbstractStorable: NewAbstractStorable()}
}
func (hda *HardDeleteAuditable) IsNew() bool {
	return hda.New
}
func (hda *HardDeleteAuditable) PreSave(ctx core.RequestContext) error {
	hda.New = (hda.CreatedBy == "")
	return nil
}

func (hda *HardDeleteAuditable) SetCreatedAt(val time.Time) {
	hda.CreatedAt = val
}
func (hda *HardDeleteAuditable) GetCreatedAt() time.Time {
	return hda.CreatedAt
}

func (hda *HardDeleteAuditable) SetUpdatedAt(val time.Time) {
	hda.UpdatedAt = val
}
func (hda *HardDeleteAuditable) GetUpdatedAt() time.Time {
	return hda.UpdatedAt
}

func (hda *HardDeleteAuditable) SetUpdatedBy(val string) {
	hda.UpdatedBy = val
}
func (hda *HardDeleteAuditable) GetUpdatedBy() string {
	return hda.UpdatedBy
}

func (hda *HardDeleteAuditable) SetCreatedBy(val string) {
	hda.CreatedBy = val
}
func (hda *HardDeleteAuditable) GetCreatedBy() string {
	return hda.CreatedBy
}

type SoftDeleteAuditable struct {
	*SoftDeleteStorable `json:",inline" initialize:"SoftDeleteStorable" protobuf:"group,63,opt,name=SoftDeleteStorable,proto3"`
	New                 bool      `json:"IsNew" bson:"IsNew" protobuf:"bytes,53,opt,name=isnew,proto3"`
	CreatedBy           string    `json:"CreatedBy" bson:"CreatedBy" protobuf:"bytes,54,opt,name=createdby,proto3" gorm:"column:CreatedBy"`
	UpdatedBy           string    `json:"UpdatedBy" bson:"UpdatedBy" protobuf:"bytes,55,opt,name=updatedby,proto3" gorm:"column:UpdatedBy"`
	CreatedAt           time.Time `json:"CreatedAt" bson:"CreatedAt" protobuf:"bytes,56,opt,name=createdat,proto3" gorm:"column:CreatedAt"`
	UpdatedAt           time.Time `json:"UpdatedAt" bson:"UpdatedAt" protobuf:"bytes,57,opt,name=updatedat,proto3" gorm:"column:UpdatedAt"`
}

func NewSoftDeleteAuditable() *SoftDeleteAuditable {
	return &SoftDeleteAuditable{SoftDeleteStorable: NewSoftDeleteStorable()}
}
func (sda *SoftDeleteAuditable) IsNew() bool {
	return sda.New
}
func (sda *SoftDeleteAuditable) PreSave(ctx core.RequestContext) error {
	sda.New = (sda.CreatedBy == "")
	return nil
}

func (sda *SoftDeleteAuditable) SetCreatedAt(val time.Time) {
	sda.CreatedAt = val
}
func (sda *SoftDeleteAuditable) GetCreatedAt() time.Time {
	return sda.CreatedAt
}

func (sda *SoftDeleteAuditable) SetUpdatedAt(val time.Time) {
	sda.UpdatedAt = val
}
func (sda *SoftDeleteAuditable) GetUpdatedAt() time.Time {
	return sda.UpdatedAt
}

func (sda *SoftDeleteAuditable) SetUpdatedBy(val string) {
	sda.UpdatedBy = val
}
func (sda *SoftDeleteAuditable) GetUpdatedBy() string {
	return sda.UpdatedBy
}

func (sda *SoftDeleteAuditable) SetCreatedBy(val string) {
	sda.CreatedBy = val
}
func (sda *SoftDeleteAuditable) GetCreatedBy() string {
	return sda.CreatedBy
}
