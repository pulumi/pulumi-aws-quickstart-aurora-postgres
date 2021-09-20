package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/kms"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an AWS resource (S3 Bucket)
		bucket, err := s3.NewBucket(ctx, "my-bucket", nil)
		if err != nil {
			return err
		}

		/*******************************
		 * Networking
		 *******************************/
		databaseName := "tempDbName"

		// entireInternet := "0.0.0.0/0"
		// dbAccessCidr := "0.0.0.0/0"
		dbPort := 5432
		// @fixme - add family map from:
		//					https://github.com/aws-quickstart/quickstart-amazon-aurora-postgresql/blob/9b20909f8a8139b84b149788d9124bf5a0321c62/templates/aurora_postgres.template.yaml#L9
		family := "aurora-postgresql11"
		engineVersion := "11.9"
		isStorageEncrypted := true

		// "engine" is NOT a parameter and can't be changed for this quickstart
		engine := "aurora-postgresql"

		// Cluster parameters
		backupRetentionPeriod := 5

		// Aurora database
		dbMasterUsername := "foo"
		dbMasterPassword := "example-dummy-that-will-be-changed-later"

		dbAutoMinorVersionUpgrade := true
		// @fixme - choose db instance class
		dbInstanceClass := "db.t3.medium"

		notificationEmail := "dummyemail1122334745566@example.org"

		/**********
		 * SNS Topic
		 **********/
		snsTopic, snsTopicErr := sns.NewTopic(ctx, "sns-topic", &sns.TopicArgs{
			DisplayName: pulumi.String("Database Name"),
		})

		if snsTopicErr != nil {
			return snsTopicErr
		}

		_, snsTopicSubscriptionErr := sns.NewTopicSubscription(ctx, "sns-topic-subscription", &sns.TopicSubscriptionArgs{
			Topic:    snsTopic.ID(),
			Endpoint: pulumi.String(notificationEmail),
			Protocol: pulumi.String("email"),
		})

		if snsTopicSubscriptionErr != nil {
			return snsTopicSubscriptionErr
		}

		account, callerIdentityErr := aws.GetCallerIdentity(ctx, nil, nil)

		if callerIdentityErr != nil {
			return callerIdentityErr
		}

		// KMS Encryption Key
		keyPolicy, keyPolicyErr := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Id":      ctx.Stack(),
			"Statement": []map[string]interface{}{
				{
					"Principal": map[string]interface{}{
						// @todo - get stack account, remove the hardcoded id
						"AWS": "arn:aws:iam::" + account.AccountId + ":root",
					},
					"Action": []string{
						"kms:*",
					},
					"Effect":   "Allow",
					"Resource": "*",
				},
			},
		})

		if keyPolicyErr != nil {
			return keyPolicyErr
		}

		kmsKey, kmsKeyErr := kms.NewKey(ctx, "database-kms-key", &kms.KeyArgs{
			// @todo - open a Pulumi issue to add the "DeletionPolicy" functionality
			// DeletionPolicy: pulumi.String("Retain"),
			Policy: pulumi.String(keyPolicy),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("database-kms-key"),
			},
		})

		if kmsKeyErr != nil {
			return kmsKeyErr
		}

		_, kmsKeyAliasErr := kms.NewAlias(ctx, "database-kms-key-alias", &kms.AliasArgs{
			// @fixme - replace "stack-name" with the actual stack name
			Name:        pulumi.String("alias/stack-name"),
			TargetKeyId: kmsKey.ID(),
		})

		if kmsKeyAliasErr != nil {
			return kmsKeyAliasErr
		}

		// @fixme - vpcId for security group
		// dbSecurityGroup, err := ec2.NewSecurityGroup(ctx, "allowTls", &ec2.SecurityGroupArgs{
		// 	Description: pulumi.String("Allow TLS inbound traffic to database port"),
		// 	// @fixme - use other vpc id
		// 	VpcId: pulumi.Any(aws_vpc.Main.Id),
		// 	Ingress: ec2.SecurityGroupIngressArray{
		// 		&ec2.SecurityGroupIngressArgs{
		// 			Description: pulumi.String("TLS DB port for db access cidr"),
		// 			FromPort:    pulumi.Int(dbPort),
		// 			ToPort:      pulumi.Int(dbPort),
		// 			Protocol:    pulumi.String("tcp"),
		// 			CidrBlocks: pulumi.StringArray{
		// 				pulumi.String(dbAccessCidr),
		// 			},
		// 		},
		// 		&ec2.SecurityGroupIngressArgs{
		// 			Description: pulumi.String("TLS DB port for db access cidr"),
		// 			FromPort:    pulumi.Int(0),
		// 			ToPort:      pulumi.Int(0),
		// 			Protocol:    pulumi.String("all"),
		// 			S
		// 		},
		// 	},
		// 	Egress: ec2.SecurityGroupEgressArray{
		// 		&ec2.SecurityGroupEgressArgs{
		// 			// protocol "all" ignores the from port and to port.
		// 			FromPort: pulumi.Int(0),
		// 			ToPort:   pulumi.Int(0),
		// 			Protocol: pulumi.String("all"),
		// 			CidrBlocks: pulumi.StringArray{
		// 				pulumi.String(entireInternet),
		// 			},
		// 		},
		// 	},
		// 	Tags: pulumi.StringMap{
		// 		"Name": pulumi.String("allow_db_access"),
		// 	},
		// })

		dbParameterGroup, _ := rds.NewParameterGroup(ctx, "db-parameter-group", &rds.ParameterGroupArgs{
			Family:      pulumi.String(family),
			Description: pulumi.String("Aurora PG Database Instance Parameter Group for Cloudformation Stack " + ctx.Stack()),
			Parameters: rds.ParameterGroupParameterArray{
				rds.ParameterGroupParameterArgs{
					Name:  pulumi.String("log_rotation_age"),
					Value: pulumi.String("1440"),
				},
				rds.ParameterGroupParameterArgs{
					Name:  pulumi.String("log_rotation_size"),
					Value: pulumi.String("102400"),
				},
			},
		})

		clusterParameterGroup, parameterGroupErr := rds.NewClusterParameterGroup(ctx, "parameter-group", &rds.ClusterParameterGroupArgs{
			Family: pulumi.String(family),
			Parameters: rds.ClusterParameterGroupParameterArray{
				rds.ClusterParameterGroupParameterArgs{
					Name:  pulumi.String("rds.force_ssl"),
					Value: pulumi.String("1"),
				},
			},
		})

		if parameterGroupErr != nil {
			return parameterGroupErr
		}

		dbCluster, clusterErr := rds.NewCluster(ctx, "postgresql", &rds.ClusterArgs{
			AvailabilityZones: pulumi.StringArray{
				pulumi.String("us-east-1a"),
				pulumi.String("us-east-1b"),
				pulumi.String("us-east-1c"),
			},
			DbClusterParameterGroupName: clusterParameterGroup.Name,
			BackupRetentionPeriod:       pulumi.Int(backupRetentionPeriod),
			// DBSubnetGroupname
			// EnableCloudwatchLogsExports
			SkipFinalSnapshot: pulumi.Bool(true),
			DatabaseName:      pulumi.String(databaseName),
			Engine:            pulumi.String(engine),
			EngineVersion:     pulumi.String(engineVersion),
			KmsKeyId:          kmsKey.Arn,
			MasterUsername:    pulumi.String(dbMasterUsername),
			MasterPassword:    pulumi.String(dbMasterPassword),
			Port:              pulumi.Int(dbPort),
			StorageEncrypted:  pulumi.Bool(isStorageEncrypted),
			// tags
			// @fixme - create security group and uncomment this line
			// VpcSecurityGroupIds: pulumi.StringArray{dbSecurityGroup.ID()},
		})

		if clusterErr != nil {
			return clusterErr
		}

		dbInstance1, dbInstance1Err := rds.NewClusterInstance(ctx, "aurora-database-1", &rds.ClusterInstanceArgs{
			// AllocatedStorage was not included in the quickstart guide but error was thrown saying: "allocated_storage": required field is not set"
			// AllocatedStorage is ignored by Aurora, so we're setting it to 1 here as a dummy default.
			//   See: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-database-instance.html#cfn-rds-dbinstance-allocatedstorage
			AutoMinorVersionUpgrade: pulumi.Bool(dbAutoMinorVersionUpgrade),
			// @fixme - DbClusterIdentifier not found
			ClusterIdentifier:    dbCluster.ID(),
			DbParameterGroupName: dbParameterGroup.ID(),
			InstanceClass:        pulumi.String(dbInstanceClass),
			Engine:               pulumi.String(engine),
			EngineVersion:        pulumi.String(engineVersion),
			PubliclyAccessible:   pulumi.Bool(false),
			// @fixme add tags
			// Tags
		})

		_, cpuUtilization1AlarmErr := cloudwatch.NewMetricAlarm(ctx, "cpu-alarm", &cloudwatch.MetricAlarmArgs{
			ActionsEnabled:   pulumi.Bool(true),
			AlarmActions:     pulumi.Array{snsTopic.ID()},
			AlarmDescription: pulumi.String("CPU_Utilization"),
			Dimensions: pulumi.StringMap{
				"DBInstanceIdentifier": dbInstance1.ID(),
			},
			MetricName:         pulumi.String("CPUUtilization"),
			Statistic:          pulumi.String("Maximum"),
			Namespace:          pulumi.String("AWS/RDS"),
			Threshold:          pulumi.Float64Ptr(80),
			Unit:               pulumi.String("Percent"),
			ComparisonOperator: pulumi.String("GreaterThanOrEqualToThreshold"),
			Period:             pulumi.Int(60),
			EvaluationPeriods:  pulumi.Int(5),
			TreatMissingData:   pulumi.String("notBreaching"),
		})

		if cpuUtilization1AlarmErr != nil {
			return cpuUtilization1AlarmErr
		}

		_, maxUsedTxIDsAlarm1Err := cloudwatch.NewMetricAlarm(ctx, "max-used-tx-alarm", &cloudwatch.MetricAlarmArgs{
			ActionsEnabled:   pulumi.Bool(true),
			AlarmActions:     pulumi.Array{snsTopic.ID()},
			AlarmDescription: pulumi.String("Maximum Used Transaction IDs"),
			Dimensions: pulumi.StringMap{
				"DBInstanceIdentifier": dbInstance1.ID(),
			},
			MetricName:         pulumi.String("MaximumUsedTransactionIDs"),
			Statistic:          pulumi.String("Average"),
			Namespace:          pulumi.String("AWS/RDS"),
			Threshold:          pulumi.Float64Ptr(600000000),
			Unit:               pulumi.String("Count"),
			ComparisonOperator: pulumi.String("GreaterThanOrEqualToThreshold"),
			Period:             pulumi.Int(60),
			EvaluationPeriods:  pulumi.Int(5),
			TreatMissingData:   pulumi.String("notBreaching"),
		})

		if maxUsedTxIDsAlarm1Err != nil {
			return maxUsedTxIDsAlarm1Err
		}

		_, freeLocalStorageAlarm1Err := cloudwatch.NewMetricAlarm(ctx, "free-local-storage-alarm", &cloudwatch.MetricAlarmArgs{
			ActionsEnabled:   pulumi.Bool(true),
			AlarmActions:     pulumi.Array{snsTopic.ID()},
			AlarmDescription: pulumi.String("Free Local Storage"),
			Dimensions: pulumi.StringMap{
				"DBInstanceIdentifier": dbInstance1.ID(),
			},
			MetricName:         pulumi.String("FreeLocalStorage"),
			Statistic:          pulumi.String("Average"),
			Namespace:          pulumi.String("AWS/RDS"),
			Threshold:          pulumi.Float64Ptr(5368709120),
			Unit:               pulumi.String("Bytes"),
			ComparisonOperator: pulumi.String("LessThanOrEqualToThreshold"),
			Period:             pulumi.Int(60),
			EvaluationPeriods:  pulumi.Int(5),
			TreatMissingData:   pulumi.String("notBreaching"),
		})
		if freeLocalStorageAlarm1Err != nil {
			return freeLocalStorageAlarm1Err
		}

		if dbInstance1Err != nil {
			return dbInstance1Err
		}

		_, clusterEventSubscriptionErr := rds.NewEventSubscription(ctx, "cluster-event-subscription", &rds.EventSubscriptionArgs{
			SnsTopic:        snsTopic.ID(),
			SourceType:      pulumi.String("db-cluster"),
			EventCategories: pulumi.ToStringArray([]string{"failover", "failure", "notification"}),
			SourceIds: pulumi.StringArray{
				dbCluster.ID(),
			},
		})

		if clusterEventSubscriptionErr != nil {
			return clusterEventSubscriptionErr
		}

		_, instanceEventSubscriptionErr := rds.NewEventSubscription(ctx, "instance-event-subscription", &rds.EventSubscriptionArgs{
			SnsTopic:        snsTopic.ID(),
			SourceType:      pulumi.String("db-instance"),
			EventCategories: pulumi.ToStringArray([]string{"availability", "configuration change", "deletion", "failover", "failure", "maintenance", "notification", "recovery"}),
			SourceIds: pulumi.StringArray{
				dbInstance1.ID(),
			},
		})

		if instanceEventSubscriptionErr != nil {
			return instanceEventSubscriptionErr
		}

		_, dbParameterGroupEventSubscriptionErr := rds.NewEventSubscription(ctx, "parameter-group-event-subscription", &rds.EventSubscriptionArgs{
			SnsTopic:        snsTopic.ID(),
			SourceType:      pulumi.String("db-parameter-group"),
			EventCategories: pulumi.ToStringArray([]string{"configuration change"}),
			SourceIds: pulumi.StringArray{
				dbParameterGroup.ID(),
			},
		})

		if dbParameterGroupEventSubscriptionErr != nil {
			return dbParameterGroupEventSubscriptionErr
		}

		// Export the name of the bucket
		ctx.Export("bucketName", bucket.ID())
		return nil
	})
}
