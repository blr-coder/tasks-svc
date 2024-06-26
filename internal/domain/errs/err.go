package errs

import (
	"errors"
	"fmt"
	"reflect"
)

type DomainError struct {
	Message string `json:"message"`
	Params  Params `json:"params"`
}

type Params map[string]string

func (e *DomainError) Error() string {
	return fmt.Sprintf("%s. Params: %v", e.Message, e.Params)
}

func (e *DomainError) Is(tgt error) bool {
	var target *DomainError
	if !errors.As(tgt, &target) {
		return false
	}
	return reflect.DeepEqual(e, target)
}

func (e *DomainError) AddParam(key string, value string) {
	e.Params[key] = value
}

func (e *DomainError) WithParam(key, value string) *DomainError {
	e.AddParam(key, value)
	return e
}
func (e *DomainError) WithParams(params map[string]string) *DomainError {
	for key, value := range params {
		e.AddParam(key, value)
	}
	return e
}

type DomainNotFoundError struct {
	*DomainError
}

func NewDomainNotFoundError() *DomainError {
	return &DomainError{
		Message: fmt.Sprintf("Record not found"),
		Params:  map[string]string{},
	}
}
