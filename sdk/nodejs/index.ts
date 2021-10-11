// *** WARNING: this file was generated by Pulumi SDK Generator. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "./utilities";

// Export members:
export * from "./cluster";
export * from "./provider";

// Import resources to register:
import { Cluster } from "./cluster";

const _module = {
    version: utilities.getVersion(),
    construct: (name: string, type: string, urn: string): pulumi.Resource => {
        switch (type) {
            case "aws-quickstart-aurora-postgres:index:Cluster":
                return new Cluster(name, <any>undefined, { urn })
            default:
                throw new Error(`unknown resource type ${type}`);
        }
    },
};
pulumi.runtime.registerResourceModule("aws-quickstart-aurora-postgres", "index", _module)

import { Provider } from "./provider";

pulumi.runtime.registerResourcePackage("aws-quickstart-aurora-postgres", {
    version: utilities.getVersion(),
    constructProvider: (name: string, type: string, urn: string): pulumi.ProviderResource => {
        if (type !== "pulumi:providers:aws-quickstart-aurora-postgres") {
            throw new Error(`unknown provider type ${type}`);
        }
        return new Provider(name, <any>undefined, { urn });
    },
});