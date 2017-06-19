package data

import (
	"fmt"
	"laatoo/sdk/errors"
)

const (
	DATA_ERROR_CONNECTION      = "Data_Error_Connection"
	DATA_ERROR_OPERATION       = "Data_Error_Operation"
	DATA_ERROR_NOT_IMPLEMENTED = "Data_Error_Not_Implemented"
	DATA_ERROR_ID_NOT_FOUND    = "Data_Error_ID_Not_Found"
)

func init() {
	errors.RegisterCode(DATA_ERROR_CONNECTION, errors.ERROR, fmt.Errorf("Could not connect to the database."))
	errors.RegisterCode(DATA_ERROR_NOT_IMPLEMENTED, errors.ERROR, fmt.Errorf("Method not implemented for the service."))
	errors.RegisterCode(DATA_ERROR_OPERATION, errors.ERROR, fmt.Errorf("Error occured while executing a database operation."))
	errors.RegisterCode(DATA_ERROR_ID_NOT_FOUND, errors.ERROR, fmt.Errorf("Id not provided for the entity."))
}
