package search

import "time"

type Searchable interface {
	GetId() string
	GetType() string
}

type BaseSearchDocument struct {
	Title       string
	Id          string
	Type        string
	TextContent string
	Date        time.Time
}

func (bs *BaseSearchDocument) GetId() string {
	return bs.Id
}

func (bs *BaseSearchDocument) GetType() string {
	return bs.Type
}
