package validator

import "cmp"

func ValidateValueInRange[T cmp.Ordered](v, low, high T) bool {
	return v >= low && v <= high
}
