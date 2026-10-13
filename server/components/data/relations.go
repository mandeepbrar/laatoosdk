package data

import (
	"laatoo.io/sdk/server/core"
)

// A RELATION is one edge between two records: a storableref followed from the record that holds it
// (From) to the record it names (To), typed by the relationshipname the field declares.
//
// Relations answer the question a stored reference cannot answer from its own record: which records
// point AT this one. A forward reference is read off the record that holds it; the reverse direction
// has nothing to read, so it needs an adjacency store -- graphedges.Edge, maintained from every
// entity's storableref fields. A RelationComponent is that store's query surface, and
// DataManager.QueryRelations is how any caller reaches it without knowing the store's schema.
//
// RELATIONS, NOT RECORDS. A query returns the edges, and a caller wanting the far records reads them
// through their own components (GetMulti on the far entity's DataComponent). That is deliberate on
// two counts. A mixed list of records would lose which relation reached each one. And the far
// records' own components are where their tenancy and soft-delete apply -- so an edge whose far
// record has since been soft-deleted or moved to another tenant resolves to nothing when read, rather
// than to a record the caller may not see. The adjacency store is raw adjacency; the components are
// the authority on what is live.

// RelationDirection says which way a relation query follows edges from the records it starts at.
type RelationDirection string

const (
	// RelationIncoming follows edges that END at the starting records: the records that reference
	// them. This is the direction a stored reference cannot answer from its own record.
	RelationIncoming RelationDirection = "incoming"
	// RelationOutgoing follows edges that START at the starting records: what they reference.
	RelationOutgoing RelationDirection = "outgoing"
	// RelationBoth follows edges in both directions.
	RelationBoth RelationDirection = "both"
)

// RelationQuery asks which records are related to Nodes.
type RelationQuery struct {
	// Nodes are the records the query starts from, as references. Each needs Id and Type: the type
	// is part of an edge's identity, because one adjacency store is shared by every entity type on
	// a dataconnection and an id alone can name another type's record. DataConnection on a node is
	// not consulted -- the connection a query runs against is the DataManager call's.
	Nodes []StorableRef
	// Direction is required. Its zero value is REFUSED rather than defaulted: incoming and outgoing
	// answer different questions, and a default would answer one when the caller meant the other.
	Direction RelationDirection
	// Relationship narrows the query to one edge type -- the relationshipname a storableref field
	// declares. Empty means any relationship.
	Relationship string
	// OtherType narrows the far end to one entity type: the source for an incoming query, the
	// target for an outgoing one. Empty means any type.
	OtherType string
	// Limit is the most relations the CALLER can use. A component refuses a query matching more than
	// the smaller of Limit and its own bound, and reads no more than one past that -- so a caller
	// that would refuse a larger answer anyway does not pay for reading up to the component's bound
	// first. Zero or negative means no caller limit: the component's own bound applies.
	Limit int
}

// Relation is one edge: From holds a storableref naming To, declared with Relationship.
//
// From and To carry the Id and Type of each end, and DataConnection where the store recorded it.
// They are references, not records -- see the note at the top of this file.
type Relation struct {
	From         StorableRef
	To           StorableRef
	Relationship string
}

// RelationComponent answers relation queries for ONE dataconnection's adjacency store. graphedges
// implements it and registers one per Edge store it maintains (DataManager.RegisterRelationComponent).
//
// It is an optional interface beside DataComponent and never a method on it, for the reason
// NavigatingComponent is: adding a method to DataComponent breaks its implementors SILENTLY, because
// Go checks satisfaction at the assertion site.
//
// AN EDGE LIVES ON ITS SOURCE'S CONNECTION. An adjacency store co-located with a dataconnection
// holds the edges whose SOURCE records live there, so an incoming query answers only for sources on
// that connection. A caller with sources on several connections asks each one.
type RelationComponent interface {
	// QueryRelations returns the relations the query names, and the total matching.
	//
	// pageNum is 1-based and pageSize/pageNum of -1, -1 returns every relation, as Get does. An
	// empty result is an ANSWER -- nothing is related -- and is returned with a nil error. The
	// adjacency store is read through its own component, so the edges themselves are scoped to the
	// caller's tenant; the far records are not read at all.
	//
	// An implementation refuses a query it cannot answer correctly -- a zero Direction, a node with
	// no Id or Type, a result past its bound or past the query's Limit -- rather than answering part
	// of it.
	QueryRelations(ctx core.RequestContext, query RelationQuery, pageSize int, pageNum int) (relations []Relation, totalrecs int, err error)
}
