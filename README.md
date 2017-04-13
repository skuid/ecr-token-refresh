# ecr-token-refresh

`ecr-token-refresh` is a utility for refreshing access tokens to an AWS ECR Registry on a regular interval. It is designed to be used as a sidecar for Spinnaker's Clouddriver service. It's responsible for refreshing the tokens and writing their values to a file. Ideally, these files would be written to a volume shared between Clouddriver and `ecr-token-refresh`.


## Configuration

`ecr-token-refresh` is configured via a configuration file. The default path for this file is at `/opt/config/ecr-token-refresh/config.yaml`. See the sample configuration below.

```yaml
interval: 30m # defines refresh interval
registries: # list of registries to refresh
    - registryId: "12345667"
      region: "us-west-2"
      passwordFile: "/opt/passwords/my-registry.pass"
    - registryId: "0987654321"
      region: "eu-central-1"
      passwordFile: "/opt/passwords/my-registry-in-eu.pass"
```

## Usage

You can use `ecr-token-refresh` with Docker or as a standalone binary. In either case, either drop the config file in the default path or pass the `--config` flag with the path to the configuration when starting the process.


See the `examples` directory for sample Clouddriver ReplicaSets, ConfigMaps and `clouddriver-local` config!