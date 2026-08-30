package core

import (
	"laatoo.io/sdk/ctx"
	"laatoo.io/sdk/utils"
)

type StorableConfig struct {
	ObjectType        string
	LabelField        string
	PartialLoadFields []string
	FullLoadFields    []string
	PreSave           bool
	PostSave          bool
	PostUpdate        bool
	PostLoad          bool
	Trackable         bool
	Collection        string
	Cacheable         bool
	RefOps            bool
	Workflow          bool
	Multitenant       bool
	// Namespace places this entity's data component in a DATA namespace, forming the other half
	// of the (namespace, object) key the DataManager registry uses. Declared in the entity YAML as
	// `namespace`, and empty means NAMESPACE_DEFAULT.
	//
	// A MODULE SETTING OVERRIDES IT, exactly as it does for Multitenant and SoftDelete above --
	// see BaseComponent.Initialize, where all three read the module configuration first and fall
	// back to this. That is what lets one entity be instantiated per connection by several module
	// instances, each in its own namespace, without the entity file having to know about any of
	// them.
	Namespace string
	// SoftDelete makes Delete mark the record rather than remove it: the data service updates
	// the entity's Deleted field to true, and every read excludes it. Every generated entity
	// already embeds data.DeletionInfo unconditionally, so the field exists in storage whether
	// or not this is set — turning it on is a configuration change and never a migration.
	SoftDelete bool
}

