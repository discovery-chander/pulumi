package routing

import "context"

type contextKey int

// LoggingContext defines the fields that show up in logging output.
type LoggingContext = map[string]interface{}

const (
	keyLoggingContext contextKey = iota
)

// WithLoggingContext adds the logging context to the context.
func WithLoggingContext(ctx context.Context, loggingContext LoggingContext) context.Context {
	return context.WithValue(ctx, keyLoggingContext, loggingContext)
}

// GetLoggingContext retrieves the logging context from the context, otherwise returns an empty context.
func GetLoggingContext(ctx context.Context) LoggingContext {
	if loggingContext := ctx.Value(keyLoggingContext); loggingContext != nil {
		return loggingContext.(LoggingContext)
	}
	return make(LoggingContext)
}
