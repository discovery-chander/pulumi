package utils

import (
	"encoding/json"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/alecthomas/jsonschema"
	"github.com/xeipuuv/gojsonschema"
)

// GenerateSchema generates the schema of a json given a Go type. Using alecthomas/jsonschema it's possible to create
// annotations that add various parameters to the schema as well as extra parameters
func GenerateSchema(goType interface{}) ([]byte, error) {
	jsonSchema, err := json.Marshal(jsonschema.Reflect(goType))
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate schema")
	}
	return jsonSchema, nil
}

// ValidateRequest validates a request given the go type to use as a schema reference to and the request body
// to be validated
func ValidateRequest(goType interface{}, requestBody []byte) ([]gojsonschema.ResultError, error) {
	jobJsonSchema, err := GenerateSchema(goType)
	if err != nil {
		return nil, errors.Wrap(err, "unable to validate schema")
	}
	schemaLoader := gojsonschema.NewBytesLoader(jobJsonSchema)
	requestLoader := gojsonschema.NewStringLoader(string(requestBody))
	result, err := gojsonschema.Validate(schemaLoader, requestLoader)
	if err != nil {
		return nil, errors.Wrap(err, "unable to validate schema")
	}
	if !result.Valid() {
		return result.Errors(), nil
	}
	return nil, nil
}
