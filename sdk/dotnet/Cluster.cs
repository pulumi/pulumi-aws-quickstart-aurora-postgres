// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.AwsQuickStartAuroraPostgres
{
    [AwsQuickStartAuroraPostgresResourceType("aws-quickstart-aurora-postgres:index:Cluster")]
    public partial class Cluster : Pulumi.ComponentResource
    {
        /// <summary>
        /// Create a Cluster resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Cluster(string name, ClusterArgs args, ComponentResourceOptions? options = null)
            : base("aws-quickstart-aurora-postgres:index:Cluster", name, args ?? new ClusterArgs(), MakeResourceOptions(options, ""), remote: true)
        {
        }

        private static ComponentResourceOptions MakeResourceOptions(ComponentResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new ComponentResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = ComponentResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
    }

    public sealed class ClusterArgs : Pulumi.ResourceArgs
    {
        [Input("availabilityZoneNames", required: true)]
        private List<Input<string>>? _availabilityZoneNames;

        /// <summary>
        /// List of Availability Zone names to use to create the DB Cluster.
        /// </summary>
        public List<Input<string>> AvailabilityZoneNames
        {
            get => _availabilityZoneNames ?? (_availabilityZoneNames = new List<Input<string>>());
            set => _availabilityZoneNames = value;
        }

        /// <summary>
        /// Set this parameter to true if you want to enable your DB
        /// instances to receive minor DB engine version upgrades
        /// automatically when upgrades become available.
        /// </summary>
        [Input("dbAutoMinorVersionUpgrade")]
        public bool? DbAutoMinorVersionUpgrade { get; set; }

        /// <summary>
        /// The number of days to retain automatic database snapshots.
        /// To disable automatic backups, set this parameter to 0. Default is 35 days
        /// </summary>
        [Input("dbBackupRetentionPeriod")]
        public int? DbBackupRetentionPeriod { get; set; }

        /// <summary>
        /// The number of days to retain automatic database snapshots.
        /// To disable automatic backups, set this parameter to 0.
        /// </summary>
        [Input("dbEncryptedEnabled")]
        public Input<bool>? DbEncryptedEnabled { get; set; }

        /// <summary>
        /// The number of days to retain automatic database snapshots.
        /// To disable automatic backups, set this parameter to 0.
        /// </summary>
        [Input("dbEngineVersion", required: true)]
        public string DbEngineVersion { get; set; } = null!;

        /// <summary>
        /// The DB (compute and memory capacity) class for the database
        /// instances.
        /// </summary>
        [Input("dbInstanceClass", required: true)]
        public string DbInstanceClass { get; set; } = null!;

        [Input("dbMasterPassword", required: true)]
        private string? _dbMasterPassword;

        /// <summary>
        /// The password for the database administrator account (8-64
        /// character string)
        /// </summary>
        public string? DbMasterPassword
        {
            get => _dbMasterPassword;
            set
            {
                var emptySecret = Output.CreateSecret(0);
                _dbMasterPassword = Output.Tuple<string?, int>(value, emptySecret).Apply(t => t.Item1);
            }
        }

        /// <summary>
        /// The user name for the database administrator account. This is
        /// an alphanumeric string of 1-16 characters. The user name
        /// must start with an uppercase or lowercase letter (A-Z, a-z).
        /// </summary>
        [Input("dbMasterUsername", required: true)]
        public string DbMasterUsername { get; set; } = null!;

        /// <summary>
        /// The name of the Aurora DB to provision. This is an
        /// alphanumeric string of 5-64 characters.
        /// </summary>
        [Input("dbName", required: true)]
        public string DbName { get; set; } = null!;

        /// <summary>
        /// The number of db instances to launch as part of the cluster. Defaults to 1.
        /// </summary>
        [Input("dbNumDbClusterInstances")]
        public double? DbNumDbClusterInstances { get; set; }

        /// <summary>
        /// The family of the DB parameter group (e.g. aurora-postgresql11).
        /// </summary>
        [Input("dbParameterGroupFamily", required: true)]
        public string DbParameterGroupFamily { get; set; } = null!;

        /// <summary>
        /// The port that you want to access the database through. The DB
        /// instance will listen on this port for connections. This value
        /// must be in the range 1115-65535. Default is 5432
        /// </summary>
        [Input("dbPort")]
        public double? DbPort { get; set; }

        /// <summary>
        /// The ID of the custom security group you want to use in your
        /// existing VPC (e.g., sg-7f16e910).
        /// </summary>
        [Input("dbSecurityGroupID")]
        public Input<string>? DbSecurityGroupID { get; set; }

        /// <summary>
        /// Set this parameter to `false` if you want to disable Amazon
        /// Aurora Cluster and Instance level event subscriptions. You
        /// might want to disable it if you are testing or running
        /// continuous integration (CI) processes.
        /// </summary>
        [Input("enableEventSubscription")]
        public bool? EnableEventSubscription { get; set; }

        /// <summary>
        /// The ID of the private subnet in Availability Zone 1 in your
        /// existing VPC (e.g., subnet-a0246dcd).
        /// </summary>
        [Input("privateSubnetID1", required: true)]
        public Input<string> PrivateSubnetID1 { get; set; } = null!;

        /// <summary>
        /// The ID of the private subnet in Availability Zone 2 in your
        /// existing VPC (e.g., subnet-b58c3d67).
        /// </summary>
        [Input("privateSubnetID2", required: true)]
        public Input<string> PrivateSubnetID2 { get; set; } = null!;

        /// <summary>
        /// The email that is used to configure an SNS topic for sending
        /// CloudWatch alarms and Amazon RDS event notifications. This
        /// must be a valid email address. Required if enableEventSubscription is true.
        /// </summary>
        [Input("snsNotificationEmail")]
        public string? SnsNotificationEmail { get; set; }

        /// <summary>
        /// The ID of your existing VPC (e.g., vpc-0343606e) where you
        /// want to deploy the Aurora database.
        /// </summary>
        [Input("vpcID", required: true)]
        public Input<string> VpcID { get; set; } = null!;

        public ClusterArgs()
        {
            DbAutoMinorVersionUpgrade = false;
            DbEncryptedEnabled = true;
            EnableEventSubscription = true;
        }
    }
}