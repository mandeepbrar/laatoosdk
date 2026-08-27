package search

import "time"

// Searchable is a document a search provider can index, retrieve and delete.
//
// The interface exposes identity only — there is no field accessor — so each provider decides for
// itself what text of the concrete document it indexes: bleve and googlesearch hand the struct
// straight to the engine, while sonicsearch reflects over its exported string fields
// (laatoomodules/search/dev/plugins/sonicsearch/src/server/go/sonicsearchservice.go:152-183).
// Anything a query must be able to match therefore has to be an exported string field on the
// concrete type; unexported and non-string data is not searchable anywhere.
//
// A Searchable coming BACK from Search is never the caller's own type. Every provider returns
// *BaseSearchDocument: bleve rebuilds one by reflection (blevesearchservice.go:150-164),
// googlesearch decodes into one (googlesearchservice.go:110-120), and sonicsearch, which stores no
// document content at all, returns one holding just an id and type
// (sonicsearchservice.go:249-262). Asserting a result back to the indexed type panics — read the
// record from its system of record using GetId.
type Searchable interface {
	// GetId returns the document's identifier within its type.
	//
	// It is half of the storage key, never the whole one — see GetType.
	GetId() string

	// GetType returns the document's type name, which combines with GetId to form the key every
	// provider stores under: fmt.Sprintf("%s_%s", GetType(), GetId()) in bleve
	// (blevesearchservice.go:83-94 and Delete at :166-174) and googlesearch
	// (googlesearchservice.go:63-76, :128-137), and the same pairing in sonicsearch
	// (sonicsearchservice.go:186-204, :266-274).
	//
	// Both halves must be non-empty and stable. An empty type still yields a well-formed key
	// ("_<id>"), so every empty-type document sharing an id silently overwrites the others, and a
	// Delete issued with a different type than the Index used removes nothing and reports success.
	GetType() string

	// GetTenant returns the tenant the document belongs to — and NO PROVIDER READS IT.
	//
	// Implementing it does not scope anything: it takes no part in the storage key, and
	// sonicsearch deliberately excludes the Tenant field from the indexed text so that a query
	// cannot match on it (identityFields, sonicsearchservice.go:140-148). Search isolation is the
	// `bucket` argument passed to Index/Search/Delete, which the caller supplies; a caller that
	// passes the same bucket for two tenants gets one shared index, and this accessor will not
	// have prevented it.
	//
	// It survives only as a plain field on a stored document, for a caller that chooses to read
	// and filter on it after the fact.
	GetTenant() string
}

// BaseSearchDocument is the fixed, provider-neutral document shape — and the concrete type every
// provider returns from Search regardless of what was indexed.
//
// The numbered Text and Date fields are deliberately generic slots: the indexed text is taken from
// exported string fields, so mapping meaningful content onto Text1..Text6 is what makes it
// findable. The Date fields are not: bleve's result reconstruction copies string-kind fields only
// (blevesearchservice.go:154-162), so Date1..Date3 come back zero on every bleve search hit even
// when they were indexed.
type BaseSearchDocument struct {
	Title  string
	Id     string
	Type   string
	Text1  string
	Text2  string
	Text3  string
	Text4  string
	Text5  string
	Text6  string
	User   string
	UserId string
	Tenant string
	Date1  time.Time
	Date2  time.Time
	Date3  time.Time
}

func (bs *BaseSearchDocument) GetId() string {
	return bs.Id
}

func (bs *BaseSearchDocument) GetType() string {
	return bs.Type
}

func (bs *BaseSearchDocument) GetTenant() string {
	return bs.Tenant
}
