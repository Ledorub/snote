package common

func MergeValidationErrors(v Validator) ValidationErrors {
	nonFieldErrors := v.GetNonFieldErrors()
	fieldErrors := v.GetFieldErrors()
	list := make(ValidationErrors, len(nonFieldErrors)+len(fieldErrors))

	i := 0
	for _, msg := range nonFieldErrors {
		item := map[string]string{"message": msg}
		list[i] = item
		i++
	}
	for field, msg := range fieldErrors {
		item := map[string]string{
			"field":   field,
			"message": msg,
		}
		list[i] = item
		i++
	}
	return list
}
