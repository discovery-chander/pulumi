package main

import (
	"log"

	"github.com/EurosportDigital/global-transcoding-platform/infra/modules"
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Setup a private VPC (with subnets) for our environment.
		vpcInfo, err := modules.SetupVpc(ctx)
		if err != nil {
			return errors.Wrap(err, "setting up VPC")
		}

		// Setup an RDS instance.
		err = modules.SetupRDS(ctx, vpcInfo)
		if err != nil {
			return errors.Wrap(err, "setting up RDS")
		}

		// Setup an ECS cluster.
		err = modules.SetupECSCluster(ctx)
		if err != nil {
			return errors.Wrap(err, "setting up cluster")
		}

		log.Print("reading gtp:user value from pulumi context")
		accountID, userExists := ctx.GetConfig("gtp:user")

		if userExists {
			err = modules.CreateDatadogIntegration(ctx, "dataDogIntegration", accountID)
			if err != nil {
				return errors.Wrap(err, "setting up datadog integration")
			}
		}

		_, err = modules.NewIngestionPipeline(ctx)
		if err != nil {
			return errors.Wrap(err, "setting up ingestion-pipeline")
		}

		err = modules.SetupS3InputOutput(ctx)
		if err != nil {
			return errors.Wrap(err, "setting up s3-pipeline")
		}

		_, err = modules.CreateJobSchedulingQueue(ctx)
		if err != nil {
			return errors.Wrap(err, "setting up job scheduling queue")
		}

		return nil
	})
}
