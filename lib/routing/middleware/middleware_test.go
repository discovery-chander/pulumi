package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/monitoring"
	"github.com/EurosportDigital/global-transcoding-platform/lib/monitoring/mocks"
	"github.com/EurosportDigital/global-transcoding-platform/lib/routing"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestMetricMiddleware(t *testing.T) {
	t.Run("Should capture metrics from handler", func(t *testing.T) {
		var (
			expectedPathTemplate = "/api/model/{id}"
			expectedPath         = "/api/model/1"
			expectedMethod       = http.MethodGet
			expectedStatusCode   = http.StatusTeapot
		)
		testRouter := mux.NewRouter()
		testRouter.HandleFunc(expectedPathTemplate, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(expectedStatusCode)
		}).Methods(expectedMethod)
		testWriter := httptest.NewRecorder()
		testRequest := httptest.NewRequest(expectedMethod, expectedPath, nil)
		mockReporter := &mocks.MetricsReporter{}
		defer mockReporter.AssertExpectations(t)
		mockReporter.On("Histogram", "endpoint", float64(1),
			monitoring.Tag{
				Key:   "route",
				Value: expectedPathTemplate,
			},
			monitoring.Tag{
				Key:   "method",
				Value: expectedMethod,
			},
			monitoring.Tag{
				Key:   "status_code",
				Value: strconv.Itoa(expectedStatusCode),
			})

		testRouter.Use(MetricMiddleware(mockReporter))
		testRouter.ServeHTTP(testWriter, testRequest)
	})
	t.Run("Should use hardcoded value when there is a failure in retrieving path template", func(t *testing.T) {
		testWriter := httptest.NewRecorder()
		testRequest := httptest.NewRequest(http.MethodGet, "/api/job/1", nil)
		mockReporter := &mocks.MetricsReporter{}
		defer mockReporter.AssertExpectations(t)
		mockReporter.On("Histogram", "endpoint", float64(1),
			monitoring.Tag{
				Key:   "route",
				Value: failedToGetPathTemplateMessage,
			},
			monitoring.Tag{
				Key:   "method",
				Value: http.MethodGet,
			},
			monitoring.Tag{
				Key:   "status_code",
				Value: "200",
			})
		MetricMiddleware(mockReporter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(testWriter, testRequest)
	})
}

func TestLoggingContextMiddleware(t *testing.T) {
	t.Run("Should capture metrics from handler", func(t *testing.T) {
		var (
			expectedPathTemplate = "/api/model/{id}"
			expectedPath         = "/api/model/1"
			expectedMethod       = http.MethodGet
			expectedContext      = routing.LoggingContext{
				"path":         expectedPath,
				"pathTemplate": expectedPathTemplate,
				"method":       expectedMethod,
			}
		)
		testRouter := mux.NewRouter()
		testRouter.HandleFunc(expectedPathTemplate, func(w http.ResponseWriter, r *http.Request) {
			loggingContext := routing.GetLoggingContext(r.Context())
			require.EqualValues(t, expectedContext, loggingContext)
		}).Methods(expectedMethod)
		testWriter := httptest.NewRecorder()
		testRequest := httptest.NewRequest(expectedMethod, expectedPath, nil)

		testRouter.Use(LoggingContextMiddleware())
		testRouter.ServeHTTP(testWriter, testRequest)
	})
	t.Run("Should use hardcoded value when there is a failure in retrieving path template", func(t *testing.T) {
		testWriter := httptest.NewRecorder()
		testRequest := httptest.NewRequest(http.MethodGet, "/api/job/1", nil)
		LoggingContextMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggingContext := routing.GetLoggingContext(r.Context())
			require.EqualValues(t, failedToGetPathTemplateMessage, loggingContext["pathTemplate"])
		})).ServeHTTP(testWriter, testRequest)
	})
}
