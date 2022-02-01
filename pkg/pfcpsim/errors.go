package pfcpsim

import "fmt"

type pfcpSimError struct {
	message string
	error   error
}

func (e *pfcpSimError) Error() string {
	return fmt.Sprintf("Message: %v, Error: %v", e.message, e.error)
}

func NewInvalidCauseError(err error) *pfcpSimError {
	return &pfcpSimError{
		message: "Invalid Cause from response",
		error:   err,
	}
}

func NewAssociationInactiveError(err error) *pfcpSimError {
	return &pfcpSimError{
		message: "Could not complete operation. Association is not active",
		error:   err,
	}
}

func NewTimeoutExpiredError(err error) *pfcpSimError {
	return &pfcpSimError{
		message: "Timeout has expired",
		error:   err,
	}
}

func NewInvalidResponseError(err error) *pfcpSimError {
	return &pfcpSimError{
		message: "Invalid response received",
		error:   err,
	}
}
