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

type AbstractStorable struct {
	Id    string `json:"Id" bson:"Id" sql:"type:varchar(100); primary key; unique;index" gorm:"primary_key"`
	Empty string `json:"-" bson:"-" sql:"-"`
}

func NewAbstractStorable() AbstractStorable {
	return AbstractStorable{Id: uuid.NewV4().String()}
}
func (as *AbstractStorable) Initialize(ctx ctx.Context, conf config.Config) error {
	as.Id = uuid.NewV4().String()
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
	utils.SetObjectFields(obj, val)
}
func (as *AbstractStorable) IsDeleted() bool {
	return false
}
func (as *AbstractStorable) Delete() {
}

func (as *AbstractStorable) IsMultitenant() bool {
	return false
}

func (as *AbstractStorable) Join(item Storable) {
}
func (as *AbstractStorable) Config() *StorableConfig {
	return nil
}

type SoftDeleteStorable struct {
	AbstractStorable `bson:",inline"`
	Deleted          bool `json:"Deleted" bson:"Deleted"`
}

func NewSoftDeleteStorable() SoftDeleteStorable {
	return SoftDeleteStorable{NewAbstractStorable(), false}
}
func (sds *SoftDeleteStorable) IsDeleted() bool {
	return sds.Deleted
}
func (sds *SoftDeleteStorable) Delete() {
	sds.Deleted = true
}

type HardDeleteAuditable struct {
	AbstractStorable `bson:",inline"`
	New              bool      `json:"IsNew" bson:"IsNew"`
	CreatedBy        string    `json:"CreatedBy" bson:"CreatedBy" gorm:"column:CreatedBy"`
	UpdatedBy        string    `json:"UpdatedBy" bson:"UpdatedBy" gorm:"column:UpdatedBy"`
	CreatedAt        time.Time `json:"CreatedAt" bson:"CreatedAt" gorm:"column:CreatedAt"`
	UpdatedAt        time.Time `json:"UpdatedAt" bson:"UpdatedAt" gorm:"column:UpdatedAt"`
}

func NewHardDeleteAuditable() HardDeleteAuditable {
	return HardDeleteAuditable{AbstractStorable: NewAbstractStorable()}
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
func (hda *HardDeleteAuditable) SetUpdatedAt(val time.Time) {
	hda.UpdatedAt = val
}

func (hda *HardDeleteAuditable) SetUpdatedBy(val string) {
	hda.UpdatedBy = val
}
func (hda *HardDeleteAuditable) SetCreatedBy(val string) {
	hda.CreatedBy = val
}
func (hda *HardDeleteAuditable) GetCreatedBy() string {
	return hda.CreatedBy
}

type SoftDeleteAuditable struct {
	SoftDeleteStorable `bson:",inline"`
	New                bool      `json:"IsNew" bson:"IsNew"`
	CreatedBy          string    `json:"CreatedBy" bson:"CreatedBy" gorm:"column:CreatedBy"`
	UpdatedBy          string    `json:"UpdatedBy" bson:"UpdatedBy" gorm:"column:UpdatedBy"`
	CreatedAt          time.Time `json:"CreatedAt" bson:"CreatedAt" gorm:"column:CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt" bson:"UpdatedAt" gorm:"column:UpdatedAt"`
}

func NewSoftDeleteAuditable() SoftDeleteAuditable {
	return SoftDeleteAuditable{SoftDeleteStorable: NewSoftDeleteStorable()}
}
func (hda *SoftDeleteAuditable) IsNew() bool {
	return hda.New
}
func (hda *SoftDeleteAuditable) PreSave(ctx core.RequestContext) error {
	hda.New = (hda.CreatedBy == "")
	return nil
}
func (hda *SoftDeleteAuditable) SetCreatedAt(val time.Time) {
	hda.CreatedAt = val
}
func (hda *SoftDeleteAuditable) SetUpdatedAt(val time.Time) {
	hda.UpdatedAt = val
}

func (hda *SoftDeleteAuditable) SetUpdatedBy(val string) {
	hda.UpdatedBy = val
}
func (hda *SoftDeleteAuditable) SetCreatedBy(val string) {
	hda.CreatedBy = val
}
func (hda *SoftDeleteAuditable) GetCreatedBy() string {
	return hda.CreatedBy
}
