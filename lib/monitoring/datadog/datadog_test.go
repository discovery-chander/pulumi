package datadog

import (
	"errors"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	"github.com/EurosportDigital/global-transcoding-platform/lib/monitoring"
	"github.com/EurosportDigital/global-transcoding-platform/lib/monitoring/datadog/mocks"
	"github.com/stretchr/testify/assert"
)

type reporterTestCase struct {
	name    string
	execute func(*mocks.ClientInterface, Reporter)
}

const hostname = "gtp.discovery.internal"

var errExpected = errors.New("error")

// Tests that the methods exposed through the monitoring.MetricsReporter interface correctly delegate
// to Datadog's statsd client and log any errors it returns
func TestMetricsReporting(t *testing.T) {
	metric, tags := "TestMetric", []monitoring.Tag{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}}
	expectedTags := []string{"key1:value1", "key2:value2"}
	// Each test case defined below creates the assertion of the call to the mock statsd client and then
	// invokes the appropriate method on the monitoring.MetricsReporter instance
	tests := []reporterTestCase{
		{"Count", func(m *mocks.ClientInterface, r Reporter) {
			value := int64(123)
			m.On("Count", metric, value, expectedTags, 1.0).Return(errExpected)
			r.Count(metric, value, tags...)
		}},
		{"Increment", func(m *mocks.ClientInterface, r Reporter) {
			m.On("Incr", metric, expectedTags, 1.0).Return(errExpected)
			r.Increment(metric, tags...)
		}},
		{"Decrement", func(m *mocks.ClientInterface, r Reporter) {
			m.On("Decr", metric, expectedTags, 1.0).Return(errExpected)
			r.Decrement(metric, tags...)
		}},
		{"Gauge", func(m *mocks.ClientInterface, r Reporter) {
			value := 1.67
			m.On("Gauge", metric, value, expectedTags, 1.0).Return(errExpected)
			r.Gauge(metric, value, tags...)
		}},
		{"Histogram", func(m *mocks.ClientInterface, r Reporter) {
			value := 1.67
			m.On("Histogram", metric, value, expectedTags, 1.0).Return(errExpected)
			r.Histogram(metric, value, tags...)
		}},
		{"Distribution", func(m *mocks.ClientInterface, r Reporter) {
			value := 1.67
			m.On("Distribution", metric, value, expectedTags, 1.0).Return(errExpected)
			r.Distribution(metric, value, tags...)
		}},
		{"Set", func(m *mocks.ClientInterface, r Reporter) {
			value := "set_value"
			m.On("Set", metric, value, expectedTags, 1.0).Return(errExpected)
			r.Set(metric, value, tags...)
		}},
		{"Time", func(m *mocks.ClientInterface, r Reporter) {
			value, _ := time.ParseDuration("345ms")
			m.On("Timing", metric, value, expectedTags, 1.0).Return(errExpected)
			r.Time(metric, value, tags...)
		}},
		{"TimeSince", func(m *mocks.ClientInterface, r Reporter) {
			now := time.Now()
			value := now.Add(-time.Minute)
			since = now.Sub
			defer func() { since = time.Since }()
			m.On("Timing", metric, time.Minute, expectedTags, 1.0).Return(errExpected)
			r.TimeSince(metric, value, tags...)
		}},
	}

	runReporterTests(t, tests)
}

// Tests that the methods exposed through the monitoring.EventReporter interface correctly delegate
// to Datadog's statsd client and log any errors it returns
func TestEventReporting(t *testing.T) {
	event := &monitoring.Event{
		AggregationKey: "agg_key",
		Description:    "description",
		Tags:           []monitoring.Tag{{Key: "key1", Value: "value1"}, {Key: "key2", Value: "value2"}, {Key: "key3", Value: "value3"}},
		Priority:       monitoring.Normal,
		Source:         "tests",
		Timestamp:      time.Now(),
		Title:          "TestEvent",
		Type:           monitoring.Error,
	}
	statsdEvent := &statsd.Event{
		AggregationKey: event.AggregationKey,
		AlertType:      statsd.Error,
		Hostname:       hostname,
		Priority:       statsd.Normal,
		SourceTypeName: event.Source,
		Tags:           []string{"key1:value1", "key2:value2", "key3:value3"},
		Text:           event.Description,
		Timestamp:      event.Timestamp,
		Title:          event.Title,
	}

	// Each test case defined below creates the assertion of the call to the mock statsd client and then
	// invokes the appropriate method on the monitoring.EventReporter instance
	tests := []reporterTestCase{
		{"Emit", func(m *mocks.ClientInterface, r Reporter) {
			m.On("Event", statsdEvent).Return(errExpected)
			r.Emit(event)
		}},
		{"EmitSimple", func(m *mocks.ClientInterface, r Reporter) {
			m.On("SimpleEvent", statsdEvent.Title, statsdEvent.Text).Return(errExpected)
			r.EmitSimple(event.Title, event.Description)
		}},
	}

	runReporterTests(t, tests)
}

