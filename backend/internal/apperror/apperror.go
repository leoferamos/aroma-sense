package apperror

import "errors"

// DomainError carries a stable code/message for business errors without binding to HTTP transport.
type DomainError struct {
	Err     error
	Code    string
	Message string
}

func (d *DomainError) Error() string {
	if d == nil {
		return ""
	}
	if d.Err != nil {
		return d.Err.Error()
	}
	return d.Message
}

func (d *DomainError) Unwrap() error { return d.Err }

// NewDomain creates a DomainError with a code/message, leaving transport (HTTP status) to callers.
func NewDomain(err error, code, message string) *DomainError {
	return &DomainError{Err: err, Code: code, Message: message}
}

// NewCodeMessage builds a DomainError from code/message without needing a prebuilt error value.
func NewCodeMessage(code, message string) *DomainError {
	return &DomainError{Err: errors.New(message), Code: code, Message: message}
}
