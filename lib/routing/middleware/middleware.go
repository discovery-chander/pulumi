package middleware

import (
	"net/http"
	"strconv"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/EurosportDigital/global-transcoding-platform/lib/logger"
	"github.com/EurosportDigital/global-transcoding-platform/lib/monitoring"
	"github.com/EurosportDigital/global-transcoding-platform/lib/routing"
	"github.com/gorilla/mux"
)

const failedToGetPathTemplateMessage = "failed_to_get_template"

// MuxMiddlewareFunc is an alias to the type mux.Use expects.
type MuxMiddlewareFunc = func(next http.Handler) http.Handler

// MetricMiddleware automatically collects a metric for the endpoint.
func MetricMiddleware(reporter monitoring.MetricsReporter) MuxMiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pathTemplate, err := getPathTemplate(r)
			if err != nil {
				logger.Error(err, "Could not get path template, route may not have been setup correctly")
				pathTemplate = failedToGetPathTemplateMessage
			}
			statusWrapper := &statusResponseWriter{w, 200}
			defer func() {
				reporter.Histogram("endpoint", 1, []monitoring.Tag{
					{
						Key:   "route",
						Value: pathTemplate,
					},
					{
						Key:   "method",
						Value: r.Method,
					},
					{
						Key:   "status_code",
						Value: strconv.Itoa(statusWrapper.statusCode),
					},
				}...)
			}()

			next.ServeHTTP(statusWrapper, r)
		})
	}
}

func LoggingContextMiddleware() MuxMiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pathTemplate, err := getPathTemplate(r)
			if err != nil {
				logger.Error(err, "Could not get path template, route may not have been setup correctly")
				pathTemplate = failedToGetPathTemplateMessage
			}

			loggingContext := routing.LoggingContext{
				"path":         r.URL.String(),
				"pathTemplate": pathTemplate,
				"method":       r.Method,
			}

			requestWithContext := r.WithContext(routing.WithLoggingContext(r.Context(), loggingContext))
			next.ServeHTTP(w, requestWithContext)
		})
	}
}

func getPathTemplate(r *http.Request) (string, error) {
	currentRoute := mux.CurrentRoute(r)
	if currentRoute == nil {
		return "", errors.New("could not get route")
	}

	pathTemplate, err := currentRoute.GetPathTemplate()
	if err != nil {
		return "", errors.Wrap(err, "could not get path template")
	}

	return pathTemplate, nil
}
