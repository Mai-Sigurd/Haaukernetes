package models

// swagger:model CommonError
type CommonError struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model CommonSuccess
type CommonSuccess struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// ErrHandler returns error message response
func ErrHandler(errormessage string) *CommonError {
	response := CommonError{}
	response.Status = 0
	response.Message = errormessage
	return &response
}
