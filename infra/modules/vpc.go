package modules

import (
	"log"

	"github.com/EurosportDigital/global-transcoding-platform/infra/resources"
	"github.com/EurosportDigital/global-transcoding-platform/infra/utils"
	"github.com/EurosportDigital/global-transcoding-platform/lib/errors"
	"github.com/pulumi/pulumi-aws/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

type VpcInfo struct {
	VpcId       pulumi.StringPtrInput
	SubnetIds   pulumi.StringArray
	IntraSecGrp pulumi.StringInput
}

func SetupVpc(ctx *pulumi.Context) (*VpcInfo, error) {
	vpcInfo := &VpcInfo{}
	// Create a new private VPC with a given CIDR block.
	vpc, err := resources.CreateVpc(ctx, "vpc", "10.20.0.0/16")
	if err != nil {
		return nil, errors.Wrap(err, "creating vpc")
	}

	vpcInfo.VpcId = vpc.ID().ToStringPtrOutput()

	// Create a subnet attached to the above VPC in AZ a.
	subnet1, err := resources.CreateSubnet(ctx, vpc.ID().ToStringOutput(), "snet1", "10.20.16.0/20", "us-west-2a", true)
	if err != nil {
		return nil, errors.Wrap(err, "creating snet1")
	}

	// Create a subnet attached to the above VPC in AZ b.
	subnet2, err := resources.CreateSubnet(ctx, vpc.ID().ToStringOutput(), "snet2", "10.20.32.0/20", "us-west-2b", true)
	if err != nil {
		return nil, errors.Wrap(err, "creating snet2")
	}

	// Create a subnet attached to the above VPC in AZ c.
	subnet3, err := resources.CreateSubnet(ctx, vpc.ID().ToStringOutput(), "snet3", "10.20.48.0/20", "us-west-2c", true)
	if err != nil {
		return nil, errors.Wrap(err, "creating snet3")
	}

	// Create a subnet attached to the above VPC in AZ d.
	subnet4, err := resources.CreateSubnet(ctx, vpc.ID().ToStringOutput(), "snet4", "10.20.64.0/20", "us-west-2d", true)
	if err != nil {
		return nil, errors.Wrap(err, "creating snet4")
	}

	subnetIds := pulumi.StringArray{subnet1.ID().ToStringOutput(), subnet2.ID().ToStringOutput(), subnet3.ID().ToStringOutput(), subnet4.ID().ToStringOutput()}

	vpcInfo.SubnetIds = subnetIds

	// Create a new Internet Gateway and attach it to our VPC.
	inetGw, err := resources.CreateInternetGateway(ctx, vpc.ID().ToStringPtrOutput(), "gw")
	if err != nil {
		return nil, errors.Wrap(err, "creating internet gateway")
	}

	// Create a new routing table that will allow passing traffic through above created Internet Gateway.
	rt, err := resources.CreateRouteTable(ctx, vpc.ID().ToStringOutput(), inetGw.ID().ToStringPtrOutput(), "pubTbl", "0.0.0.0/0")
	if err != nil {
		return nil, errors.Wrap(err, "creating routing table")
	}

	// Create associations between subnets and our route table
	assoc1, err := ec2.NewRouteTableAssociation(ctx, utils.CreateResourceName(ctx, "rtSnAs1"), &ec2.RouteTableAssociationArgs{
		RouteTableId: rt.ID().ToStringOutput(),
		SubnetId:     subnet1.ID().ToStringOutput(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating rtSnAs1")
	}
	assoc2, err := ec2.NewRouteTableAssociation(ctx, utils.CreateResourceName(ctx, "rtSnAs2"), &ec2.RouteTableAssociationArgs{
		RouteTableId: rt.ID().ToStringOutput(),
		SubnetId:     subnet2.ID().ToStringOutput(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating rtSnAs2")
	}
	assoc3, err := ec2.NewRouteTableAssociation(ctx, utils.CreateResourceName(ctx, "rtSnAs3"), &ec2.RouteTableAssociationArgs{
		RouteTableId: rt.ID().ToStringOutput(),
		SubnetId:     subnet3.ID().ToStringOutput(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating rtSnAs3")
	}
	assoc4, err := ec2.NewRouteTableAssociation(ctx, utils.CreateResourceName(ctx, "rtSnAs4"), &ec2.RouteTableAssociationArgs{
		RouteTableId: rt.ID().ToStringOutput(),
		SubnetId:     subnet4.ID().ToStringOutput(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating rtSnAs4")
	}

	intraVPCSecGroup, err := resources.CreateSecurityGroup(ctx, vpcInfo.VpcId, "intVPCSecGrp", []string{"0.0.0.0/0"}, []string{}, 0, true)
	if err != nil {
		return nil, errors.Wrap(err, "creating intVPCSecGrp")
	}

	vpcInfo.IntraSecGrp = intraVPCSecGroup.ID().ToStringOutput()

	out := pulumi.Sprintf("Associations: %s, %s, %s, %s created.", assoc1.ID().ToStringOutput(), assoc2.ID().ToStringOutput(), assoc3.ID().ToStringOutput(), assoc4.ID().ToStringOutput())
	log.Print(out.ToStringOutput())

	ctx.Export("vpcId", vpc.ID().ToStringOutput())
	ctx.Export("subnetIds", subnetIds)
	ctx.Export("inetGwId", inetGw.ID().ToStringOutput())
	ctx.Export("pubRouteId", rt.ID().ToStringOutput())
	ctx.Export("intraVPCSecGroup", intraVPCSecGroup.ID().ToStringOutput())

	return vpcInfo, nil
}
