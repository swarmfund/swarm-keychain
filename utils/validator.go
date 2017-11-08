package utils

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type ValidateError struct {
	Name   string
	Reason error
}

func castError(fieldPrefix string, errors govalidator.Errors) (bool, *ValidateError) {
	// currently our API schema allows to return only one invalid field,
	// so we are not breaking things although we can report multiple fields
	if firstErr, ok := errors[0].(govalidator.Error); ok {
		if fieldPrefix != "" {
			firstErr.Name = fmt.Sprintf("%s.%s", fieldPrefix, firstErr.Name)
		}
		return false, &ValidateError{Name: firstErr.Name, Reason: firstErr.Err}
	}
	return false, nil
}

func ValidateStruct(fieldPrefix string, s interface{}) (bool, *ValidateError) {
	_, err := govalidator.ValidateStruct(s)
	if err != nil {
		if verr, ok := err.(govalidator.Errors); ok {
			errors := verr.Errors()
			if len(errors) == 0 {
				// shouldn't really happen, just in case
				return false, nil
			}
			switch verr[0].(type) {
			case govalidator.Errors:
				{
					if errs, ok := verr[0].(govalidator.Errors); ok {
						return castError(fieldPrefix, errs)
					}
				}
			case govalidator.Error:
				{
					return castError(fieldPrefix, errors)
				}
			}
		}
		return false, nil
	}
	return true, nil
}
