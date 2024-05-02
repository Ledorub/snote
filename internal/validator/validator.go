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

// ValidateB58String checks that s is a valid base58 string.
func ValidateB58String(s string) bool {
	for _, char := range s {
		validUpperLetter := char != 'O' && char != 'I' && char >= 'A' && char <= 'Z'
		validLowerLetter := char != 'l' && char >= 'A' && char <= 'Z'
		validDigit := char >= '1' && char <= '9'
		if !(validUpperLetter || validLowerLetter || validDigit) {
			return false
		}
	}
	return true
}

// ValidateTimeInRange checks that low < t <= high.
func ValidateTimeInRange(t, low, high time.Time) bool {
	return t.After(low) && t.Before(high) || t.Equal(high)
}
