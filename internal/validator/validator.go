package validator

import (
	"cmp"
	"time"
)

type ValidationError string

func (e ValidationError) Error() string {
	return string(e)
}

type Validator struct {
	Errors []ValidationError
}

func (v *Validator) CheckIsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(message string) {
	v.Errors = append(v.Errors, ValidationError(message))
}

func (v *Validator) Check(ok bool, message string) {
	if !ok {
		v.AddError(message)
	}
}

func (v *Validator) GetErrors() []ValidationError {
	return v.Errors
}

func New() *Validator {
	return &Validator{}
}

// ValidateValueInRange checks that low <= v <= high.
func ValidateValueInRange[T cmp.Ordered](v, low, high T) bool {
	return v >= low && v <= high
}

// ValidateASCIIAlphaNumeric checks that s is in A-Za-z0-9.
func ValidateASCIIAlphaNumeric(s string) bool {
	for _, char := range s {
		isLetter := char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z'
		isDigit := char >= '0' && char <= '9'
		if !(isLetter || isDigit) {
			return false
		}
	}
	return true
}

// ValidateTimeInRange checks that low < t <= high.
func ValidateTimeInRange(t, low, high time.Time) bool {
	return t.After(low) && t.Before(high) || t.Equal(high)
}
