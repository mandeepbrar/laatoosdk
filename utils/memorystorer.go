package utils

import (
	"fmt"
	"github.com/pmylund/go-cache"
)

//Memory storer class
type MemoryStorer struct {
	*cache.Cache
}

//creates a new storer for memory
func NewMemoryStorer() *MemoryStorer {
	memStor := &MemoryStorer{cache.New(cache.NoExpiration, 0)}
	return memStor
}

//Puts an item in memory
func (ms *MemoryStorer) PutObject(id string, item interface{}) error {
	ms.Set(id, item, cache.NoExpiration)
	return nil
}

//Get an item from memory
func (ms *MemoryStorer) GetObject(id string) (interface{}, error) {
	item, ok := ms.Get(id)
	if !ok {
		return nil, fmt.Errorf("Object %s not found", id)
	}
	return item, nil
}

//Delete an item from memory
func (ms *MemoryStorer) DeleteObject(id string) error {
	ms.Delete(id)
	return nil
}

//Get a list of all items in memory
func (ms *MemoryStorer) GetList() []interface{} {
	cItems := ms.Items()
	items := make([]interface{}, len(cItems))
	idx := 0
	for _, value := range cItems {
		items[idx] = value.Object
		idx++
	}
	return items
}
