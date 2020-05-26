package resources

import (
	"fmt"
	"strings"

	"github.com/EurosportDigital/global-transcoding-platform/infra/config"
	"github.com/EurosportDigital/global-transcoding-platform/infra/utils"
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/pulumi/pulumi-aws/sdk/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/ecs"
	elb "github.com/pulumi/pulumi-aws/sdk/go/aws/elasticloadbalancingv2"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/sns"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/sqs"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
	pulumiConfig "github.com/pulumi/pulumi/sdk/go/pulumi/config"
)

// CreateLoadBalancer creates an AWS load balancer with given subnets and security group
func CreateLoadBalancer(ctx *pulumi.Context, subnetIds pulumi.StringArrayInput, secGroupIds pulumi.StringArrayInput, lbName string) (*elb.LoadBalancer, error) {
	return elb.NewLoadBalancer(ctx, utils.CreateResourceName(ctx, lbName), &elb.LoadBalancerArgs{
		Subnets:        subnetIds,
		SecurityGroups: secGroupIds,
	})
}

// CreateWebListener creates an Elastic Load Balancer web listener for given load balancer and target group
func CreateWebListener(ctx *pulumi.Context, webLb *elb.LoadBalancer, targetGroup *elb.TargetGroup, listenerName string) (*elb.Listener, error) {
	return elb.NewListener(ctx, utils.CreateResourceName(ctx, listenerName), &elb.ListenerArgs{
		LoadBalancerArn: webLb.Arn,
		Port:            pulumi.Int(80),
		DefaultActions: elb.ListenerDefaultActionArray{
			elb.ListenerDefaultActionArgs{
				Type:           pulumi.String("forward"),
				TargetGroupArn: targetGroup.Arn,
			},
		},
	})
}

// CreateVpc creates a private VPC in our AWS account.
func CreateVpc(ctx *pulumi.Context, name string, cidrBlock string) (*ec2.Vpc, error) {
	vpcName := utils.CreateResourceName(ctx, name)
	return ec2.NewVpc(ctx, vpcName, &ec2.VpcArgs{
		CidrBlock:                   pulumi.String(cidrBlock),
		EnableClassiclink:           pulumi.BoolPtr(false),
		EnableClassiclinkDnsSupport: pulumi.BoolPtr(true),
		EnableDnsHostnames:          pulumi.BoolPtr(true),
		EnableDnsSupport:            pulumi.BoolPtr(true),
		Tags:                        pulumi.Map{"Name": pulumi.StringPtr(vpcName)}})
}

// GetVpc reads back the given AWS VPC.
func GetVpc(ctx *pulumi.Context) (*ec2.LookupVpcResult, error) {
	typedBool := true
	return ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{
		Default: &typedBool,
	})
}

// CreateSubnet creates subnet in a given VPC in AWS.
func CreateSubnet(ctx *pulumi.Context, vpcId pulumi.StringInput, subnetName string, cidrBlock string, az string, public bool) (*ec2.Subnet, error) {
	return ec2.NewSubnet(ctx, utils.CreateResourceName(ctx, subnetName), &ec2.SubnetArgs{
		CidrBlock:           pulumi.String(cidrBlock),
		VpcId:               vpcId,
		AvailabilityZone:    pulumi.String(az),
		MapPublicIpOnLaunch: pulumi.BoolPtr(public),
	})
}

// CreateDBSubnetGroup creates a DB subnet group that can be used by RDS in a given VPC in AWS.
func CreateDbSubnetGroup(ctx *pulumi.Context, subnetIds pulumi.StringArrayInput, dbSubnetName string) (*rds.SubnetGroup, error) {
	return rds.NewSubnetGroup(ctx, utils.CreateResourceName(ctx, dbSubnetName), &rds.SubnetGroupArgs{
		SubnetIds: subnetIds,
	})
}

// CreateRouteTable creates a new routing table associated with the given VPC.
func CreateRouteTable(ctx *pulumi.Context, vpcId pulumi.StringInput, inetGwId pulumi.StringPtrInput, tableName string, cidrBlock string) (*ec2.RouteTable, error) {
	return ec2.NewRouteTable(ctx, utils.CreateResourceName(ctx, tableName), &ec2.RouteTableArgs{
		VpcId: vpcId,
		Routes: ec2.RouteTableRouteArray{
			ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String(cidrBlock),
				GatewayId: inetGwId,
			},
		},
	})
}

// CreateInternetGateway creates internet gateway that can be attached to a VPC.
func CreateInternetGateway(ctx *pulumi.Context, vpcId pulumi.StringPtrInput, gatewayName string) (*ec2.InternetGateway, error) {
	return ec2.NewInternetGateway(ctx, utils.CreateResourceName(ctx, gatewayName), &ec2.InternetGatewayArgs{
		VpcId: vpcId,
	})
}

