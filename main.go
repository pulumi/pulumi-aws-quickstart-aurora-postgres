package main

import (
	"encoding/json"
	"errors"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/kms"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func dbFamilyMap(ctx string) (string, error) {
	switch ctx {
	case "9.6.16":
		return "aurora-postgresql9.6", nil
	case "9.6.17":
		return "aurora-postgresql9.6", nil
	case "9.6.18":
		return "aurora-postgresql9.6", nil
	case "9.6.19":
		return "aurora-postgresql9.6", nil
	case "10.11":
		return "aurora-postgresql10", nil
	case "10.12":
		return "aurora-postgresql10", nil
	case "10.13":
		return "aurora-postgresql10", nil
	case "10.14":
		return "aurora-postgresql10", nil
	case "11.6":
		return "aurora-postgresql11", nil
	case "11.7":
		return "aurora-postgresql11", nil
	case "11.8":
		return "aurora-postgresql11", nil
	case "11.9":
		return "aurora-postgresql11", nil
	case "12.4":
		return "aurora-postgresql12", nil
	default:
		return "", errors.New("aurora engine version not supported")
		// throw error
	}
}

type PostgresAuroraInstanceArgs struct {
	namespace                  string
	autoMinorVersioningUpgrade bool
	instanceClass              string
	engine                     string
	engineVersion              string
	snsTopic                   *sns.Topic
	dbCluster                  *rds.Cluster
	dbParameterGroup           *rds.ParameterGroup
}

