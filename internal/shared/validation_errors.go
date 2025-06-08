// internal/common/errors/validation_errors.go
package errors

import "strings"

type ValidationErrors struct {
	Errors []string
}

func (v *ValidationErrors) Add(entity, msg string) {
	v.Errors = append(v.Errors, entity+":"+msg)
}

func (v *ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

func (v *ValidationErrors) Error() string {
	return "validation failed: " + strings.Join(v.Errors, "; ")
}
