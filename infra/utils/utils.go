package utils

import (
	"encoding/json"
	"fmt"

	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

type EnvironmentBlockArgs struct {
	Blocks []EnvironmentArgs
}

type EnvironmentArgs struct {
	Name  string
	Value string
}

// Creates the name of a specific resource given the name of the stack and the resource
func CreateResourceName(ctx *pulumi.Context, resName string) string {
	//Read the stack name from the context
	stackName := ctx.Stack()
	//Crop the stack name and resource to accomodate to AWS Resource name limit
	if len(stackName) > 10 {
		stackName = stackName[:10]
	}
	if len(resName) > 9 {
		resName = resName[:9]
	}
	return fmt.Sprintf("%s-%s", stackName, resName)
}

// CreateEnvironmentBlockString used to create an environment variable block for EC2 task definitions using a provided map.
// Values will be mapped using the following format:
//
// {
//		"name": key,
//		"value": value,
// }
// If environment is nil an empty string will be returned.
func CreateEnvironmentBlockString(environment *EnvironmentBlockArgs) (string, error) {
	if environment == nil {
		return "[]", errors.New("Missing environment args")
	}

	// Marshal the map into a JSON string.
	empData, err := json.Marshal(environment.Blocks)
	if err != nil {
		return "[]", errors.Wrapf(err, "marshaling %v", environment.Blocks)
	}

	return string(empData), nil
}

// ToStringArrayOutput converts a generic pulumi.AnyOutput to a pulumi.StringArrayOutput.
func ToStringArrayOutput(output pulumi.AnyOutput) pulumi.StringArrayOutput {
	return output.ApplyStringArray(func(inputGeneric interface{}) []string {
		input := inputGeneric.([]interface{})
		result := make([]string, len(input))

		for i, v := range input {
			result[i] = v.(string)
		}

		return result
	})
}
