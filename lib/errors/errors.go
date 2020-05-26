// errors is meant to be a drop-in replacement for github.com/pkg/errors.
package errors

import (
	"fmt"

	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	"github.com/pkg/errors"
)

func Wrap(err error, message string) error {
	logErrorf("%v: %+v", message, err)
	return errors.Wrap(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	logErrorf("%v: %+v", message, err)
	return errors.Wrapf(err, format, args...)
}

func WithStack(err error) error {
	return errors.WithStack(err)
}

func WithMessage(err error, message string) error {
	logErrorf("%v: %v", message, err)
	return errors.WithMessage(err, message)
}

func WithMessagef(err error, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	logErrorf("%v: %v", message, err)
	return errors.WithMessagef(err, format, args...)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func New(message string) error {
	err := errors.New(message)
	logErrorf("%+v", err)
	return err
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Cause(err error) error {
	return errors.Cause(err)
}

func Errorf(format string, args ...interface{}) error {
	err := errors.Errorf(format, args...)
	logErrorf("%+v", err)
	return err
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func logErrorf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
