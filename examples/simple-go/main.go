package main

import (
	quickstartAuroraPostgres "github.com/pulumi/pulumi-aws-quickstart-aurora-postgres/sdk/go/aws"
	quickstartVpc "github.com/pulumi/pulumi-aws-quickstart-vpc/sdk/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		privateSubnet1ACidr := "10.0.1.0/24"
		privateSubnet2ACidr := "10.0.2.0/24"

		databaseNotificationEmail := "your-email@example.org"
		enableEventSubscription := true

		vpc, err := quickstartVpc.NewVpc(ctx, "simple-vpc", &quickstartVpc.VpcArgs{
			CidrBlock: "10.0.0.0/16",
			AvailabilityZoneConfig: []quickstartVpc.AvailabilityZoneArgs{
				quickstartVpc.AvailabilityZoneArgs{
					AvailabilityZone:   "us-east-1a",
					PrivateSubnetACidr: &privateSubnet1ACidr,
				},
				quickstartVpc.AvailabilityZoneArgs{
					AvailabilityZone:   "us-east-1b",
					PrivateSubnetACidr: &privateSubnet2ACidr,
				},
			},
		})

		if err != nil {
			return err
		}

		cfg := config.New(ctx, "")

		dbNumDbClusterInstances := 2

		_, err = quickstartAuroraPostgres.NewCluster(ctx, "smiple-aurora-postgres", &quickstartAuroraPostgres.ClusterArgs{
			VpcID:                   vpc.VpcID,
			DbEngineVersion:         "11.9",
			DbInstanceClass:         "db.t3.medium",
			AvailabilityZoneNames:   pulumi.ToStringArray([]string{"us-east-1a", "us-east-1b"}),
			DbNumDbClusterInstances: &dbNumDbClusterInstances,
			DbMasterUsername:        "mainuser",
			SnsNotificationEmail:    &databaseNotificationEmail,
			EnableEventSubscription: &enableEventSubscription,
			DbMasterPassword:        cfg.RequireSecret("dbPassword"),
			DbParameterGroupFamily:  "aurora-postgresql11",
			PrivateSubnetID1:        vpc.PrivateSubnetIDs.Index(pulumi.Int(0)),
			PrivateSubnetID2:        vpc.PrivateSubnetIDs.Index(pulumi.Int(1)),
		})

		if err != nil {
			return err
		}

		return nil
	})
}
