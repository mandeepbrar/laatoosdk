package search

import "time"

type Searchable interface {
	GetId() string
	GetType() string
	GetTenant() string
}

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
