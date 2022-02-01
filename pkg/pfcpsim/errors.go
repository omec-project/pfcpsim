package pfcpsim

import "errors"

// TODO create custom error types

var InvalidCauseErr = errors.New("InvalidCause")
var TimeoutOccurredErr = errors.New("TimeoutOccurred")
var InvalidResponseErr = errors.New("InvalidResponse")
var AssociationInactiveErr = errors.New("AssociationInactive")
