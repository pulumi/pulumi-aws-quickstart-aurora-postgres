# coding=utf-8
# *** WARNING: this file was generated by Pulumi SDK Generator. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import warnings
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
from . import _utilities

__all__ = ['ClusterArgs', 'Cluster']

@pulumi.input_type
class ClusterArgs:
    def __init__(__self__, *,
                 availability_zone_names: Sequence[pulumi.Input[str]],
                 db_engine_version: str,
                 db_instance_class: str,
                 db_master_password: pulumi.Input[str],
                 db_master_username: str,
                 db_name: str,
                 db_parameter_group_family: str,
                 private_subnet_id1: pulumi.Input[str],
                 private_subnet_id2: pulumi.Input[str],
                 vpc_id: pulumi.Input[str],
                 db_auto_minor_version_upgrade: Optional[bool] = None,
                 db_backup_retention_period: Optional[int] = None,
                 db_encrypted_enabled: Optional[bool] = None,
                 db_num_db_cluster_instances: Optional[int] = None,
                 db_port: Optional[float] = None,
                 db_security_group_id: Optional[pulumi.Input[str]] = None,
                 enable_event_subscription: Optional[bool] = None,
                 sns_notification_email: Optional[str] = None):
        """
        The set of arguments for constructing a Cluster resource.
        :param Sequence[pulumi.Input[str]] availability_zone_names: List of Availability Zone names to use to create the DB Cluster.
        :param str db_engine_version: The version of the database engine.
        :param str db_instance_class: The DB (compute and memory capacity) class for the database
               instances.
        :param pulumi.Input[str] db_master_password: The password for the database administrator account (8-64
               character string)
        :param str db_master_username: The user name for the database administrator account. This is
               an alphanumeric string of 1-16 characters. The user name
               must start with an uppercase or lowercase letter (A-Z, a-z).
        :param str db_name: The name of the Aurora DB to provision. This is an
               alphanumeric string of 5-64 characters.
        :param str db_parameter_group_family: The family of the DB parameter group (e.g. aurora-postgresql11).
        :param pulumi.Input[str] private_subnet_id1: The ID of the private subnet in Availability Zone 1 in your
               existing VPC (e.g., subnet-a0246dcd).
        :param pulumi.Input[str] private_subnet_id2: The ID of the private subnet in Availability Zone 2 in your
               existing VPC (e.g., subnet-b58c3d67).
        :param pulumi.Input[str] vpc_id: The ID of your existing VPC (e.g., vpc-0343606e) where you
               want to deploy the Aurora database.
        :param bool db_auto_minor_version_upgrade: Set this parameter to true if you want to enable your DB
               instances to receive minor DB engine version upgrades
               automatically when upgrades become available.
        :param int db_backup_retention_period: The number of days to retain automatic database snapshots.
               To disable automatic backups, set this parameter to 0. Default is 35 days
        :param bool db_encrypted_enabled: Set this parameter to false if you don’t want to encrypt the
               database at rest. Defaults to `true`.
        :param int db_num_db_cluster_instances: The number of db instances to launch as part of the cluster. Defaults to 1.
        :param float db_port: The port that you want to access the database through. The DB
               instance will listen on this port for connections. This value
               must be in the range 1115-65535. Default is 5432
        :param pulumi.Input[str] db_security_group_id: The ID of the custom security group you want to use in your
               existing VPC (e.g., sg-7f16e910).
        :param bool enable_event_subscription: Set this parameter to `false` if you want to disable Amazon
               Aurora Cluster and Instance level event subscriptions. You
               might want to disable it if you are testing or running
               continuous integration (CI) processes.
        :param str sns_notification_email: The email that is used to configure an SNS topic for sending
               CloudWatch alarms and Amazon RDS event notifications. This
               must be a valid email address. Required if enableEventSubscription is true.
        """
        pulumi.set(__self__, "availability_zone_names", availability_zone_names)
        pulumi.set(__self__, "db_engine_version", db_engine_version)
        pulumi.set(__self__, "db_instance_class", db_instance_class)
        pulumi.set(__self__, "db_master_password", db_master_password)
        pulumi.set(__self__, "db_master_username", db_master_username)
        pulumi.set(__self__, "db_name", db_name)
        pulumi.set(__self__, "db_parameter_group_family", db_parameter_group_family)
        pulumi.set(__self__, "private_subnet_id1", private_subnet_id1)
        pulumi.set(__self__, "private_subnet_id2", private_subnet_id2)
        pulumi.set(__self__, "vpc_id", vpc_id)
        if db_auto_minor_version_upgrade is not None:
            pulumi.set(__self__, "db_auto_minor_version_upgrade", db_auto_minor_version_upgrade)
        if db_backup_retention_period is not None:
            pulumi.set(__self__, "db_backup_retention_period", db_backup_retention_period)
        if db_encrypted_enabled is not None:
            pulumi.set(__self__, "db_encrypted_enabled", db_encrypted_enabled)
        if db_num_db_cluster_instances is not None:
            pulumi.set(__self__, "db_num_db_cluster_instances", db_num_db_cluster_instances)
        if db_port is not None:
            pulumi.set(__self__, "db_port", db_port)
        if db_security_group_id is not None:
            pulumi.set(__self__, "db_security_group_id", db_security_group_id)
        if enable_event_subscription is not None:
            pulumi.set(__self__, "enable_event_subscription", enable_event_subscription)
        if sns_notification_email is not None:
            pulumi.set(__self__, "sns_notification_email", sns_notification_email)

    @property
    @pulumi.getter(name="availabilityZoneNames")
    def availability_zone_names(self) -> Sequence[pulumi.Input[str]]:
        """
        List of Availability Zone names to use to create the DB Cluster.
        """
        return pulumi.get(self, "availability_zone_names")

    @availability_zone_names.setter
    def availability_zone_names(self, value: Sequence[pulumi.Input[str]]):
        pulumi.set(self, "availability_zone_names", value)

    @property
    @pulumi.getter(name="dbEngineVersion")
    def db_engine_version(self) -> str:
        """
        The version of the database engine.
        """
        return pulumi.get(self, "db_engine_version")

    @db_engine_version.setter
    def db_engine_version(self, value: str):
        pulumi.set(self, "db_engine_version", value)

    @property
    @pulumi.getter(name="dbInstanceClass")
    def db_instance_class(self) -> str:
        """
        The DB (compute and memory capacity) class for the database
        instances.
        """
        return pulumi.get(self, "db_instance_class")

    @db_instance_class.setter
    def db_instance_class(self, value: str):
        pulumi.set(self, "db_instance_class", value)

    @property
    @pulumi.getter(name="dbMasterPassword")
    def db_master_password(self) -> pulumi.Input[str]:
        """
        The password for the database administrator account (8-64
        character string)
        """
        return pulumi.get(self, "db_master_password")

    @db_master_password.setter
    def db_master_password(self, value: pulumi.Input[str]):
        pulumi.set(self, "db_master_password", value)

    @property
    @pulumi.getter(name="dbMasterUsername")
    def db_master_username(self) -> str:
        """
        The user name for the database administrator account. This is
        an alphanumeric string of 1-16 characters. The user name
        must start with an uppercase or lowercase letter (A-Z, a-z).
        """
        return pulumi.get(self, "db_master_username")

    @db_master_username.setter
    def db_master_username(self, value: str):
        pulumi.set(self, "db_master_username", value)

    @property
    @pulumi.getter(name="dbName")
    def db_name(self) -> str:
        """
        The name of the Aurora DB to provision. This is an
        alphanumeric string of 5-64 characters.
        """
        return pulumi.get(self, "db_name")

    @db_name.setter
    def db_name(self, value: str):
        pulumi.set(self, "db_name", value)

    @property
    @pulumi.getter(name="dbParameterGroupFamily")
    def db_parameter_group_family(self) -> str:
        """
        The family of the DB parameter group (e.g. aurora-postgresql11).
        """
        return pulumi.get(self, "db_parameter_group_family")

    @db_parameter_group_family.setter
    def db_parameter_group_family(self, value: str):
        pulumi.set(self, "db_parameter_group_family", value)

    @property
    @pulumi.getter(name="privateSubnetID1")
    def private_subnet_id1(self) -> pulumi.Input[str]:
        """
        The ID of the private subnet in Availability Zone 1 in your
        existing VPC (e.g., subnet-a0246dcd).
        """
        return pulumi.get(self, "private_subnet_id1")

    @private_subnet_id1.setter
    def private_subnet_id1(self, value: pulumi.Input[str]):
        pulumi.set(self, "private_subnet_id1", value)

    @property
    @pulumi.getter(name="privateSubnetID2")
    def private_subnet_id2(self) -> pulumi.Input[str]:
        """
        The ID of the private subnet in Availability Zone 2 in your
        existing VPC (e.g., subnet-b58c3d67).
        """
        return pulumi.get(self, "private_subnet_id2")

    @private_subnet_id2.setter
    def private_subnet_id2(self, value: pulumi.Input[str]):
        pulumi.set(self, "private_subnet_id2", value)

    @property
    @pulumi.getter(name="vpcID")
    def vpc_id(self) -> pulumi.Input[str]:
        """
        The ID of your existing VPC (e.g., vpc-0343606e) where you
        want to deploy the Aurora database.
        """
        return pulumi.get(self, "vpc_id")

    @vpc_id.setter
    def vpc_id(self, value: pulumi.Input[str]):
        pulumi.set(self, "vpc_id", value)

    @property
    @pulumi.getter(name="dbAutoMinorVersionUpgrade")
    def db_auto_minor_version_upgrade(self) -> Optional[bool]:
        """
        Set this parameter to true if you want to enable your DB
        instances to receive minor DB engine version upgrades
        automatically when upgrades become available.
        """
        return pulumi.get(self, "db_auto_minor_version_upgrade")

    @db_auto_minor_version_upgrade.setter
    def db_auto_minor_version_upgrade(self, value: Optional[bool]):
        pulumi.set(self, "db_auto_minor_version_upgrade", value)

    @property
    @pulumi.getter(name="dbBackupRetentionPeriod")
    def db_backup_retention_period(self) -> Optional[int]:
        """
        The number of days to retain automatic database snapshots.
        To disable automatic backups, set this parameter to 0. Default is 35 days
        """
        return pulumi.get(self, "db_backup_retention_period")

    @db_backup_retention_period.setter
    def db_backup_retention_period(self, value: Optional[int]):
        pulumi.set(self, "db_backup_retention_period", value)

    @property
    @pulumi.getter(name="dbEncryptedEnabled")
    def db_encrypted_enabled(self) -> Optional[bool]:
        """
        Set this parameter to false if you don’t want to encrypt the
        database at rest. Defaults to `true`.
        """
        return pulumi.get(self, "db_encrypted_enabled")

    @db_encrypted_enabled.setter
    def db_encrypted_enabled(self, value: Optional[bool]):
        pulumi.set(self, "db_encrypted_enabled", value)

    @property
    @pulumi.getter(name="dbNumDbClusterInstances")
    def db_num_db_cluster_instances(self) -> Optional[int]:
        """
        The number of db instances to launch as part of the cluster. Defaults to 1.
        """
        return pulumi.get(self, "db_num_db_cluster_instances")

    @db_num_db_cluster_instances.setter
    def db_num_db_cluster_instances(self, value: Optional[int]):
        pulumi.set(self, "db_num_db_cluster_instances", value)

    @property
    @pulumi.getter(name="dbPort")
    def db_port(self) -> Optional[float]:
        """
        The port that you want to access the database through. The DB
        instance will listen on this port for connections. This value
        must be in the range 1115-65535. Default is 5432
        """
        return pulumi.get(self, "db_port")

    @db_port.setter
    def db_port(self, value: Optional[float]):
        pulumi.set(self, "db_port", value)

    @property
    @pulumi.getter(name="dbSecurityGroupID")
    def db_security_group_id(self) -> Optional[pulumi.Input[str]]:
        """
        The ID of the custom security group you want to use in your
        existing VPC (e.g., sg-7f16e910).
        """
        return pulumi.get(self, "db_security_group_id")

    @db_security_group_id.setter
    def db_security_group_id(self, value: Optional[pulumi.Input[str]]):
        pulumi.set(self, "db_security_group_id", value)

    @property
    @pulumi.getter(name="enableEventSubscription")
    def enable_event_subscription(self) -> Optional[bool]:
        """
        Set this parameter to `false` if you want to disable Amazon
        Aurora Cluster and Instance level event subscriptions. You
        might want to disable it if you are testing or running
        continuous integration (CI) processes.
        """
        return pulumi.get(self, "enable_event_subscription")

    @enable_event_subscription.setter
    def enable_event_subscription(self, value: Optional[bool]):
        pulumi.set(self, "enable_event_subscription", value)

    @property
    @pulumi.getter(name="snsNotificationEmail")
    def sns_notification_email(self) -> Optional[str]:
        """
        The email that is used to configure an SNS topic for sending
        CloudWatch alarms and Amazon RDS event notifications. This
        must be a valid email address. Required if enableEventSubscription is true.
        """
        return pulumi.get(self, "sns_notification_email")

    @sns_notification_email.setter
    def sns_notification_email(self, value: Optional[str]):
        pulumi.set(self, "sns_notification_email", value)


