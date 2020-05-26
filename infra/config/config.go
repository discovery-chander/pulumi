package config

import "github.com/pulumi/pulumi/sdk/go/pulumi"

// ToPulumiStringArray converts a simple string array to Pulumi-readable string array
func ToPulumiStringArray(strings []string) pulumi.StringArrayInput {
	var res []pulumi.StringInput
	for _, s := range strings {
		res = append(res, pulumi.String(s))
	}
	return pulumi.StringArray(res)
}
