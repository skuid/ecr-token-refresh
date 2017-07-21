# ecr-token-refresh

`ecr-token-refresh` is a utility for refreshing access tokens to an AWS ECR Registry on a regular interval. It is designed to be used as a sidecar for Spinnaker's Clouddriver service. It's responsible for refreshing the tokens and writing their values to a file. Ideally, these files would be written to a volume shared between Clouddriver and `ecr-token-refresh`.


### A special note to Halyard users!

If you wish to use ECR with your Halyard install of Spinnaker, please read this. If you're deploying Spinnaker on a VM, you can take advantage of `ecr-token-refresh`. All you need to do is install/configure it on the same VM that Spinnaker is running on and configure your ECR Docker registries like so: 

```
$ hal config provider docker-registry account add my-docker-registry \
    --address https://<aws-account-number>.dkr.ecr.<region>.aws.com\
    --repositories $REPOSITORIES \
    --username AWS \
    --password-file /opt/passwords/my-docker-registry.pass 
```

Note, the password file _is not_ a file that will be read by Halyard. It will be read by Spinnaker each time it interacts with your registry.

For those wanting to deploy Spinnaker to Kubernetes, there is an [open issue](https://github.com/spinnaker/halyard/issues/597) to make `ecr-token-refresh` part of the Spinnaker deployment. If you'd like to see this happen more quickly, give it a vote or file a PR! 



## Docker Image
`quay.io/skuid/ecr-token-refresh:latest`

## Configuration

`ecr-token-refresh` is configured via a configuration file. The default path for this file is at `/opt/config/ecr-token-refresh/config.yaml`. See the sample configuration below.

```yaml
interval: 30m # defines refresh interval
registries: # list of registries to refresh
    - registryId: "12345667" # typically AWS account ID
      region: "us-west-2" 
      passwordFile: "/opt/passwords/my-registry.pass"
    - registryId: "0987654321"
      region: "eu-central-1"
      passwordFile: "/opt/passwords/my-registry-in-eu.pass"
```

## Usage

You can use `ecr-token-refresh` with Docker or as a standalone binary. In either case, either drop the config file in the default path or pass the `--config` flag with the path to the configuration when starting the process.

See the `examples` directory for sample Clouddriver ReplicaSets, ConfigMaps and `clouddriver-local` config!

### Standalone

If you've installed Spinnaker on a VM, and you don't want to run `ecr-token-refresh` inside of a container, you can download the binary [here](https://github.com/skuid/ecr-token-refresh/releases/tag/1.0.0). Once installed on the same machine as Spinnaker, you can run it like this:

```
$ ./ecr-token-refresh --config {path-to-config}
```

### Docker

To run `ecr-token-refresh` in a Docker container (but _not_ on Kubernetes), you must mount your configuration file into the container as well as a directory for the application to write password files to. You can accomplish this by running the following command:

```
$ docker run -d \
	-v /opt/ecr/passwords:/opt/passwords \
	-v $(pwd)/config.yaml:/opt/config/ecr-token-refresh/config.yaml \
	quay.io/skuid/ecr-token-refresh:latest
```

This will mount your config file at the default config location as well as a directory to write the passwords to (based on the example config file).


## ECR Repository Configuration

In order to use `ecr-token-refresh` with ECR you will need to configure your ECR repository to allow access via IAM. The only permission that `ecr-token-refresh` needs to function is `ecr:GetAuthorizationToken`. Under the Permissions tab, create a new Policy giving the appropriate user or role permission. It's your choice how you want to assign these on the machine. We rely on the machine role that `ecr-token-refresh` is running on having access to the repository. You can also use Access and Secret Keys if you are so inclined.