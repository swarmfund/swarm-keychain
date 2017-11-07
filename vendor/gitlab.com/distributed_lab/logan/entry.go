package logan

import (
	"fmt"

	goerrors "github.com/go-errors/errors"
	"github.com/sirupsen/logrus"
	"io"
)

var (
	ErrorKey = logrus.ErrorKey
	StackKey = "stack"
)

type Entry struct {
	*logrus.Entry
}

func (e *Entry) SetLevel(level Level) *Entry {
	e.Entry.Level = logrus.Level(level)
	return e
}

func (e *Entry) SetLoggerOut(writer io.Writer) *Entry {
	e.Entry.Logger.Out = writer
	return e
}

// WithRecover creates `go-errors.Error` error from the `recoverData`
// and returns Entry with this error and its stack.
func (e *Entry) WithRecover(recoverData interface{}) *Entry {
	err := toError(recoverData)
	goErr := goerrors.Wrap(err, 4)
	return e.WithStack(goErr).WithError(goErr)
}

func (e *Entry) WithError(err error) *Entry {
	fielderError, ok := err.(FieldedErrorI)
	if ok {
		e = &Entry{
			e.Entry.WithFields(logrus.Fields(fielderError.Fields())),
		}
	}

	return &Entry{
		Entry: e.Entry.WithError(err),
	}
}

func (e *Entry) WithField(key string, value interface{}) *Entry {
	fieldedEntity, ok := value.(FieldedEntityI)

	if ok {
		return e.WithFields(obtainFields(key, fieldedEntity))
	} else {
		// It's just a plain field.
		return &Entry{
			e.Entry.WithField(key, value),
		}
	}
}

func (e *Entry) WithFields(fields F) *Entry {
	return &Entry{e.Entry.WithFields(logrus.Fields(fields))}
}

func (e *Entry) WithStack(err error) *Entry {
	return e.WithField(StackKey, getStack(err))
}

// Debugf logs a message at the debug severity.
func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Entry.Debugf(format, args...)
}

// Debug logs a message at the debug severity.
func (e *Entry) Debug(args ...interface{}) {
	e.Entry.Debug(args...)
}

// Infof logs a message at the Info severity.
func (e *Entry) Infof(format string, args ...interface{}) {
	e.Entry.Infof(format, args...)
}

// Info logs a message at the Info severity.
func (e *Entry) Info(args ...interface{}) {
	e.Entry.Info(args...)
}

// Warnf logs a message at the Warn severity.
func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Entry.Warnf(format, args...)
}

// Warn logs a message at the Warn severity.
func (e *Entry) Warn(args ...interface{}) {
	e.Entry.Warn(args...)
}

// Errorf logs a message at the Error severity.
func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Entry.Errorf(format, args...)
}

// Error logs a message at the Error severity.
func (e *Entry) Error(args ...interface{}) {
	e.Entry.Error(args...)
}

// Panicf logs a message at the Panic severity.
func (e *Entry) Panicf(format string, args ...interface{}) {
	e.Entry.Panicf(format, args...)
}

// Panic logs a message at the Panic severity.
func (e *Entry) Panic(args ...interface{}) {
	e.Entry.Panic(args...)
}

// getStack returns the stack, as a string, if one can be extracted from `err`.
// If provided `err` is NOT of go-errors.Error type - returns "unknown".
func getStack(err error) string {
	if stackProvider, ok := err.(Stackable); ok {
		return string(stackProvider.Stack())
	}

	return "unknown"
}

// toError creates error from the string representation of `data`, if `data` is not error itself.
func toError(data interface{}) error {
	err, ok := data.(error)
	if !ok {
		err = fmt.Errorf("%s", data)
	}

	return err
}
