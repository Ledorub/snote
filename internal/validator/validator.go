package validator

type number interface {
	~int | ~uint | ~int16 | ~uint16 | ~int32 | ~uint32 | ~int64 | ~uint64
}

func ValidateIntInRange[T number](v, low, high T) bool {
	return v >= low && v <= high
}
