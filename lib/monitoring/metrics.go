package monitoring

import (
	"time"
)

// Metrics can be used to report metrics to the underlying metrics collector.
var Metrics MetricsReporter

// MetricsReporter is an interface that defines the methods expected to be implemented by metrics reporting clients.
type MetricsReporter interface {
	// Count tracks how many times something happened per reporting interval.
	Count(name string, value int64, tags ...Tag)

	// Increment is simply Count with a value of 1.
	Increment(name string, tags ...Tag)

	// Decrement is simply Count with a value of -1.
	Decrement(name string, tags ...Tag)

	// Gauge measures the value of a metric at the time of reporting.
	Gauge(name string, value float64, tags ...Tag)

	// Histogram tracks the statistical distribution of a set of values submitted to this reporter.
	Histogram(name string, value float64, tags ...Tag)

	// Distribution tracks the statistical distribution of a set of values submitted to all reporters that report to the same monitoring system.
	Distribution(name string, value float64, tags ...Tag)

	// Set counts the number of unique elements in a group.
	Set(name string, value string, tags ...Tag)

	// Time reports how long an operation took to complete.
	Time(name string, value time.Duration, tags ...Tag)

	// TimeSince reports the amount of time from the start time to now.
	TimeSince(name string, startTime time.Time, tags ...Tag)
}
