package datadog

import (
	"io"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	"github.com/EurosportDigital/global-transcoding-platform/lib/monitoring"
)

type reporter struct {
	client   statsd.ClientInterface
	hostname string
}

// Reporter is an interface that defines the functionality provided by the Datadog client.
type Reporter interface {
	monitoring.MetricsReporter
	monitoring.EventReporter
	monitoring.HealthCheckReporter
	io.Closer
}

// These are stored as variables here so that they can easily be overridden in tests.
var since = time.Since
var logError = logger.Errorf

// InstallReporter initializes the default statsd client and uses it to install a new Datadog reporter.
func InstallReporter(agentAddress string, hostname string, namespace string) (Reporter, error) {
	client, err := statsd.New(agentAddress, statsd.WithNamespace(namespace))
	if err != nil {
		return nil, errors.Wrap(err, "creating datadog statsd client")
	}

	return InstallReporterUsingClient(client, hostname), nil
}

// InstallReporterUsingClient installs a new Datadog reporter, utilizing a specific statsd client, that provides metrics, event, and health check reporting.
func InstallReporterUsingClient(client statsd.ClientInterface, hostname string) Reporter {
	newReporter := &reporter{client: client, hostname: hostname}
	monitoring.Metrics = newReporter
	monitoring.Events = newReporter
	monitoring.HealthChecks = newReporter

	return newReporter
}

func (r *reporter) Count(name string, value int64, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Count(name, value, formatTags(tags), 1) })
}

func (r *reporter) Increment(name string, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Incr(name, formatTags(tags), 1) })
}

func (r *reporter) Decrement(name string, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Decr(name, formatTags(tags), 1) })
}

func (r *reporter) Gauge(name string, value float64, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Gauge(name, value, formatTags(tags), 1) })
}

func (r *reporter) Histogram(name string, value float64, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Histogram(name, value, formatTags(tags), 1) })
}

func (r *reporter) Distribution(name string, value float64, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Distribution(name, value, formatTags(tags), 1) })
}

func (r *reporter) Set(name string, value string, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Set(name, value, formatTags(tags), 1) })
}

func (r *reporter) Time(name string, value time.Duration, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Timing(name, value, formatTags(tags), 1) })
}

func (r *reporter) TimeSince(name string, startTime time.Time, tags ...monitoring.Tag) {
	r.withLogging(func() error { return r.client.Timing(name, since(startTime), formatTags(tags), 1) })
}

func (r *reporter) Emit(e *monitoring.Event) {
	r.withLogging(func() error {
		return r.client.Event(
			&statsd.Event{
				AggregationKey: e.AggregationKey,
				AlertType:      toEventAlertType(e.Type),
				Hostname:       r.hostname,
				Priority:       toEventPriority(e.Priority),
				SourceTypeName: e.Source,
				Tags:           formatTags(e.Tags),
				Text:           e.Description,
				Timestamp:      e.Timestamp,
				Title:          e.Title,
			})
	})
}

func (r *reporter) EmitSimple(title, description string) {
	r.withLogging(func() error { return r.client.SimpleEvent(title, description) })
}

func (r *reporter) Report(hc *monitoring.HealthCheck) {
	r.withLogging(func() error {
		return r.client.ServiceCheck(
			&statsd.ServiceCheck{
				Hostname:  r.hostname,
				Message:   hc.Message,
				Name:      hc.Name,
				Status:    toDatadogStatus(hc.Status),
				Tags:      formatTags(hc.Tags),
				Timestamp: hc.Timestamp,
			})
	})
}

func (r *reporter) ReportSimple(name string, status monitoring.HealthCheckStatus) {
	r.withLogging(func() error { return r.client.SimpleServiceCheck(name, toDatadogStatus(status)) })
}

func (r *reporter) Close() error {
	return r.client.Close()
}

func (r *reporter) withLogging(f func() error) {
	if err := f(); err != nil {
		logError("Error occurred when reporting to Datadog statsd client: %+v", err)
	}
}

func toEventAlertType(level monitoring.EventType) statsd.EventAlertType {
	switch level {
	case monitoring.Success:
		return statsd.Success
	case monitoring.Info:
		return statsd.Info
	case monitoring.Warning:
		return statsd.Warning
	case monitoring.Error:
		return statsd.Error
	default:
		return ""
	}
}

func toEventPriority(priority monitoring.EventPriority) statsd.EventPriority {
	switch priority {
	case monitoring.Low:
		return statsd.Low
	case monitoring.Normal:
		return statsd.Normal
	default:
		return ""
	}
}

func toDatadogStatus(status monitoring.HealthCheckStatus) statsd.ServiceCheckStatus {
	switch status {
	case monitoring.Ok:
		return statsd.Ok
	case monitoring.Warn:
		return statsd.Warn
	case monitoring.Critical:
		return statsd.Critical
	default:
		return statsd.Unknown
	}
}

func formatTags(tags []monitoring.Tag) []string {
	result := make([]string, len(tags))

	for i := range tags {
		result[i] = tags[i].Key + ":" + tags[i].Value
	}

	return result
}
