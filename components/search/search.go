package search

import "laatoo/sdk/core"

type SearchComponent interface {
	//Index a searchable document
	Index(ctx core.RequestContext, s Searchable) error
	//Index a searchable document
	Search(ctx core.RequestContext, query string) ([]Searchable, error)
	//Delete a searchable document
	Delete(ctx core.RequestContext, s Searchable) error
}
