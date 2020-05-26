package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// ZeroLogger is an interface for the Error, Info and Warn zerolog functions.
type ZeroLogger interface {
	Error() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Debug() *zerolog.Event
	With() zerolog.Context
}

var loggerInstance ZeroLogger

// SetLogger creates a logger instance with the given ZeroLogger config.
func SetLogger(baseLogger ZeroLogger) {
	loggerInstance = baseLogger
}

// logger() is a private function that uses existing loggerInstance if available and creates a new one if not.
func logger() ZeroLogger {
	if loggerInstance != nil {
		return loggerInstance
	}

	zerolog.TimeFieldFormat = time.RFC3339
	var baseLogger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	SetLogger(&baseLogger)
	return loggerInstance
}

var (
	invalidArgMessage      = "Invalid arg: %s"
	invalidArgValueMessage = "Invalid value for argument: %s: %v"
	missingArgMessage      = "Missing arg: %s"
)

// InvalidArg is a standard error message.
func InvalidArg(argumentName string) {
	logger().Error().Msgf(invalidArgMessage, argumentName)
}

// InvalidArgValue is a standard error message.
func InvalidArgValue(argument string, argumentValue string) {
	logger().Error().Msgf(invalidArgValueMessage, argument, argumentValue)
}

// MissingArg is a standard error message.
func MissingArg(argument string) {
	logger().Error().Msgf(missingArgMessage, argument)
}

func Error(err error, message string) {
	logger().Error().Msgf("%v: %+v", message, err)
}

// Error is a standard error message.
func Errorf(msgFormat string, v ...interface{}) {
	logger().Error().Msgf(msgFormat, v...)
}

func ErrorWithFields(err error, message string, fields map[string]interface{}) {
	log := logger().With().Fields(fields).Logger()
	log.Error().Msgf("%v: %+v", message, err)
}

// Info is a standard info message.
func Infof(msgFormat string, v ...interface{}) {
	logger().Info().Msgf(msgFormat, v...)
}

func InfoWithFields(message string, fields map[string]interface{}) {
	log := logger().With().Fields(fields).Logger()
	log.Info().Msg(message)
}

// Warn is a standard warn message.
func Warnf(msgFormat string, v ...interface{}) {
	logger().Warn().Msgf(msgFormat, v...)
}

func WarnWithFields(message string, fields map[string]interface{}) {
	log := logger().With().Fields(fields).Logger()
	log.Warn().Msg(message)
}

func Debugf(msgFormat string, v ...interface{}) {
	logger().Debug().Msgf(msgFormat, v...)
}

func DebugWithFields(message string, fields map[string]interface{}) {
	log := logger().With().Fields(fields).Logger()
	log.Debug().Msg(message)
}