// Storable is the object a data service persists. Implementations are GENERATED, not hand-written:
// the codegen emits ReadAll/WriteAll plus Config for each entity and inherits everything else from
// the data.StorageInfo struct it embeds. Read data/storageinfo.go for the behaviour behind most of
// these methods; the entity's own file overrides only Config.
//
// Two invariants underpin the whole interface, and both are established by the object factory
// rather than by any constructor you can call:
//
//   - the object must be created through ctx.CreateObject("plugin.Entity"), never as a struct
//     literal. The factory runs Constructor and then SetSelfReference
//     (laatooserver/src/core/objectadapter.go:95-100);
//   - selfRef, set by that second call, is what GetLabel and GetObjectRef reflect over. On a
//     struct literal it is nil and BOTH PANIC on an unchecked type assertion
//     (storageinfo.go:68, :112) — not an error return, a panic in the request path.
type Storable interface {
	// Constructor assigns the record an identity. StorageInfo's implementation calls
	// ctx.CreateUUID() and only when Id is still empty (storageinfo.go:51-55), so an id already
	// set — by a caller, or by a decode — survives.
	//
	// The object factory calls it, immediately followed by SetSelfReference, and only for types
	// that satisfy both (objectadapter.go:95-100, :222-225). Do not call it yourself: calling it
	// without also calling SetSelfReference produces an object with an id whose GetLabel and
	// GetObjectRef panic.
	Constructor(ctx.Context)

	// Config returns the entity's storage declaration — object type, collection, label field, and
	// the flags (Trackable, Multitenant, SoftDelete, PreSave/PostSave/PostLoad, Cacheable) the
	// data component reads at Initialize to decide which hooks and scoping apply
	// (datastore/dev/plugins/common/basecomponent.go:65-122).
	//
	// The codegen emits an override per entity. StorageInfo's own implementation returns NIL
	// (storageinfo.go:107-109), so a hand-written Storable that forgets to override it compiles,
	// registers, and then panics inside GetObjectRef, which dereferences the result with no nil
	// check (storageinfo.go:113-114). GetLabel is the only method that guards against it.
	Config() *StorableConfig

	// GetId returns the record's primary key — the value stored as "Id" on every provider.
	GetId() string

	// SetId overwrites the primary key in memory. Put calls it to force the record onto the id the
	// caller named (mongodatabase mongodataservice.go:286) before the upsert. It performs no
	// validation and does not move an already-persisted row: calling it and then saving writes a
	// SECOND record and leaves the original.
	SetId(string)

	// GetLabel returns the record's human-readable label, read by reflection from the field named
	// by Config().LabelField (storageinfo.go:67-76).
	//
	// It never returns an error and never reports a misconfiguration. An entity with no LabelField
	// returns "". A LabelField naming a field that does not exist returns the literal
	// "<invalid Value>", and one naming a non-string field returns "<int Value>" or similar,
	// because reflect.Value.String() substitutes those rather than panicking. It DOES panic if the
	// object was not created through the factory — see the note on the interface.
	GetLabel() string

	// GetVersion returns StorageInfo.Version.
	//
	// Nothing in the SDK, the server or any datastore module ever assigns that field, and
	// StorageInfo.WriteAll does not serialise it (storageinfo.go:125-131) — only Id. In practice
	// it is always the empty string, and the only consumer is GetObjectRef, which copies it into
	// the StorableRef it builds. Treat it as reserved rather than as a concurrency token.
	GetVersion() string

	// SetValues applies a field/value map onto obj — the path every provider's update-by-map takes
	// after loading the record (mongodatabase mongodataservice.go:434, sqldatabase
	// sqldataservice.go:481, boltdatabase kvdataservice.go:699).
	//
	// It SILENTLY DROPS six keys before applying anything: Id, IsNew, CreatedBy, UpdatedBy,
	// CreatedAt and UpdatedAt (storageinfo.go:92-97). A form or API payload trying to set identity
	// or audit fields is ignored with no error — deliberate, since audit stamps belong to
	// data.Track — but indistinguishable from a value that was applied.
	//
	// It also MUTATES the caller's map in place, via delete on the argument, so a map reused after
	// the call has lost those keys.
	//
	// obj is the target the fields are written to and is normally the entity itself. Passing
	// anything else writes into that object instead, with no complaint: jsonbdatabase passes
	// stor.GetObjectRef() (jsonbdataservice.go:549), a freshly built *data.StorableRef, so the
	// values land on a throwaway struct and the entity it then Puts is unchanged. That looks like
	// a bug in that provider, not a supported use.
	SetValues(ctx.Context, interface{}, utils.StringMap) error

	// PreSave runs before the record is written, giving the entity a last chance to derive fields
	// or reject the save by returning an error — which aborts the write on every provider.
	//
	// It runs ONLY when the component has presave enabled (StorableConfig.PreSave, overridable by
	// the module's presave setting), and it is paired with a synchronous CONF_PRESAVE_MSG message
	// that rules subscribe to. StorageInfo's default returns nil (storageinfo.go:82-84), so an
	// entity that does not override it has no hook regardless of the flag.
	PreSave(ctx ctx.Context) error

	// PostSave runs after a successful write, alongside the synchronous CONF_NEWOBJ_MSG /
	// CONF_POSTSAVE_MSG message, and only when the component has postsave enabled.
	//
	// It runs AFTER the row is committed, so returning an error here reports a failure to the
	// caller without undoing the write — unless the whole call sits inside Transaction. Default is
	// a nil return (storageinfo.go:85-87).
	PostSave(ctx ctx.Context) error

	// PostLoad runs on each record after it is decoded from storage, when the component has
	// postload enabled.
	//
	// ITS ERROR IS DISCARDED on the read path: mongodatabase calls stor.PostLoad(ctx) inside a
	// bare loop with no assignment (mongodataservice_get.go:246-250), so a PostLoad that fails
	// still yields the record to the caller as if it had succeeded. Do not use it to enforce an
	// invariant. Default is a nil return (storageinfo.go:88-90).
	PostLoad(ctx ctx.Context) error

	// IsMultitenant is INERT. StorageInfo returns false unconditionally (storageinfo.go:101-103),
	// the codegen emits no override for any entity, and no production code calls it — the only
	// implementations outside the SDK are test stubs in laatooserver/src/testutils.
	//
	// Tenancy is decided elsewhere and never consults this: the data component reads
	// StorableConfig.Multitenant (or the module's multitenant setting) at Initialize
	// (basecomponent.go:106-111) and stamps or filters from ctx.GetUser().GetTenant(). Do not
	// branch on this method — it answers false on a fully tenant-scoped entity.
	IsMultitenant() bool

	// Join is INERT in the same way. StorageInfo's implementation is an empty body
	// (storageinfo.go:105-106) and no generated entity overrides it, so merging a looked-up record
	// into this one does nothing.
	//
	// Its only two callers are the dataadapter join and lookupstorable services
	// (join.go:81, lookupstorable.go:84), and both then respond with the LOOKUP MAP rather than
	// the records Join was supposed to enrich (join.go:87), so even a working implementation would
	// not reach the client through them.
	Join(item Storable)

	// GetObjectRef returns a *data.StorableRef pointing at this record — {Id, Type, Name, Version}
	// — the form a storableref field stores instead of an embedded copy (storageinfo.go:111-115).
	//
	// It is declared as interface{} but every implementation returns *data.StorableRef, and
	// callers assert that directly (datapipeline objectmapprocessor maptoobjectprocessor.go:140).
	//
	// It PANICS on an object not built through the factory (nil selfRef) and on one whose Config
	// is StorageInfo's nil-returning default — neither is guarded. The returned ref's Entity field
	// is nil: it is populated later, by the dataset Expand machinery.
	GetObjectRef() interface{}
}