func provisionInstance(ctx *pulumi.Context, postgresAuroraInstanceArgs *PostgresAuroraInstanceArgs) (clusterInstance *rds.ClusterInstance, err error) {
	dbInstance, dbInstanceErr := rds.NewClusterInstance(ctx, "aurora-database-"+postgresAuroraInstanceArgs.namespace, &rds.ClusterInstanceArgs{
		AutoMinorVersionUpgrade: pulumi.Bool(postgresAuroraInstanceArgs.autoMinorVersioningUpgrade),
		ClusterIdentifier:       postgresAuroraInstanceArgs.dbCluster.ID(),
		DbParameterGroupName:    postgresAuroraInstanceArgs.dbParameterGroup.ID(),
		InstanceClass:           pulumi.String(postgresAuroraInstanceArgs.instanceClass),
		Engine:                  pulumi.String(postgresAuroraInstanceArgs.engine),
		EngineVersion:           pulumi.String(postgresAuroraInstanceArgs.engineVersion),
		PubliclyAccessible:      pulumi.Bool(false),

		// @fixme add tags
		// Tags
	})

	_, cpuUtilizationAlarmErr := cloudwatch.NewMetricAlarm(ctx, "cpu-alarm-"+postgresAuroraInstanceArgs.namespace, &cloudwatch.MetricAlarmArgs{
		ActionsEnabled:   pulumi.Bool(true),
		AlarmActions:     pulumi.Array{postgresAuroraInstanceArgs.snsTopic.ID()},
		AlarmDescription: pulumi.String("CPU_Utilization"),
		Dimensions: pulumi.StringMap{
			"DBInstanceIdentifier": dbInstance.ID(),
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

	if cpuUtilizationAlarmErr != nil {
		return nil, cpuUtilizationAlarmErr
	}

	_, maxUsedTxIDsAlarmErr := cloudwatch.NewMetricAlarm(ctx, "max-used-tx-alarm-"+postgresAuroraInstanceArgs.namespace, &cloudwatch.MetricAlarmArgs{
		ActionsEnabled:   pulumi.Bool(true),
		AlarmActions:     pulumi.Array{postgresAuroraInstanceArgs.snsTopic.ID()},
		AlarmDescription: pulumi.String("Maximum Used Transaction IDs"),
		Dimensions: pulumi.StringMap{
			"DBInstanceIdentifier": dbInstance.ID(),
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

	if maxUsedTxIDsAlarmErr != nil {
		return nil, maxUsedTxIDsAlarmErr
	}

	_, freeLocalStorageAlarmErr := cloudwatch.NewMetricAlarm(ctx, "free-local-storage-alarm-"+postgresAuroraInstanceArgs.namespace, &cloudwatch.MetricAlarmArgs{
		ActionsEnabled:   pulumi.Bool(true),
		AlarmActions:     pulumi.Array{postgresAuroraInstanceArgs.snsTopic.ID()},
		AlarmDescription: pulumi.String("Free Local Storage"),
		Dimensions: pulumi.StringMap{
			"DBInstanceIdentifier": dbInstance.ID(),
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

	if freeLocalStorageAlarmErr != nil {
		return nil, freeLocalStorageAlarmErr
	}

	if dbInstanceErr != nil {
		return nil, dbInstanceErr
	}

	return dbInstance, nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		/********************************
		 * Networking imports for VPC and Subnets
		 *   This might need to be redone later
		 ********************************/

		// @fixme move these to env variables or function inputs
		pVpcId := "vpc-0c1f4105aea0108b6"
		pDbAvailabilityZones := []string{
			"us-east-1a",
			"us-east-1b",
		}
		pDbSubnetIds := []string{
			"subnet-0d0f0bca649e2ca05",
			"subnet-0de5edd7b1c3410a3",
		}
		// END ENV VARIABLES

		// Parameter Subscription
		// eventSubscription: true | false
		pEnableEventSubscription := true

		/*******************************
		 * Networking
		 *******************************/
		// NOTE: if this is set to nil, then Aurora won't create a database (it will create the instances, just not the logical database)
		pDatabaseName := "tempDbName"

		entireInternet := "0.0.0.0/0"
		pDbAccessCidr := "0.0.0.0/0"

		// "ID of the security group (e.g., sg-0234se). One will be created for you if left empty."
		// var pCustomSecurityGroupId *string = nil

		pDbPort := 5432
		// @fixme - add family map from:
		//					https://github.com/aws-quickstart/quickstart-amazon-aurora-postgresql/blob/9b20909f8a8139b84b149788d9124bf5a0321c62/templates/aurora_postgres.template.yaml#L9
		pEngineVersion := "11.9"

		dbFamily, dbFamilyErr := dbFamilyMap(pEngineVersion)

		if nil != dbFamilyErr {
			return dbFamilyErr
		}

		isStorageEncrypted := false

		engine := "aurora-postgresql"

		// Cluster parameters
		pDbBackupRetentionPeriod := 35

		// Aurora database
		pDbMasterUsername := "pgAdmin"
		pDbMasterUserPassword := "example-dummy-that-will-be-changed-later"

		pDbEnableLogExport := true

		pDbAutoMinorVersionUpgrade := true
		pDbInstanceClass := "db.t3.medium"

		pNotificationEmail := "dummyemail1122334745566@example.org"

		pEnvironmentStage := "environment"
		pApplication := "application"
		pApplicationVersion := "applicationVersion"
		pProjectCostCenter := "projectCostCenter"
		// @fixme - use "public | private | confidentiality | pii/phi"
		pConfidentiality := "public"
		// @fixme - use "hipaa | sox | fips | other | ''
		pCompliance := ""

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
			Endpoint: pulumi.String(pNotificationEmail),
			Protocol: pulumi.String("email"),
		})

		if snsTopicSubscriptionErr != nil {
			return snsTopicSubscriptionErr
		}

		account, callerIdentityErr := aws.GetCallerIdentity(ctx, nil, nil)

		if callerIdentityErr != nil {
			return callerIdentityErr
		}

		// var kmsKey *kms.Key = nil

		if isStorageEncrypted {
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
		}

		// This can be a user parameter input or it can be created by pulumi
		// var dbSecurityGroupId string
		// var createSecurityGroup = true

		// if nil == pCustomSecurityGroupId {
		// 	createSecurityGroup = true
		// }

		// var  *pulumi.StringArray

		dbSecurityGroup, dbSecurityGroupErr := ec2.NewSecurityGroup(ctx, "allowTls", &ec2.SecurityGroupArgs{
			Description: pulumi.String("Allow TLS inbound traffic to database port"),
			VpcId:       pulumi.String(pVpcId),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Description: pulumi.String("TLS DB port for db access cidr"),
					FromPort:    pulumi.Int(pDbPort),
					ToPort:      pulumi.Int(pDbPort),
					Protocol:    pulumi.String("tcp"),
					CidrBlocks: pulumi.StringArray{
						pulumi.String(pDbAccessCidr),
					},
				},
				&ec2.SecurityGroupIngressArgs{
					Description: pulumi.String("Allow traffic with itself"),
					FromPort:    pulumi.Int(0),
					ToPort:      pulumi.Int(0),
					Protocol:    pulumi.String("all"),
					Self:        pulumi.Bool(true),
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					// protocol "all" ignores the from port and to port.
					FromPort: pulumi.Int(0),
					ToPort:   pulumi.Int(0),
					Protocol: pulumi.String("all"),
					CidrBlocks: pulumi.StringArray{
						pulumi.String(entireInternet),
					},
				},
			},
			Tags: pulumi.StringMap{
				"Name": pulumi.String("allow_db_access"),
			},
		})

		if dbSecurityGroupErr != nil {
			return dbSecurityGroupErr
		}

		// dbSecurityGroupSelfIngress, dbSecurityGroupSelfIngressErr := ec2.NewSecurityGroupRule(ctx, "allow-self", &ec2.SecurityGroupRuleArgs{

		// }
		// securityGroupIngress := ec2.NewSecurityGroupIngress(ctx, "")

		// // @fixme - how to cast this id into a string?
		// securityGroupIds :=
		// // }
		// // else {
		// // 	dbSecurityGroupId = *pCustomSecurityGroupId
		// // 	// @fixme - what is the correct way to cast this id?
		// // 	k := pulumi.ToStringArray([]string{dbSecurityGroupId})

		// // 	securityGroupIds = &k
		// // }

		// // dbSecurityGroupSelfIngressRule, dbSecurityGroupSelfIngressRuleErr := ec2.NewSecurityGroupRule(ctx, "self-ingress-rule", &ec2.SecurityGroupRuleArgs{
		// // 		Desription: pulumi.String("TLS DB port for db access cidr"),
		// // 		Self: ,
		// // 		// 	FromPort:    pulumi.Int(0),
		// // 		// 	ToPort:      pulumi.Int(0),
		// // 		// 	Protocol:    pulumi.String("all"),
		// // 		// 	S
		// // })

		dbParameterGroup, _ := rds.NewParameterGroup(ctx, "db-parameter-group", &rds.ParameterGroupArgs{
			Family:      pulumi.String(dbFamily),
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
			Family: pulumi.String(dbFamily),
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

		auroraDbSubnetGroup, auroraDbSubnetGroupErr := rds.NewSubnetGroup(ctx, "aurora-subnet-group", &rds.SubnetGroupArgs{
			Description: pulumi.String("Subnets available for the Amazon Aurora database instances"),
			SubnetIds:   pulumi.ToStringArray(pDbSubnetIds),
		})

		if auroraDbSubnetGroupErr != nil {
			return auroraDbSubnetGroupErr
		}

		// No exports by default
		enabledCloudwatchLogsExports := []string{""}

		if pDbEnableLogExport {
			enabledCloudwatchLogsExports = []string{"postgresql"}
		}

		dbCluster, clusterErr := rds.NewCluster(ctx, "postgresql", &rds.ClusterArgs{
			AvailabilityZones:            pulumi.ToStringArray(pDbAvailabilityZones),
			DbClusterParameterGroupName:  clusterParameterGroup.Name,
			BackupRetentionPeriod:        pulumi.Int(pDbBackupRetentionPeriod),
			DbSubnetGroupName:            auroraDbSubnetGroup.Name,
			EnabledCloudwatchLogsExports: pulumi.ToStringArray(enabledCloudwatchLogsExports),
			SkipFinalSnapshot:            pulumi.Bool(true),
			DatabaseName:                 pulumi.String(pDatabaseName),
			Engine:                       pulumi.String(engine),
			EngineVersion:                pulumi.String(pEngineVersion),
			// @fixme - kms key is inside an "if block". Had trouble defining the variable outside the "if" block and using it here. Need a ternary!
			// KmsKeyId:                     kmsKey.KeyId,
			MasterUsername:      pulumi.String(pDbMasterUsername),
			MasterPassword:      pulumi.String(pDbMasterUserPassword),
			Port:                pulumi.Int(pDbPort),
			StorageEncrypted:    pulumi.Bool(isStorageEncrypted),
			VpcSecurityGroupIds: pulumi.StringArray{dbSecurityGroup.ID()},
			Tags: pulumi.StringMap{
				"EnvironmentStage":   pulumi.String(pEnvironmentStage),
				"Application":        pulumi.String(pApplication),
				"ApplicationVersion": pulumi.String(pApplicationVersion),
				"ProjectCostCenter":  pulumi.String(pProjectCostCenter),
				"Confidentiality":    pulumi.String(pConfidentiality),
				"Compliance":         pulumi.String(pCompliance),
			},
		}, pulumi.DependsOn(nil))

		if clusterErr != nil {
			return clusterErr
		}

		// // @fixme - 1. test that this approach works
		// //          2. move all conditional parameters to this "set after constructor" approach
		// if nil != encryptionKeyArn {
		// 	dbCluster.KmsKeyId = *encryptionKeyArn
		// }

		dbInstance1, dbInstance1Err := provisionInstance(ctx, &PostgresAuroraInstanceArgs{
			namespace:                  "1",
			autoMinorVersioningUpgrade: pDbAutoMinorVersionUpgrade,
			instanceClass:              pDbInstanceClass,
			engine:                     engine,
			engineVersion:              pEngineVersion,
			snsTopic:                   snsTopic,
			dbCluster:                  dbCluster,
			dbParameterGroup:           dbParameterGroup,
		})

		if dbInstance1Err != nil {
			return dbInstance1Err
		}

		dbInstance2, dbInstance2Err := provisionInstance(ctx, &PostgresAuroraInstanceArgs{
			namespace:                  "2",
			autoMinorVersioningUpgrade: pDbAutoMinorVersionUpgrade,
			instanceClass:              pDbInstanceClass,
			engine:                     engine,
			engineVersion:              pEngineVersion,
			snsTopic:                   snsTopic,
			dbCluster:                  dbCluster,
			dbParameterGroup:           dbParameterGroup,
		})

		if dbInstance2Err != nil {
			return dbInstance2Err
		}

		dbInstance3, dbInstance3Err := provisionInstance(ctx, &PostgresAuroraInstanceArgs{
			namespace:                  "3",
			autoMinorVersioningUpgrade: pDbAutoMinorVersionUpgrade,
			instanceClass:              pDbInstanceClass,
			engine:                     engine,
			engineVersion:              pEngineVersion,
			snsTopic:                   snsTopic,
			dbCluster:                  dbCluster,
			dbParameterGroup:           dbParameterGroup,
		})

		if dbInstance3Err != nil {
			return dbInstance3Err
		}

		if pEnableEventSubscription {
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
					dbInstance2.ID(),
					dbInstance3.ID(),
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
		}

		// Export the name of the bucket
		// ctx.Export("bucketName", bucket.ID())
		return nil
	})
}
