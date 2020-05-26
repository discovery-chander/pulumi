package monitoring

import (
	"time"
)

// HealthChecks can be used to report health checks to the underlying health check collector.
var HealthChecks HealthCheckReporter

// HealthCheckReporter is an interface that defines the methods expected to be implemented by health check reporting clients.
type HealthCheckReporter interface {
	// Report sends the provided HealthCheck.
	Report(hc *HealthCheck)

	// ReportSimple sends a health check with the provided name and status.
	ReportSimple(name string, status HealthCheckStatus)
}

// HealthCheck encapsulates the data representing a single health check.
type HealthCheck struct {
	// Name of the health check.
	Name string
	// Status of the system being checked.
	Status HealthCheckStatus
	// Timestamp is the time the check was performed.
	Timestamp time.Time
	// Message describes the state of the system being checked.
	Message string
	// Tags for the health check.
	Tags []Tag
}

// HealthCheckStatus indicates the status of the system being checked.
type HealthCheckStatus int

const (
	// Ok indicates the system being checked is operating normally.
	Ok HealthCheckStatus = iota
	// Warn indicates the system being checked is experiencing non-critical problems.
	Warn
	// Critical indicates the system being checked is not functioning correctly.
	Critical
	// Unknown indicates the status of the system being checked could not be determined.
	Unknown
)
