package utils

import (
	"github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi/sdk/go/pulumi/config"
)

type ServiceEnv struct {
	Env           string
	DatadogAPIKey string
}

func GetServiceEnvironment(ctx *pulumi.Context) *ServiceEnv {
	cfg := config.New(ctx, "")
	env := cfg.Require("env")
	datadogAPIKey := config.Require(ctx, "datadog:apiKey")

	return &ServiceEnv{
		Env:           env,
		DatadogAPIKey: datadogAPIKey,
	}
}

func GetConnectionStringFromBaseInfra(mainInfra *pulumi.StackReference) pulumi.StringOutput {
	rdsEndpoint := mainInfra.GetOutput(pulumi.String("rdsEndpoint"))
	rdsPort := mainInfra.GetOutput(pulumi.String("rdsPort"))
	rdsUsername := mainInfra.GetOutput(pulumi.String("rdsUsername"))
	rdsDbName := mainInfra.GetOutput(pulumi.String("rdsDbName"))
	rdsPassword := mainInfra.GetOutput(pulumi.String("rdsPassword"))
	return pulumi.Sprintf("host=%s port=%.f user=%s dbname=%s password=%s", rdsEndpoint, rdsPort, rdsUsername, rdsDbName, rdsPassword)
}
