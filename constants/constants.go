package constants

const (
	AUTHORIZATIONHANDLER = "authorizationhandler"
	ALLOWALLSERVICES     = "allowallservices"
	ACCESSPERMISSION     = "accesspermission"
	BARREDROLES          = "disallowedroles"
	JWTSECRET            = "__jwtsecret"
	AUTHHEADER           = "__authheader"
	USER                 = "__user"
	ROLE                 = "__role"
	SYSTEMROLE           = "__system"
	TENANT               = "__tenant"
	AUTHCLAIMS           = "__authclaims"
	CONF_PVTKEYPATH      = "pvtkey"
	CONF_PUBLICKEYPATH   = "publickey"
	ENCODING             = "encoding"
	BASEDIR              = "basedir"
	CONFIGDIR            = "configdir"
	MODULEDIR            = "moduledir"
	HTTP_METHOD          = "httpmethod"
)

type HttpChannelMethod int

const (
	GET HttpChannelMethod = iota
	PUT
	POST
	DELETE
	NONE
)

const (
	OBJECTTYPE_STRINGMAP  = "stringmap"
	OBJECTTYPE_STRINGSMAP = "stringsmap"
	OBJECTTYPE_MAP        = "map"
	OBJECTTYPE_STRINGARR  = "stringarray"
	OBJECTTYPE_BYTES      = "bytes"
	OBJECTTYPE_FILES      = "files"
	OBJECTTYPE_STRING     = "string"
	OBJECTTYPE_DATETIME   = "datetime"
	OBJECTTYPE_CONFIG     = "config"
	OBJECTTYPE_CONFIGARR  = "configarray"
	OBJECTTYPE_BOOL       = "bool"
	OBJECTTYPE_INT        = "int"
	OBJECTTYPE_OBJECT     = "object"
)