// GetSubnets reads back the default AWS Subnets
func GetSubnets(ctx *pulumi.Context, vpcId string) (*ec2.GetSubnetIdsResult, error) {
	return ec2.GetSubnetIds(ctx, &ec2.GetSubnetIdsArgs{
		VpcId: vpcId,
	})
}

// CreateSecurityGroup creates a security group with provided CidrBlocks
func CreateSecurityGroup(ctx *pulumi.Context, vpcId pulumi.StringPtrInput, securityGroupName string, egressCidrs []string, ingressCidrs []string, port int, self bool) (*ec2.SecurityGroup, error) {
	if !self {
		return ec2.NewSecurityGroup(ctx, utils.CreateResourceName(ctx, securityGroupName), &ec2.SecurityGroupArgs{
			VpcId: vpcId,
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: config.ToPulumiStringArray(egressCidrs),
				},
			},
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(port),
					ToPort:     pulumi.Int(port),
					CidrBlocks: config.ToPulumiStringArray(ingressCidrs),
				},
			},
		})
	} else {
		return ec2.NewSecurityGroup(ctx, utils.CreateResourceName(ctx, securityGroupName), &ec2.SecurityGroupArgs{
			VpcId: vpcId,
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: config.ToPulumiStringArray(egressCidrs),
				},
			},
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol: pulumi.String("-1"),
					FromPort: pulumi.Int(0),
					ToPort:   pulumi.Int(0),
					Self:     pulumi.BoolPtr(self),
				},
			},
		})
	}
}

// CreateCluster creates a Fargate Cluster in AWS with provided name
func CreateCluster(ctx *pulumi.Context, clusterName string, clusterArgs *ecs.ClusterArgs) (*ecs.Cluster, error) {
	return ecs.NewCluster(ctx, utils.CreateResourceName(ctx, clusterName), clusterArgs)
}

// CreateIamRole creates an IAM role with provided AssumeRolePolicy JSON
func CreateIamRole(ctx *pulumi.Context, roleName string, policy pulumi.StringInput) (*iam.Role, error) {
	return iam.NewRole(ctx, utils.CreateResourceName(ctx, roleName), &iam.RoleArgs{
		AssumeRolePolicy: policy,
	})
}

// CreateIamPolicy creates an IAM policy with provided policy JSON
func CreateIamPolicy(ctx *pulumi.Context, policyName string, policy pulumi.StringInput) (*iam.Policy, error) {
	return iam.NewPolicy(ctx, utils.CreateResourceName(ctx, policyName), &iam.PolicyArgs{
		Policy: policy,
	})
}

// AttachPolicyToRole attaches given policy to given IAM role
func AttachPolicyToRole(ctx *pulumi.Context, role *iam.Role, policyArn pulumi.StringInput, policyName string) (*iam.RolePolicyAttachment, error) {
	return iam.NewRolePolicyAttachment(ctx, utils.CreateResourceName(ctx, policyName), &iam.RolePolicyAttachmentArgs{
		Role:      role.Name,
		PolicyArn: policyArn,
	})
}

// CreateTargetGroup creates an Elastic Load Balancer Target Group with given VPC
func CreateTargetGroup(ctx *pulumi.Context, vpcId pulumi.StringPtrInput, targetGroupName string, hcPath string) (*elb.TargetGroup, error) {
	return elb.NewTargetGroup(ctx, utils.CreateResourceName(ctx, targetGroupName), &elb.TargetGroupArgs{
		Port:       pulumi.Int(80),
		Protocol:   pulumi.String("HTTP"),
		TargetType: pulumi.String("ip"),
		VpcId:      vpcId,
		HealthCheck: &elb.TargetGroupHealthCheckArgs{
			Path: pulumi.StringPtr(hcPath),
		},
	})
}

// NewRdsInstance creates a RDS PostgreSQL instance - remember database Name in RDSArgs can only contain alphanumeric values
func NewRdsInstance(ctx *pulumi.Context, name string, username string, password string, instanceClass string, allocatedStorage int, engineVersion string, dbSubnetGroupId pulumi.StringPtrInput, secGroupIds pulumi.StringArrayInput) (*rds.Instance, error) {
	dbName := utils.CreateResourceName(ctx, strings.ToLower(name))
	databaseName := strings.Replace(dbName, "-", "", -1)
	return rds.NewInstance(ctx, dbName, &rds.InstanceArgs{
		ApplyImmediately:    pulumi.BoolPtr(true),
		SkipFinalSnapshot:   pulumi.BoolPtr(true),
		Engine:              pulumi.StringPtr("postgres"),
		EngineVersion:       pulumi.StringPtr(engineVersion),
		AllocatedStorage:    pulumi.IntPtr(allocatedStorage),
		InstanceClass:       pulumi.String(instanceClass),
		Name:                pulumi.StringPtr(databaseName),
		Identifier:          pulumi.StringPtr(dbName),
		Username:            pulumi.StringPtr(username),
		Password:            pulumi.StringPtr(password),
		PubliclyAccessible:  pulumi.BoolPtr(true),
		DbSubnetGroupName:   dbSubnetGroupId,
		VpcSecurityGroupIds: secGroupIds,
	}, pulumi.Protect(true))
}

