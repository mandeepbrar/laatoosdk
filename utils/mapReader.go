package utils

type MapReader struct {
	MapToRead StringMap
	Mappings  StringsMap
}

func (rdr *MapReader) ReadStringFromMap(key string) (string, bool) {
	mappedKey, ok := rdr.Mappings["Id"]
	if ok {
		return rdr.MapToRead.GetString(mappedKey)
	}
	return "", false
}
