package data

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/datatypes"
	"laatoo.io/sdk/server/auth"
)

type TenantInfo struct {
	datatypes.Serializable
	TenantName string `json:"TenantName" protobuf:"bytes,73,opt,name=tenantname,proto3" bson:"TenantName" sql:"type:varchar(100);"`
	TenantId   string `json:"TenantId" protobuf:"bytes,61,opt,name=tenantid,proto3" bson:"TenantId" sql:"type:varchar(100);"`
}

func (ti *TenantInfo) GetTenantId() string {
	return ti.TenantId
}
func (ti *TenantInfo) GetTenantName() string {
	return ti.TenantName
}

func (ti *TenantInfo) SetTenant(tenantid, tenantname string) {
	ti.TenantId = tenantid
	ti.TenantName = tenantname
}

func (ti *TenantInfo) SetTenantInfo(inf auth.TenantInfo) {
	if inf != nil {
		ti.SetTenant(inf.GetTenantId(), inf.GetTenantName())
	}
}

func (ti *TenantInfo) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error
	if err = rdr.ReadString(c, cdc, "TenantId", &ti.TenantId); err != nil {
		return err
	}
	if err = rdr.ReadString(c, cdc, "TenantName", &ti.TenantName); err != nil {
		return err
	}
	return nil
}

func (ti *TenantInfo) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	var err error
	if err = wtr.WriteString(c, cdc, "TenantId", &ti.TenantId); err != nil {
		return err
	}
	if err = wtr.WriteString(c, cdc, "TenantName", &ti.TenantName); err != nil {
		return err
	}
	/*var err error
	if err = wtr.WriteString(c, cdc, "Tenant", &ti.Tenant); err != nil {
		return err
	}*/
	return nil
}
