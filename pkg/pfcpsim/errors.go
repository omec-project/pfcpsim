package pfcpsim

import (
	"fmt"
	"strings"
)

type pfcpSimError struct {
	message string
	error   []error
}

func (e *pfcpSimError) unwrap() string {
	errMsg := strings.Builder{}
	errMsg.WriteString("")

	for i, e := range e.error {
		errMsg.WriteString(fmt.Sprintf("\n\t- Error %v: %v", i, e))
	}

	return errMsg.String()
}

func (e *pfcpSimError) Error() string {
	return fmt.Sprintf("Message: %v. %v", e.message, e.unwrap())
}

func NewInvalidCauseError(err ...error) *pfcpSimError {
	return &pfcpSimError{
		message: "Invalid Cause from response",
		error:   err,
	}
}

func NewNotEnoughSessionsError(err ...error) *pfcpSimError {
	return &pfcpSimError{
		message: "Not enough active sessions",
		error:   err,
	}
}

func NewAssociationInactiveError(err ...error) *pfcpSimError {
	return &pfcpSimError{
		message: "Could not complete operation: Association is not active",
		error:   err,
	}
}

func NewTimeoutExpiredError(err ...error) *pfcpSimError {
	return &pfcpSimError{
		message: "Timeout has expired",
		error:   err,
	}
}

func NewInvalidResponseError(err ...error) *pfcpSimError {
	return &pfcpSimError{
		message: "Invalid response received",
		error:   err,
	}
}