// CreateSNSTopic creates a new SNS topic with the given name and options
func CreateSNSTopic(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*sns.Topic, error) {
	return sns.NewTopic(ctx, utils.CreateResourceName(ctx, name), &sns.TopicArgs{}, opts...)
}

// CreateSQSQueue creates a new SQS queue with the given name and options; and, optionally, a DLQ
func CreateSQSQueue(
	ctx *pulumi.Context,
	name string,
	createDLQ bool,
	policy pulumi.StringPtrInput,
	opts ...pulumi.ResourceOption,
) (*sqs.Queue, *pulumi.StringOutput, error) {

	args := sqs.QueueArgs{}

	if createDLQ {
		dlq, _, err := CreateSQSQueue(ctx, name+"-dlq", false, nil, opts...)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to create DLQ")
		}

		args.RedrivePolicy = pulumi.Sprintf(`{
			"deadLetterTargetArn": "%s",
			"maxReceiveCount": 5
		}`, dlq.Arn)
	}

	if policy != nil {
		args.Policy = policy
	}

	queue, err := sqs.NewQueue(ctx, utils.CreateResourceName(ctx, name), &args, opts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to create SQS queue")
	}

	region, err := aws.GetRegion(ctx, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get current AWS region")
	}

	identity, err := aws.GetCallerIdentity(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to get AWS account ID")
	}

	queueURL := pulumi.Sprintf("https://sqs.%v.amazonaws.com/%v/%v", region.Name, identity.AccountId, queue.Name)

	return queue, &queueURL, nil
}

// CreateSQSToSNSTopicSubscription subscribes the given SQS queue to the given SNS topic
func CreateSQSToSNSTopicSubscription(ctx *pulumi.Context, name string, topicARN pulumi.StringInput, queueARN pulumi.StringInput, opts ...pulumi.ResourceOption) (*sns.TopicSubscription, error) {
	return sns.NewTopicSubscription(ctx, utils.CreateResourceName(ctx, name), &sns.TopicSubscriptionArgs{
		Topic:              topicARN,
		Protocol:           pulumi.String("sqs"),
		Endpoint:           queueARN,
		RawMessageDelivery: pulumi.Bool(true),
	}, opts...)
}

func CreateS3Bucket(ctx *pulumi.Context, bucketName string) (*s3.Bucket, error) {
	cfg := pulumiConfig.New(ctx, "")
	env := cfg.Require("env")

	bucketName = fmt.Sprintf("%s-%s", env, bucketName)

	return s3.NewBucket(ctx, bucketName, &s3.BucketArgs{
		Bucket: pulumi.StringPtr(bucketName),
	})
}

// CreateRolePolicy creates an IAM policy attached to a role
func CreateRolePolicy(ctx *pulumi.Context, role *iam.Role, policyName string, policy pulumi.StringInput) (*iam.RolePolicy, error) {
	return iam.NewRolePolicy(ctx, utils.CreateResourceName(ctx, policyName), &iam.RolePolicyArgs{
		Policy: policy,
		Role:   role.Name,
	})
}

// GetSQSWriterPolicyStatement returns an IAM policy document statement that grants permission to write to an SQS queue.
func GetSQSWriterPolicyStatement(queueArn pulumi.StringInput) iam.GetPolicyDocumentStatementOutput {
	return queueArn.ToStringOutput().ApplyT(func(arn string) iam.GetPolicyDocumentStatement {
		return iam.GetPolicyDocumentStatement{
			Actions:   []string{"sqs:SendMessage"},
			Resources: []string{arn},
		}
	}).(iam.GetPolicyDocumentStatementOutput)
}

// GetSQSReaderPolicyStatement returns an IAM policy document statement that grants permission to read from an SQS queue.
func GetSQSReaderPolicyStatement(queueArn pulumi.StringInput) iam.GetPolicyDocumentStatementOutput {
	return queueArn.ToStringOutput().ApplyT(func(arn string) iam.GetPolicyDocumentStatement {
		return iam.GetPolicyDocumentStatement{
			Actions:   []string{"sqs:ReceiveMessage", "sqs:DeleteMessage"},
			Resources: []string{arn},
		}
	}).(iam.GetPolicyDocumentStatementOutput)
}
