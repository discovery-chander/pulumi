package logger

import (
	stderrors "errors"
	"os"
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/logger/mocks"
	"github.com/rs/zerolog"
)

const message = "test message"
const argumentName = "testName"
const argumentValue = "testValue"

func TestInfoLogger(t *testing.T) {
	mockZerolog := mocks.ZeroLogger{}
	var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	var event = logger.Info()
	mockZerolog.On("Info").Return(event)
	SetLogger(&mockZerolog)
	Infof(message)
	mockZerolog.AssertExpectations(t)
}

func TestErrorLogger(t *testing.T) {
	mockZerolog := mocks.ZeroLogger{}
	var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	var event = logger.Error()
	mockZerolog.On("Error").Return(event)
	SetLogger(&mockZerolog)
	Errorf(message)
	mockZerolog.AssertExpectations(t)
}

func TestWarnLogger(t *testing.T) {
	mockZerolog := mocks.ZeroLogger{}
	var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	var event = logger.Warn()
	mockZerolog.On("Warn").Return(event)
	SetLogger(&mockZerolog)
	Warnf(message)
	mockZerolog.AssertExpectations(t)
}

func TestInvalidArgLogger(t *testing.T) {
	mockZerolog := mocks.ZeroLogger{}
	var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	var event = logger.Error()
	mockZerolog.On("Error").Return(event)
	SetLogger(&mockZerolog)
	InvalidArg(argumentName)
	mockZerolog.AssertExpectations(t)
}

func TestInvalidArgValueLogger(t *testing.T) {
	mockZerolog := mocks.ZeroLogger{}
	var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	var event = logger.Error()
	mockZerolog.On("Error").Return(event)
	SetLogger(&mockZerolog)
	InvalidArgValue(argumentName, argumentValue)
	mockZerolog.AssertExpectations(t)
}

func TestMissingArgLogger(t *testing.T) {
	mockZerolog := mocks.ZeroLogger{}
	var logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	var event = logger.Error()
	mockZerolog.On("Error").Return(event)
	SetLogger(&mockZerolog)
	MissingArg(message)
	mockZerolog.AssertExpectations(t)
}

func TestLoggingWithFields(t *testing.T) {
	mockZerolog := &mocks.ZeroLogger{}
	SetLogger(mockZerolog)

	funcsToTest := []func(string, map[string]interface{}){
		InfoWithFields,
		WarnWithFields,
		DebugWithFields,
	}

	fields := map[string]interface{}{"foo": "bar"}
	for _, f := range funcsToTest {
		mockZerolog.On("With").Return(zerolog.Context{})
		f("message", fields)
	}

	mockZerolog.On("With").Return(zerolog.Context{})
	ErrorWithFields(stderrors.New("some error"), "message", fields)

	mockZerolog.AssertExpectations(t)
}
