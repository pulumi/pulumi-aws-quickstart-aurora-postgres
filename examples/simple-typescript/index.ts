import * as pulumi from "@pulumi/pulumi";
import * as vpc from "@pulumi/aws-quickstart-vpc";
import * as aurora from "@pulumi/aws-quickstart-aurora-postgres";

const dbMasterPassword = "?SimpleExamplePassword12345?"
const dbPasswordSecret = pulumi.secret(dbMasterPassword);

const databaseNotificationEmail = "your-email@example.com";

const multiAvailabilityZoneVpc = new vpc.Vpc("example-aurora-vpc", {
    cidrBlock: "10.0.0.0/16",
    availabilityZoneConfig: [{
        availabilityZone: "us-east-1a",
        publicSubnetCidr: "10.0.128.0/20",
        privateSubnetACidr: "10.0.32.0/19",
    }, {
        availabilityZone: "us-east-1b",
        privateSubnetACidr: "10.0.64.0/19",
    }]
})


const multiAvaialabilityZoneAuroraCluster = new aurora.Cluster("example-aurora-cluster", {
  vpcID: multiAvailabilityZoneVpc.vpcID,
  dbName: "myDemoDatabase",
  dbEngineVersion: "11.9",
  dbInstanceClass: "db.t3.medium",
  availabilityZoneNames: ["us-east-1a", "us-east-1b"],
  dbNumDbClusterInstances: 2,
  dbMasterUsername: "mainuser",
  snsNotificationEmail: databaseNotificationEmail,
  enableEventSubscription: true,
  dbMasterPassword: dbPasswordSecret,
  dbParameterGroupFamily: "aurora-postgresql11",
  privateSubnetID1: multiAvailabilityZoneVpc.privateSubnetIDs.apply(x => x![0]),
  privateSubnetID2: multiAvailabilityZoneVpc.privateSubnetIDs.apply(x => x![1]),
});