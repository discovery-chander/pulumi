package monitoring

// A Tag is a key value pair that can be used to add dimensions to metrics, events, and health checks.
type Tag struct {
	Key   string
	Value string
}
