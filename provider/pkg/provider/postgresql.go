// Copyright 2016-2021, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"encoding/json"
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/kms"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/sns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ClusterArgs struct {
	DatabaseName              string `pulumi:"dbName"`
	DBAutoMinorVersionUpgrade bool   `pulumi:"dbAutoMinorVersionUpgrade"`
	DbBackupRetentionPeriod   int32  `pulumi:"dbBackupRetentionPeriod"`
	DbEngineVersion           string `pulumi:"dbEngineVersion"`
	DBInstanceClass           string `pulumi:"dbInstanceClass"`
	DBMasterUsername          string `pulumi:"dbMasterUsername"`
	DBMasterUserPassword      string `pulumi:"dbMasterUserPassword"`
	DbPort                    int32  `pulumi:"dbPort"`
	DbEncryptionEnabled       bool   `pulumi:"dbEncryptedEnabled"`
	DbParameterGroupFamily    string `pulumi:"dbParameterGroupFamily"`
	NumDbClusterInstances     int    `pulumi:"numNumDbClusterInstances"`

	EnableEventSubscription bool   `pulumi:"enableEventSubscription"`
	SnsNotificationEmail    string `pulumi:"snsNotificationEmail"`

	AvailabilityZoneNames []string           `pulumi:"availabilityZoneNames"`
	VpcID                 pulumi.StringInput `pulumi:"vpcID"`
	PrivateSubnetID1      pulumi.StringInput `pulumi:"privateSubnetID1"`
	PrivateSubnetID2      pulumi.StringInput `pulumi:"privateSubnetID2"`
	DbSecurityGroupID     pulumi.StringInput `pulumi:"dbSecurityGroupID"`
}

type Cluster struct {
	pulumi.ResourceState
}

