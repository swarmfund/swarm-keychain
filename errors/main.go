package errors

import (
	"fmt"

	"github.com/go-errors/errors"
)

// FromPanic extracts the err from the result of a recover() call.
func FromPanic(rec interface{}) error {
	err, ok := rec.(error)
	if !ok {
		err = fmt.Errorf("%s", rec)
	}

	return errors.Wrap(err, 4)
}

// Stack returns the stack, as a string, if one can be extracted from `err`.
func Stack(err error) string {

	if stackProvider, ok := err.(*errors.Error); ok {
		return string(stackProvider.Stack())
	}

	return "unknown"
}
