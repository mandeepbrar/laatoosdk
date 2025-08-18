package data

import (
	"laatoo.io/sdk/server/auth"
)

// Object stored by data service
type Multitenant interface {
	auth.TenantInfo
	GetTenantInfo() auth.TenantInfo
	SetTenant(tenantid, tenantname string)
	SetTenantInfo(inf auth.TenantInfo)
}

/*
type TenantInfo struct {
	Tenant     string `json:"Tenant" protobuf:"bytes,61,opt,name=tenant,proto3" bson:"Tenant" sql:"type:varchar(100);"`
	TenantName string `json:"TenantName" protobuf:"bytes,73,opt,name=tenantname,proto3" bson:"Tenant" sql:"type:varchar(100);"`
}

func (ti *TenantInfo) GetTenantId() string {
	return ti.Tenant
}
func (ti *TenantInfo) GetTenantName() string {
	return ti.TenantName
}

func (ti *TenantInfo) SetTenant(tenantid, tenantname string) {
	ti.Tenant = tenantid
	ti.TenantName = tenantname
}
func (ti *TenantInfo) ReadAll(c ctx.Context, cdc datatypes.Codec, rdr datatypes.SerializableReader) error {
	var err error
	if err = rdr.ReadString(c, cdc, "Tenant", &ti.Tenant); err != nil {
		return err
	}
	return nil
}

func (ti *TenantInfo) WriteAll(c ctx.Context, cdc datatypes.Codec, wtr datatypes.SerializableWriter) error {
	/*var err error
	if err = wtr.WriteString(c, cdc, "Tenant", &ti.Tenant); err != nil {
		return err
	}
	return nil
}*/
