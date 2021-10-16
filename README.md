# Pulumi AWS Aurora Postgres

Easily deploy an Aurora Postgres database with accompanying features like alarms, logging, encryption, and multi-AZ redundancy. This component is based on the best practices recommended by AWS in the [Modular architecture for Amazon Aurora PostgreSQL Quickstart](https://aws.amazon.com/quickstart/architecture/aurora-postgresql/)


# Examples

See the `/examples` directory for more

Go:
```go
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
```

Typescript
```typescript

const multiAvaialabilityZoneAuroraCluster = new aurora.Cluster("example-aurora-cluster", {
  vpcID: multiAvailabilityZoneVpc.vpcID,
  dbName: "myDemoDatabase",
  dbEngineVersion: "11.9",
  dbInstanceClass: "db.t3.medium",
  availabilityZoneNames: ["us-east-1a", "us-east-1b"],
  dbNumDbClusterInstances: 2,
  dbMasterUsername: "mainuser",
  snsNotificationEmail: emailAddress,
  enableEventSubscription: true,
  dbMasterPassword: dbPasswordSecret,
  dbParameterGroupFamily: "aurora-postgresql11",
  privateSubnetID1: multiAvailabilityZoneVpc.privateSubnetIDs.apply(x => x![0]),
  privateSubnetID2: multiAvailabilityZoneVpc.privateSubnetIDs.apply(x => x![1]),
});

```