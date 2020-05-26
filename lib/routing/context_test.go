package routing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithLoggingContext(t *testing.T) {
	t.Run("Should add a logging context", func(t *testing.T) {
		expectedContext := LoggingContext{
			"hello": "world",
		}
		ctx := WithLoggingContext(context.Background(), expectedContext)
		require.EqualValues(t, expectedContext, ctx.Value(keyLoggingContext))
	})
}

func TestGetLoggingContext(t *testing.T) {
	t.Run("Should retrieve the logging context when populated", func(t *testing.T) {
		expectedContext := LoggingContext{
			"hello": "world",
		}
		ctx := WithLoggingContext(context.Background(), expectedContext)
		sameContext := GetLoggingContext(ctx)
		require.Equal(t, expectedContext, sameContext)
	})
	t.Run("Should give back an empty context if not found", func(t *testing.T) {
		loggingContext := GetLoggingContext(context.Background())
		require.Empty(t, loggingContext)
	})
}
