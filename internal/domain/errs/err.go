package errs

import "fmt"

/*type DomainError interface {
	Code() string
	Message() string
	Error() string
}*/

type DomainError struct {
	Message string `json:"message"`
	Params  Params `json:"params"`
}

func (e *DomainError) AddParam(key string, value string) {
	e.Params[key] = value
}

func (e *DomainError) Error() string {
	return e.Message
}

type Params map[string]string

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

func NewDomainNotFoundError() *DomainNotFoundError {
	return &DomainNotFoundError{
		DomainError: &DomainError{
			Message: fmt.Sprintf("Record not found"),
			Params:  map[string]string{},
		},
	}
}
