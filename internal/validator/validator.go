package validator

import "cmp"

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) CheckIsValid() bool {
	return len(v.NonFieldErrors) == 0 && len(v.FieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[field]; !exists {
		v.FieldErrors[field] = message
	}
}

func (v *Validator) CheckField(field string, ok bool, message string) {
	if !ok {
		v.AddFieldError(field, message)
	}
}

func (v *Validator) GetNonFieldErrors() []string {
	return v.NonFieldErrors
}

func (v *Validator) GetFieldErrors() map[string]string {
	return v.FieldErrors
}

func New() *Validator {
	return &Validator{
		FieldErrors: make(map[string]string),
	}
}

func ValidateValueInRange[T cmp.Ordered](v, low, high T) bool {
	return v >= low && v <= high
}
