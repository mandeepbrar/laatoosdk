package data

import (
	"laatoo.io/sdk/server/errors"
)

const (
	DATA_ERROR_CONNECTION      = "Data_Error_Connection"
	DATA_ERROR_OPERATION       = "Data_Error_Operation"
	DATA_ERROR_NOT_IMPLEMENTED = "Data_Error_Not_Implemented"
	DATA_ERROR_ID_NOT_FOUND    = "Data_Error_ID_Not_Found"
)

func init() {
	errors.RegisterCode(DATA_ERROR_CONNECTION, "Could not connect to the database.")
	errors.RegisterCode(DATA_ERROR_NOT_IMPLEMENTED, "Method not implemented for the service.")
	errors.RegisterCode(DATA_ERROR_OPERATION, "Error occured while executing a database operation.")
	errors.RegisterCode(DATA_ERROR_ID_NOT_FOUND, "Id not provided for the entity.")
}
