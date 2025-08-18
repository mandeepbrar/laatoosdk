package datatypes

import (
	"laatoo.io/sdk/constants"
)

type DataType int

const (
	Stringmap DataType = iota
	Stringsmap
	Map
	Datetime
	Config
	Configarr
	Bytes
	Files
	Int
	String
	Stringarr
	Bool
	Object
	None
)

func (t DataType) String() string {
	switch t {
	case Stringmap:
		return constants.OBJECTTYPE_STRINGMAP
	case Stringsmap:
		return constants.OBJECTTYPE_STRINGSMAP
	case Bytes:
		return constants.OBJECTTYPE_BYTES
	case String:
		return constants.OBJECTTYPE_STRING
	case Stringarr:
		return constants.OBJECTTYPE_STRINGARR
	case Bool:
		return constants.OBJECTTYPE_BOOL
	case Files:
		return constants.OBJECTTYPE_FILES
	case Datetime:
		return constants.OBJECTTYPE_DATETIME
	case Int:
		return constants.OBJECTTYPE_INT
	case Config:
		return constants.OBJECTTYPE_CONFIG
	case Configarr:
		return constants.OBJECTTYPE_CONFIGARR
	case Map:
		return constants.OBJECTTYPE_MAP
	case Object:
		return constants.OBJECTTYPE_OBJECT
	}
	return "None"
}

func ConvertDataType(dtype string) DataType {
	switch dtype {
	case "":
		return None
	case constants.OBJECTTYPE_STRINGMAP:
		return Stringmap
	case constants.OBJECTTYPE_STRINGSMAP:
		return Stringsmap
	case constants.OBJECTTYPE_BYTES:
		return Bytes
	case constants.OBJECTTYPE_STRING:
		return String
	case constants.OBJECTTYPE_STRINGARR:
		return Stringarr
	case constants.OBJECTTYPE_BOOL:
		return Bool
	case constants.OBJECTTYPE_FILES:
		return Files
	case constants.OBJECTTYPE_DATETIME:
		return Datetime
	case constants.OBJECTTYPE_INT:
		return Int
	case constants.OBJECTTYPE_CONFIG:
		return Config
	case constants.OBJECTTYPE_CONFIGARR:
		return Configarr
	case constants.OBJECTTYPE_MAP:
		return Map
	default:
		return Object
	}
}
