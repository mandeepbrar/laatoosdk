package data

import (
	"laatoo.io/sdk/server/auth"
)

// Multitenant is implemented by an entity that carries a tenant partition, via the embedded
// TenantInfo struct in this package. Unlike DeletionInfo and TrackingInfo, TenantInfo is embedded
// only when the entity declares multitenant, so a type assertion to Multitenant genuinely
// distinguishes tenant-scoped entities from the rest — and every provider performs it in the
// comma-ok form before stamping or comparing a tenant.
//
// THE TENANT COMES FROM THE USER, NOT FROM THE CONTEXT. Every read, write and delete resolves it
// as ctx.GetUser().GetTenant() — never ctx.GetTenant(). See
// laatoomodules/datastore/dev/plugins/common/basecomponent.go:212-214 (the read filter),
// :562-569 (ValidateTenant) and mongodatabase mongodataservice.go:659-666 (setTenant on write).
// Getting this backwards does not raise an error: it either partitions on an empty tenant and
// returns zero rows, or writes rows a later read of the correct tenant will never see.
//
// A multitenant component with no user at all now REFUSES the read
// (basecomponent.go:192-202, RequireTenantScope) rather than returning every tenant's rows, which
// is what it used to do.
type Multitenant interface {
	auth.TenantInfo

	// GetTenantInfo returns the entity's own tenant as an auth.TenantInfo. The embedded
	// implementation returns the receiver itself (tenantinfo.go:26-28), so it is never nil once
	// the assertion to Multitenant has succeeded — which is why ValidateTenant can chain
	// mt.GetTenantInfo().GetTenantId() without a guard (basecomponent.go:567).
	GetTenantInfo() auth.TenantInfo

	// SetTenant assigns both halves of the partition. TenantId is the value every generated
	// filter compares against; TenantName is descriptive only and is never queried.
	//
	// It overwrites unconditionally, including with empty strings — the embedded implementation
	// is a plain two-field assignment (tenantinfo.go:21-24) with no validation that the caller is
	// entitled to the tenant it names. Providers therefore call it only from setTenant, sourcing
	// the value from the request's user; nothing else should call it on a record being saved.
	SetTenant(tenantid, tenantname string)

	// SetTenantInfo copies the partition from an existing auth.TenantInfo — the form providers
	// use on save, passing ctx.GetUser().GetTenant() straight through.
	//
	// A nil argument is a SILENT NO-OP (tenantinfo.go:31-33): the entity keeps whatever tenant it
	// already had, which for a freshly created object is the empty string. Since GetTenant()
	// returns nil for an anonymous or machine caller, a save on that path leaves the row
	// untenanted and invisible to every subsequent tenant-scoped read, with no error raised
	// anywhere. The provider's own guard is the comma-ok on ctx.GetUser() before the call
	// (mongodataservice.go:661), not anything inside this method.
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
