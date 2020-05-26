package utils

import (
	"encoding/json"
	"testing"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/stretchr/testify/require"
)

type TestCaseStruct struct {
	TestInt    int    `json:"testInt" jsonschema:"minimum=1,maximum=100"`
	TestString string `json:"testString" jsonschema:"pattern=^(testPattern:)+"`
}
type TestFailCaseStruct struct {
	TestInt int `json:"testInt" jsonschema_extras:"minimum=1"`
}

func TestCreateJsonSchema(t *testing.T) {
	t.Run("Should return a json schema when given a compatible struct", func(t *testing.T) {
		schema, err := GenerateSchema(&TestCaseStruct{})
		require.NoError(t, err)
		expectedSchemaString := "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"$ref\":\"#/definitions/TestCaseStruct\",\"definitions\":{\"TestCaseStruct\":{\"required\":[\"testInt\",\"testString\"],\"properties\":{\"testInt\":{\"maximum\":100,\"minimum\":1,\"type\":\"integer\"},\"testString\":{\"pattern\":\"^(testPattern:)+\",\"type\":\"string\"}},\"additionalProperties\":false,\"type\":\"object\"}}}"
		require.Equal(t, string(schema), expectedSchemaString, "Test schema should return expected result")
	})
}

func TestValidateJsonSchema(t *testing.T) {
	t.Run("Should return successful validation", func(t *testing.T) {
		payload, marshalErr := json.Marshal(TestCaseStruct{
			TestInt:    2,
			TestString: "testPattern:mock",
		})
		require.NoError(t, marshalErr)
		resultErrors, err := ValidateRequest(&TestCaseStruct{}, payload)
		require.NoError(t, err)
		require.Nil(t, resultErrors, "Validation should not return validation errors")
	})
	t.Run("Should return invalid range validation error", func(t *testing.T) {
		payload, marshalErr := json.Marshal(TestCaseStruct{
			TestInt:    -1,
			TestString: "testPattern:mock",
		})
		require.NoError(t, marshalErr)
		resultErrors, err := ValidateRequest(&TestCaseStruct{}, payload)
		expectedError := "testInt: Must be greater than or equal to 1"
		require.Nil(t, err, "Validating schema should not return an error")
		require.Equal(t, resultErrors[0].String(), expectedError, "Validation should return expected validation error")
	})
	t.Run("Should return invalid pattern validation error", func(t *testing.T) {
		payload, err := json.Marshal(TestCaseStruct{
			TestInt:    2,
			TestString: "invalidPattern",
		})
		require.NoError(t, err)
		resultErrors, err := ValidateRequest(&TestCaseStruct{}, payload)
		expectedError := "testString: Does not match pattern '^(testPattern:)+'"
		require.Nil(t, err, "Validating schema should not return an error")
		require.Equal(t, resultErrors[0].String(), expectedError, "Validation should return expected validation error")
	})
	t.Run("Should return required attribute missing validation error", func(t *testing.T) {
		payload, marshalErr := json.Marshal(TestFailCaseStruct{
			TestInt: 2,
		})
		require.NoError(t, marshalErr)
		resultErrors, err := ValidateRequest(&TestCaseStruct{}, payload)
		expectedError := "(root): testString is required"
		require.Nil(t, err, "Validating schema should not return an error")
		require.Equal(t, resultErrors[0].String(), expectedError, "Validation should return expected validation error")
	})
	t.Run("Should return internal failure due to malformed gojson attribute", func(t *testing.T) {
		payload, marshalErr := json.Marshal(TestCaseStruct{
			TestInt: 2,
		})
		require.NoError(t, marshalErr)
		resultErrors, err := ValidateRequest(&TestFailCaseStruct{}, payload)
		expectedInternalError := "minimum must be of a number"
		require.Nil(t, resultErrors, "Validating against a malformed schema should not return a result")
		require.EqualError(t, errors.Cause(err), expectedInternalError, "Validating schema should return expected error")
	})
}
