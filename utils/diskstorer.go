package utils

import (
	"github.com/peterbourgon/diskv"
)

//Memory storer class
type DiskStorer struct {
	*diskv.Diskv
}

//creates a new storer for memory
func NewDiskStorer(path string, maxSize uint64) *DiskStorer {
	flatTransform := func(s string) []string { return []string{} }
	dv := diskv.New(diskv.Options{
		BasePath:     path,
		Transform:    flatTransform,
		CacheSizeMax: maxSize,
	})
	dStor := &DiskStorer{dv}
	return dStor
}

//Puts an item in memory
func (ds *DiskStorer) PutObject(id string, item []byte) error {
	return ds.Write(id, item)
}

//Get an item from memory
func (ds *DiskStorer) GetObject(id string) ([]byte, error) {
	return ds.Read(id)
}

//Delete an item from memory
func (ds *DiskStorer) DeleteObject(id string) error {
	return ds.Erase(id)
}

// Keys returns a channel that will yield every key accessible by the store,
// in undefined order. If a cancel channel is provided, closing it will
// terminate and close the keys channel.
func (ds *DiskStorer) Keys(cancel <-chan struct{}) <-chan string {
	return ds.KeysPrefix("", cancel)
}
