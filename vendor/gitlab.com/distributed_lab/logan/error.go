package logan

import (
	"github.com/pkg/errors"
	"fmt"
)

// FromPanic extracts the err from the result of a recover() call.
func FromPanic(rec interface{}) error {
	err, ok := rec.(error)
	if !ok {
		err = fmt.Errorf("%s", rec)
	}

	return err
}

type FieldedErrorI interface {
	Error() string
	Fields() F
	WithField(key string, value interface{}) FieldedErrorI
	WithFields(fields F) FieldedErrorI
}

type Stackable interface {
	Stack() []byte
}

// If base is nil, Wrap returns nil.
func Wrap(base error, msg string) FieldedErrorI {
	if base == nil {
		return nil
	}

	fieldedError, ok := base.(*FError)
	if !ok {
		fieldedError = &FError{
			err:    base,
			fields: F{},
		}
	}

	fieldedError.err = errors.Wrap(fieldedError.err, msg)
	return fieldedError
}

func NewError(msg string) FieldedErrorI {
	return &FError{
		err:    errors.New(msg),
		fields: F{},
	}
}

type FError struct {
	err    error
	fields F
}

func (e *FError) Error() string {
	return e.err.Error()
}

func (e *FError) Fields() F {
	return e.fields
}

// WithField returns the same instance
func (e *FError) WithField(key string, value interface{}) FieldedErrorI {
	if e == nil {
		return nil
	}

	fieldedEntity, ok := value.(FieldedEntityI)

	if ok {
		return e.WithFields(obtainFields(key, fieldedEntity))
	} else {
		// It's just a plain field.
		e.fields[key] = value
		return e
	}
}

// WithFields returns the same instance
func (e *FError) WithFields(fields F) FieldedErrorI {
	if e == nil {
		return nil
	}

	for key, value := range fields {
		e.fields[key] = value
	}
	return e
}

func (e *FError) Cause() error {
	return errors.Cause(e.err)
}