// Tests that the methods exposed through the monitoring.HealthCheckReporter interface correctly delegate
// to Datadog's statsd client and log any errors it returns
func TestHealthCheckReporting(t *testing.T) {
	healthCheck := &monitoring.HealthCheck{
		Tags:      []monitoring.Tag{{Key: "key1", Value: "value1"}},
		Message:   "test message",
		Name:      "TestCheck",
		Status:    monitoring.Critical,
		Timestamp: time.Now(),
	}
	statsdServiceCheck := &statsd.ServiceCheck{
		Hostname:  hostname,
		Message:   healthCheck.Message,
		Name:      healthCheck.Name,
		Status:    statsd.Critical,
		Tags:      []string{"key1:value1"},
		Timestamp: healthCheck.Timestamp,
	}

	// Each test case defined below creates the assertion of the call to the mock statsd client and then
	// invokes the appropriate method on the monitoring.HealthCheckReporter instance
	tests := []reporterTestCase{
		{"Report", func(m *mocks.ClientInterface, r Reporter) {
			m.On("ServiceCheck", statsdServiceCheck).Return(errExpected)
			r.Report(healthCheck)
		}},
		{"ReportSimple", func(m *mocks.ClientInterface, r Reporter) {
			m.On("SimpleServiceCheck", statsdServiceCheck.Name, statsdServiceCheck.Status).Return(errExpected)
			r.ReportSimple(healthCheck.Name, healthCheck.Status)
		}},
	}

	runReporterTests(t, tests)
}

// Tests that Close correctly closes Datadog's statsd client and returns any error
func TestClose(t *testing.T) {
	mockClient := &mocks.ClientInterface{}
	reporter := InstallReporterUsingClient(mockClient, hostname)

	mockClient.On("Close").Return(errExpected)

	assert.Equal(t, errExpected, reporter.Close())
	mockClient.AssertExpectations(t)
}

func TestEventPriorityConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    monitoring.EventPriority
		expected statsd.EventPriority
	}{
		{name: "Low", input: monitoring.Low, expected: statsd.Low},
		{name: "Normal", input: monitoring.Normal, expected: statsd.Normal},
		{name: "Unknown", input: -1, expected: ""},
	}

	for _, test := range tests {
		input, expected := test.input, test.expected
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, expected, toEventPriority(input))
		})
	}
}

func TestEventTypeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    monitoring.EventType
		expected statsd.EventAlertType
	}{
		{name: "Info", input: monitoring.Info, expected: statsd.Info},
		{name: "Warning", input: monitoring.Warning, expected: statsd.Warning},
		{name: "Error", input: monitoring.Error, expected: statsd.Error},
		{name: "Success", input: monitoring.Success, expected: statsd.Success},
		{name: "Unknown", input: -1, expected: ""},
	}

	for _, test := range tests {
		input, expected := test.input, test.expected
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, expected, toEventAlertType(input))
		})
	}
}

func TestHealthCheckStatusConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    monitoring.HealthCheckStatus
		expected statsd.ServiceCheckStatus
	}{
		{name: "Ok", input: monitoring.Ok, expected: statsd.Ok},
		{name: "Warn", input: monitoring.Warn, expected: statsd.Warn},
		{name: "Critical", input: monitoring.Critical, expected: statsd.Critical},
		{name: "Unknown", input: monitoring.Unknown, expected: statsd.Unknown},
		{name: "UnknownInput", input: -1, expected: statsd.Unknown},
	}

	for _, test := range tests {
		input, expected := test.input, test.expected
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, expected, toDatadogStatus(input))
		})
	}
}

func TestInstallReporterReturnsReporter(t *testing.T) {
	client, err := InstallReporter("127.0.0.1:8125", hostname, "test")

	assert.NotNil(t, client)
	assert.Nil(t, err)
}

func TestInstallReporterReturnsError(t *testing.T) {
	client, err := InstallReporter("", hostname, "test")

	assert.Nil(t, client)
	assert.NotNil(t, err)
}

// Adds an error logging assertion to the test cases passed to it, runs the test cases, and verifies any mock client assertions
func runReporterTests(t *testing.T, tests []reporterTestCase) {
	for _, test := range tests {
		executeTest := test.execute
		t.Run(test.name+"CallsDatadogAndLogsAnyErrors", func(t *testing.T) {
			mockClient := &mocks.ClientInterface{}
			reporter := InstallReporterUsingClient(mockClient, hostname)
			logError = verifyErrorLog(t, "Error occurred when reporting to Datadog statsd client: %+v", errExpected)
			defer func() { logError = logger.Errorf }()

			executeTest(mockClient, reporter)

			mockClient.AssertExpectations(t)
		})
	}
}

func verifyErrorLog(t *testing.T, msg string, args ...interface{}) func(string, ...interface{}) {
	return func(actualMsg string, actualArgs ...interface{}) {
		assert.Equal(t, msg, actualMsg)
		assert.Equal(t, args, actualArgs)
	}
}
