package search

import "laatoo/sdk/core"

type SearchComponent interface {
	//Index a searchable document
	Index(ctx core.RequestContext, s Searchable) error
	//Update index
	UpdateIndex(ctx core.RequestContext, id string, stype string, u map[string]interface{}) error
	//Index a searchable document
	Search(ctx core.RequestContext, query string) ([]Searchable, error)
	//Delete a searchable document
	Delete(ctx core.RequestContext, s Searchable) error
}
