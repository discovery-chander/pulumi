package repository

import (
	stderrors "errors"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/jinzhu/gorm"
)

var (
	// ErrEntityNotFound is returned when the specified record is not found.
	ErrEntityNotFound = stderrors.New("Entity Not Found")
)

// EvaluateError determines if the error should be recognized as an ErrEntityNotFound.
func EvaluateError(err error) error {
	if gorm.IsRecordNotFoundError(err) {
		return errors.Wrap(ErrEntityNotFound, "did not find entity")
	}
	if err != nil {
		return errors.Wrap(err, "unable to retrieve entity")
	}

	return nil
}
