package db

import (
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/require"
)

func TestOpenDBConnection(t *testing.T) {
	t.Run("Should return a db when OpenDBConnection is called", func(t *testing.T) {
		mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
		mocket.Catcher.Logging = true
		mocket.Catcher.Reset()
		db, err := OpenDBConnection("CONNECTION_STRING", mocket.DriverName)
		require.NotEmpty(t, db, "The gorm db instance was returned")
		require.NoError(t, errors.Cause(err), "There were no errors returned")
	})
	t.Run("Should throw error when driver fails", func(t *testing.T) {
		mocket.Catcher.NewMock()
		mocket.Catcher.Reset()
		db, err := OpenDBConnection("CONNECTION_STRING", "")
		require.Nil(t, db)
		require.EqualError(t, errors.Cause(err), `missing "=" after "CONNECTION_STRING" in connection info string"`)
	})
	t.Run("Should fail when ConnectionString is not set", func(t *testing.T) {
		db, err := OpenDBConnection("", mocket.DriverName)
		require.Nil(t, db)
		require.EqualError(t, errors.Cause(err), "Missing arguments to open DB connection")
	})
}
