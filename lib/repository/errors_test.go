package repository

import (
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
)

func TestEvaluateError(t *testing.T) {
	t.Run("Should return respository.ErrEntityNotFound when passed gorm.ErrRecordNotFound", func(t *testing.T) {
		err := EvaluateError(gorm.ErrRecordNotFound)
		require.Error(t, err)
		require.EqualError(t, errors.Cause(err), ErrEntityNotFound.Error())
	})

	t.Run("Should bubble up passed error if not gorm.ErrRecordNotFound", func(t *testing.T) {
		expectedErr := errors.New("my error")
		err := EvaluateError(expectedErr)
		require.Error(t, err)
		require.EqualError(t, errors.Cause(err), expectedErr.Error())
	})

	t.Run("Should return nil on a nil error", func(t *testing.T) {
		require.NoError(t, EvaluateError(nil))
	})
}
