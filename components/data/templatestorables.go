package data

import (
	"github.com/twinj/uuid"
	"laatoo/sdk/core"
	"time"
)

type AbstractStorable struct {
	ID string `json:"ID" bson:"ID" sql:"size:50;unique;index" gorm:"primary_key"`
}

func NewAbstractStorable() AbstractStorable {
	return AbstractStorable{uuid.NewV4().String()}
}
func (as *AbstractStorable) Init() {
	as.ID = uuid.NewV4().String()
}
func (as *AbstractStorable) GetId() string {
	return as.ID
}
func (as *AbstractStorable) SetId(val string) {
	as.ID = val
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
func (as *AbstractStorable) IsDeleted() bool {
	return false
}
func (as *AbstractStorable) Delete() {
}

func (as *AbstractStorable) Join(item Storable) {
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
	CreatedBy        string    `json:"CreatedBy" bson:"CreatedBy"`
	UpdatedBy        string    `json:"UpdatedBy" bson:"UpdatedBy" `
	CreatedAt        time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt        time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
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
	CreatedBy          string    `json:"CreatedBy" bson:"CreatedBy"`
	UpdatedBy          string    `json:"UpdatedBy" bson:"UpdatedBy" `
	CreatedAt          time.Time `json:"CreatedAt" bson:"CreatedAt"`
	UpdatedAt          time.Time `json:"UpdatedAt" bson:"UpdatedAt"`
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
func (hda *SoftDeleteAuditable) SetUpdatedBy(val string) {
	hda.UpdatedBy = val
}
func (hda *SoftDeleteAuditable) SetCreatedBy(val string) {
	hda.CreatedBy = val
}
func (hda *SoftDeleteAuditable) GetCreatedBy() string {
	return hda.CreatedBy
}