class Cluster(pulumi.ComponentResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 availability_zone_names: Optional[Sequence[pulumi.Input[str]]] = None,
                 db_auto_minor_version_upgrade: Optional[bool] = None,
                 db_backup_retention_period: Optional[int] = None,
                 db_encrypted_enabled: Optional[bool] = None,
                 db_engine_version: Optional[str] = None,
                 db_instance_class: Optional[str] = None,
                 db_master_password: Optional[pulumi.Input[str]] = None,
                 db_master_username: Optional[str] = None,
                 db_name: Optional[str] = None,
                 db_num_db_cluster_instances: Optional[int] = None,
                 db_parameter_group_family: Optional[str] = None,
                 db_port: Optional[float] = None,
                 db_security_group_id: Optional[pulumi.Input[str]] = None,
                 enable_event_subscription: Optional[bool] = None,
                 private_subnet_id1: Optional[pulumi.Input[str]] = None,
                 private_subnet_id2: Optional[pulumi.Input[str]] = None,
                 sns_notification_email: Optional[str] = None,
                 vpc_id: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        """
        Create a Cluster resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param Sequence[pulumi.Input[str]] availability_zone_names: List of Availability Zone names to use to create the DB Cluster.
        :param bool db_auto_minor_version_upgrade: Set this parameter to true if you want to enable your DB
               instances to receive minor DB engine version upgrades
               automatically when upgrades become available.
        :param int db_backup_retention_period: The number of days to retain automatic database snapshots.
               To disable automatic backups, set this parameter to 0. Default is 35 days
        :param bool db_encrypted_enabled: Set this parameter to false if you don’t want to encrypt the
               database at rest. Defaults to `true`.
        :param str db_engine_version: The version of the database engine.
        :param str db_instance_class: The DB (compute and memory capacity) class for the database
               instances.
        :param pulumi.Input[str] db_master_password: The password for the database administrator account (8-64
               character string)
        :param str db_master_username: The user name for the database administrator account. This is
               an alphanumeric string of 1-16 characters. The user name
               must start with an uppercase or lowercase letter (A-Z, a-z).
        :param str db_name: The name of the Aurora DB to provision. This is an
               alphanumeric string of 5-64 characters.
        :param int db_num_db_cluster_instances: The number of db instances to launch as part of the cluster. Defaults to 1.
        :param str db_parameter_group_family: The family of the DB parameter group (e.g. aurora-postgresql11).
        :param float db_port: The port that you want to access the database through. The DB
               instance will listen on this port for connections. This value
               must be in the range 1115-65535. Default is 5432
        :param pulumi.Input[str] db_security_group_id: The ID of the custom security group you want to use in your
               existing VPC (e.g., sg-7f16e910).
        :param bool enable_event_subscription: Set this parameter to `false` if you want to disable Amazon
               Aurora Cluster and Instance level event subscriptions. You
               might want to disable it if you are testing or running
               continuous integration (CI) processes.
        :param pulumi.Input[str] private_subnet_id1: The ID of the private subnet in Availability Zone 1 in your
               existing VPC (e.g., subnet-a0246dcd).
        :param pulumi.Input[str] private_subnet_id2: The ID of the private subnet in Availability Zone 2 in your
               existing VPC (e.g., subnet-b58c3d67).
        :param str sns_notification_email: The email that is used to configure an SNS topic for sending
               CloudWatch alarms and Amazon RDS event notifications. This
               must be a valid email address. Required if enableEventSubscription is true.
        :param pulumi.Input[str] vpc_id: The ID of your existing VPC (e.g., vpc-0343606e) where you
               want to deploy the Aurora database.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: ClusterArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        Create a Cluster resource with the given unique name, props, and options.
        :param str resource_name: The name of the resource.
        :param ClusterArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(ClusterArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 availability_zone_names: Optional[Sequence[pulumi.Input[str]]] = None,
                 db_auto_minor_version_upgrade: Optional[bool] = None,
                 db_backup_retention_period: Optional[int] = None,
                 db_encrypted_enabled: Optional[bool] = None,
                 db_engine_version: Optional[str] = None,
                 db_instance_class: Optional[str] = None,
                 db_master_password: Optional[pulumi.Input[str]] = None,
                 db_master_username: Optional[str] = None,
                 db_name: Optional[str] = None,
                 db_num_db_cluster_instances: Optional[int] = None,
                 db_parameter_group_family: Optional[str] = None,
                 db_port: Optional[float] = None,
                 db_security_group_id: Optional[pulumi.Input[str]] = None,
                 enable_event_subscription: Optional[bool] = None,
                 private_subnet_id1: Optional[pulumi.Input[str]] = None,
                 private_subnet_id2: Optional[pulumi.Input[str]] = None,
                 sns_notification_email: Optional[str] = None,
                 vpc_id: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        if opts is None:
            opts = pulumi.ResourceOptions()
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.version is None:
            opts.version = _utilities.get_version()
        if opts.id is not None:
            raise ValueError('ComponentResource classes do not support opts.id')
        else:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = ClusterArgs.__new__(ClusterArgs)

            if availability_zone_names is None and not opts.urn:
                raise TypeError("Missing required property 'availability_zone_names'")
            __props__.__dict__["availability_zone_names"] = availability_zone_names
            __props__.__dict__["db_auto_minor_version_upgrade"] = db_auto_minor_version_upgrade
            __props__.__dict__["db_backup_retention_period"] = db_backup_retention_period
            __props__.__dict__["db_encrypted_enabled"] = db_encrypted_enabled
            if db_engine_version is None and not opts.urn:
                raise TypeError("Missing required property 'db_engine_version'")
            __props__.__dict__["db_engine_version"] = db_engine_version
            if db_instance_class is None and not opts.urn:
                raise TypeError("Missing required property 'db_instance_class'")
            __props__.__dict__["db_instance_class"] = db_instance_class
            if db_master_password is None and not opts.urn:
                raise TypeError("Missing required property 'db_master_password'")
            __props__.__dict__["db_master_password"] = None if db_master_password is None else pulumi.Output.secret(db_master_password)
            if db_master_username is None and not opts.urn:
                raise TypeError("Missing required property 'db_master_username'")
            __props__.__dict__["db_master_username"] = db_master_username
            if db_name is None and not opts.urn:
                raise TypeError("Missing required property 'db_name'")
            __props__.__dict__["db_name"] = db_name
            __props__.__dict__["db_num_db_cluster_instances"] = db_num_db_cluster_instances
            if db_parameter_group_family is None and not opts.urn:
                raise TypeError("Missing required property 'db_parameter_group_family'")
            __props__.__dict__["db_parameter_group_family"] = db_parameter_group_family
            __props__.__dict__["db_port"] = db_port
            __props__.__dict__["db_security_group_id"] = db_security_group_id
            __props__.__dict__["enable_event_subscription"] = enable_event_subscription
            if private_subnet_id1 is None and not opts.urn:
                raise TypeError("Missing required property 'private_subnet_id1'")
            __props__.__dict__["private_subnet_id1"] = private_subnet_id1
            if private_subnet_id2 is None and not opts.urn:
                raise TypeError("Missing required property 'private_subnet_id2'")
            __props__.__dict__["private_subnet_id2"] = private_subnet_id2
            __props__.__dict__["sns_notification_email"] = sns_notification_email
            if vpc_id is None and not opts.urn:
                raise TypeError("Missing required property 'vpc_id'")
            __props__.__dict__["vpc_id"] = vpc_id
        super(Cluster, __self__).__init__(
            'aws-quickstart-aurora-postgres:index:Cluster',
            resource_name,
            __props__,
            opts,
            remote=True)

