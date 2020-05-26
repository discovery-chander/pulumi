package errors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQuickCoverage(t *testing.T) {
	err := New("my error")
	require.Error(t, Wrap(err, "wrapped"))
	require.Error(t, Wrapf(err, "wrapped with value %v", 1))
	require.NoError(t, Unwrap(err))
	require.Error(t, WithStack(err))
	require.Error(t, WithMessage(err, "just this message"))
	require.Error(t, WithMessagef(err, "just this message %v", 1))
	require.Error(t, Errorf("A new error with decimal: %v", 3.14159))
	require.True(t, Is(err, err))
	var errAs error
	require.True(t, As(err, &errAs))
	require.EqualError(t, Cause(err), err.Error())
}