func NewCluster(ctx *pulumi.Context,
	name string, args *ClusterArgs, opts ...pulumi.ResourceOption) (*Cluster, error) {
	if args == nil {
		args = &ClusterArgs{}
	}

	component := &Cluster{}
	err := ctx.RegisterComponentResource("aws-quickstart-aurora-postgres:index:Cluster", name, component, opts...)
	if err != nil {
		return nil, err
	}

	account, callerIdentityErr := aws.GetCallerIdentity(ctx, nil, nil)
	if callerIdentityErr != nil {
		return nil, callerIdentityErr
	}

	var key *kms.Key
	if args.DbEncryptionEnabled {
		keyPolicy, keyPolicyErr := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Id":      ctx.Stack(),
			"Statement": []map[string]interface{}{
				{
					"Principal": map[string]interface{}{
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
			return nil, keyPolicyErr
		}

		kmsKey, kmsKeyErr := kms.NewKey(ctx, fmt.Sprintf("%s-database-kms-key", name), &kms.KeyArgs{
			Policy: pulumi.String(keyPolicy),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("database-kms-key"),
			},
		})
		if kmsKeyErr != nil {
			return nil, kmsKeyErr
		}

		_, kmsKeyAliasErr := kms.NewAlias(ctx, fmt.Sprintf("%s-database-kms-key-alias", name), &kms.AliasArgs{
			Name:        pulumi.String(fmt.Sprintf("alias/%s", name)),
			TargetKeyId: kmsKey.ID(),
		})
		if kmsKeyAliasErr != nil {
			return nil, kmsKeyAliasErr
		}

		key = kmsKey
	}

	subnetGroup, subnetGroupErr := rds.NewSubnetGroup(ctx, fmt.Sprintf("%s-db-subnet-group", name), &rds.SubnetGroupArgs{
		Description: pulumi.Sprintf("Aurora PG Database Instance Subnet Group for %s", name),
		SubnetIds: pulumi.StringArray{
			args.PrivateSubnetID1,
			args.PrivateSubnetID2,
		},
	})
	if subnetGroupErr != nil {
		return nil, subnetGroupErr
	}

	dbParameterGroup, dbParamGroupErr := rds.NewParameterGroup(ctx, fmt.Sprintf("%s-db-parameter-group", name), &rds.ParameterGroupArgs{
		Family:      pulumi.String(args.DbParameterGroupFamily),
		Description: pulumi.Sprintf("Aurora PG Database Instance Parameter Group for %s", name),
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
	if dbParamGroupErr != nil {
		return nil, dbParamGroupErr
	}

	clusterParameterGroup, parameterGroupErr := rds.NewClusterParameterGroup(ctx, fmt.Sprintf("%s-parameter-group", name), &rds.ClusterParameterGroupArgs{
		Family: pulumi.String(args.DbParameterGroupFamily),
		Parameters: rds.ClusterParameterGroupParameterArray{
			rds.ClusterParameterGroupParameterArgs{
				Name:  pulumi.String("rds.force_ssl"),
				Value: pulumi.String("1"),
			},
		},
	})
	if parameterGroupErr != nil {
		return nil, parameterGroupErr
	}

	azs := pulumi.StringArray{}
	for _, az := range args.AvailabilityZoneNames {
		azs = append(azs, pulumi.String(az))
	}

	retentionPeriod := args.DbBackupRetentionPeriod
	if retentionPeriod == 0 {
		retentionPeriod = 32
	}

	port := args.DbPort
	if port == 0 {
		port = 5432
	}

	clusterArgs := rds.ClusterArgs{
		AvailabilityZones:           azs,
		DbClusterParameterGroupName: clusterParameterGroup.Name,
		BackupRetentionPeriod:       pulumi.Int(retentionPeriod),
		SkipFinalSnapshot:           pulumi.Bool(true),
		DatabaseName:                pulumi.String(args.DatabaseName),
		Engine:                      pulumi.String("aurora-postgresql"),
		EngineVersion:               pulumi.String(args.DbEngineVersion),
		MasterUsername:              pulumi.String(args.DBMasterUsername),
		MasterPassword:              pulumi.String(args.DBMasterUserPassword),
		Port:                        pulumi.Int(port),
		StorageEncrypted:            pulumi.Bool(args.DbEncryptionEnabled),
		DbSubnetGroupName:           subnetGroup.Name,
	}
	if args.DbEncryptionEnabled {
		clusterArgs.KmsKeyId = key.Arn
	}
	if args.DbSecurityGroupID != nil {
		clusterArgs.VpcSecurityGroupIds = pulumi.StringArray{
			args.DbSecurityGroupID,
		}
	}
	dbCluster, clusterErr := rds.NewCluster(ctx, fmt.Sprintf("%s-postgresql-cluster", name), &clusterArgs)
	if clusterErr != nil {
		return nil, clusterErr
	}

	instanceCount := args.NumDbClusterInstances
	if instanceCount == 0 {
		instanceCount = 1
	}

	var topic *sns.Topic
	if args.EnableEventSubscription {
		snsTopic, snsTopicErr := sns.NewTopic(ctx, "sns-topic", &sns.TopicArgs{
			DisplayName: pulumi.String(args.DatabaseName),
		})
		if snsTopicErr != nil {
			return nil, snsTopicErr
		}

		_, snsTopicSubscriptionErr := sns.NewTopicSubscription(ctx, "sns-topic-subscription", &sns.TopicSubscriptionArgs{
			Topic:    snsTopic.ID(),
			Endpoint: pulumi.String(args.SnsNotificationEmail),
			Protocol: pulumi.String("email"),
		})
		if snsTopicSubscriptionErr != nil {
			return nil, snsTopicSubscriptionErr
		}

		topic = snsTopic
	}

	for i := 0; i < instanceCount; i++ {
		dbInstance, dbInstanceErr := rds.NewClusterInstance(ctx, fmt.Sprintf("%s-aurora-database-%d", name, i), &rds.ClusterInstanceArgs{
			AutoMinorVersionUpgrade: pulumi.Bool(args.DBAutoMinorVersionUpgrade),
			ClusterIdentifier:       dbCluster.ID(),
			DbParameterGroupName:    dbParameterGroup.ID(),
			InstanceClass:           pulumi.String(args.DBInstanceClass),
			Engine:                  pulumi.String("aurora-postgresql"),
			EngineVersion:           pulumi.String(args.DbEngineVersion),
			PubliclyAccessible:      pulumi.Bool(false),
			DbSubnetGroupName:       subnetGroup.Name,
		})
		if dbInstanceErr != nil {
			return nil, dbInstanceErr
		}

		if args.EnableEventSubscription {
			_, cpuUtilization1AlarmErr := cloudwatch.NewMetricAlarm(ctx, fmt.Sprintf("%s-cpu-alarm-%d", name, i), &cloudwatch.MetricAlarmArgs{
				ActionsEnabled:   pulumi.Bool(true),
				AlarmActions:     pulumi.Array{topic.ID()},
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
			if cpuUtilization1AlarmErr != nil {
				return nil, cpuUtilization1AlarmErr
			}

			_, maxUsedTxIDsAlarm1Err := cloudwatch.NewMetricAlarm(ctx, fmt.Sprintf("%s-max-used-tx-alarm-%d", name, i), &cloudwatch.MetricAlarmArgs{
				ActionsEnabled:   pulumi.Bool(true),
				AlarmActions:     pulumi.Array{topic.ID()},
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
			if maxUsedTxIDsAlarm1Err != nil {
				return nil, maxUsedTxIDsAlarm1Err
			}

			_, freeLocalStorageAlarm1Err := cloudwatch.NewMetricAlarm(ctx, fmt.Sprintf("%s-free-local-storage-alarm-%d", name, i), &cloudwatch.MetricAlarmArgs{
				ActionsEnabled:   pulumi.Bool(true),
				AlarmActions:     pulumi.Array{topic.ID()},
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
			if freeLocalStorageAlarm1Err != nil {
				return nil, freeLocalStorageAlarm1Err
			}

			_, clusterEventSubscriptionErr := rds.NewEventSubscription(ctx, fmt.Sprintf("%s-cluster-event-subscription-%d", name, i), &rds.EventSubscriptionArgs{
				SnsTopic:        topic.ID(),
				SourceType:      pulumi.String("db-cluster"),
				EventCategories: pulumi.ToStringArray([]string{"failover", "failure", "notification"}),
				SourceIds: pulumi.StringArray{
					dbCluster.ID(),
				},
			})
			if clusterEventSubscriptionErr != nil {
				return nil, clusterEventSubscriptionErr
			}

			_, instanceEventSubscriptionErr := rds.NewEventSubscription(ctx, fmt.Sprintf("%s-instance-event-subscription-%d", name, i), &rds.EventSubscriptionArgs{
				SnsTopic:        topic.ID(),
				SourceType:      pulumi.String("db-instance"),
				EventCategories: pulumi.ToStringArray([]string{"availability", "configuration change", "deletion", "failover", "failure", "maintenance", "notification", "recovery"}),
				SourceIds: pulumi.StringArray{
					dbInstance.ID(),
				},
			})
			if instanceEventSubscriptionErr != nil {
				return nil, instanceEventSubscriptionErr
			}

			_, dbParameterGroupEventSubscriptionErr := rds.NewEventSubscription(ctx, fmt.Sprintf("%s-parameter-group-event-subscription-%d", name, i), &rds.EventSubscriptionArgs{
				SnsTopic:        topic.ID(),
				SourceType:      pulumi.String("db-parameter-group"),
				EventCategories: pulumi.ToStringArray([]string{"configuration change"}),
				SourceIds: pulumi.StringArray{
					dbParameterGroup.ID(),
				},
			})
			if dbParameterGroupEventSubscriptionErr != nil {
				return nil, dbParameterGroupEventSubscriptionErr
			}
		}
	}

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{}); err != nil {
		return nil, err
	}

	return component, nil
}
