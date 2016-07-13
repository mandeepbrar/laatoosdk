package data

import (
	"github.com/twinj/uuid"
	"laatoo/sdk/core"
)

type AbstractStorable struct {
	Id string `json:"Id" bson:"Id"`
}

func NewAbstractStorable() AbstractStorable {
	return AbstractStorable{uuid.NewV4().String()}
}
func (as *AbstractStorable) Init() {
	as.Id = uuid.NewV4().String()
}
func (as *AbstractStorable) GetId() string {
	return as.Id
}
func (as *AbstractStorable) SetId(val string) {
	as.Id = val
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
	CreatedBy        string `json:"CreatedBy" bson:"CreatedBy"`
	UpdatedBy        string `json:"UpdatedBy" bson:"UpdatedBy" `
	UpdatedOn        string `json:"UpdatedOn" bson:"UpdatedOn"`
}

func NewHardDeleteAuditable() HardDeleteAuditable {
	return HardDeleteAuditable{NewAbstractStorable(), "", "", ""}
}
func (hda *HardDeleteAuditable) IsNew() bool {
	return hda.CreatedBy == ""
}
func (hda *HardDeleteAuditable) SetUpdatedOn(val string) {
	hda.UpdatedOn = val
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
	CreatedBy          string `json:"CreatedBy" bson:"CreatedBy"`
	UpdatedBy          string `json:"UpdatedBy" bson:"UpdatedBy" `
	UpdatedOn          string `json:"UpdatedOn" bson:"UpdatedOn"`
}

func NewSoftDeleteAuditable() SoftDeleteAuditable {
	return SoftDeleteAuditable{NewSoftDeleteStorable(), "", "", ""}
}
func (hda *SoftDeleteAuditable) IsNew() bool {
	return hda.CreatedBy == ""
}
func (hda *SoftDeleteAuditable) SetUpdatedOn(val string) {
	hda.UpdatedOn = val
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
