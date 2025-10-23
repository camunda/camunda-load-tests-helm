# Camunda Load Test Helm Chart


- [Requirements](#requirements)
- [Dependencies](#dependencies)
- [Installation](#installation)
    - [Local Kubernetes](#local-kubernetes)
- [Uninstalling Charts](#uninstalling-charts)
- [Configuration](#configuration)
    - [Global](#global)
    - [Camunda Platform](#camunda-platform)
    - [Retention Policy](#retention-policy)
    - [Worker](#worker)
    - [Starter](#publisher)
    - [Publisher](#starter)
    - [Timer](#timer)
    - [Leader Balancing](#leader-balancing)
    - [Zeebe](#zeebe)
- [Development](#development)
- [Releasing the Charts](#releasing-the-charts)


## Requirements

* [Helm](https://helm.sh/) >= 3.11.x
* Kubernetes >= 1.20+
* Minimum cluster requirements include the following to run this chart with default settings.
  All of these settings are configurable.
    * Three Kubernetes nodes to respect the default "hard" affinity settings
    * 2GB of RAM for the JVM heap

## Dependencies

The Camunda Load test Helm chart has no direct dependencies. It can ran against any Camunda 8 cluster, self-managed or against the Camunda SaaS offering.

## Installation

The first command adds the Camunda Load Test Helm charts repo, and the second installs the chart to your current 
Kubernetes context.

```shell
  helm repo add camunda-load-tests https://camunda.github.io/camunda-load-tests-helm/
  helm install this-is-a-load-test camunda-load-tests/camunda-load-tests
```

### Local Kubernetes

We recommend using Helm on KIND for local environments, as the Helm configurations are battle-tested
and much closer to production systems.

For more details, follow the Camunda Platform 8
[local Kubernetes cluster guide](https://docs.camunda.io/docs/self-managed/platform-deployment/helm-kubernetes/guides/local-kubernetes-cluster/).

### Running against Saas

If you want to run the load test against a Camunda Saas cluster, or other 
Camunda kubernetes deployment (such as in Google cloud, AWS etc.), you 
need to run through the following steps:

1. Create a new API client in the intended cluster.
2. Download environment vars Under "Env Vars" tab
3. [Source](https://linuxcommand.org/lc3_man_pages/sourceh.html) the downloaded file, containing the environment variables (e.g. `source Credentials.yaml`)
4. Run the following `install` command (please replace `PREFIX` with your initials)

```yaml
helm install PREFIX-saas-load-test camunda-load-tests/camunda-load-tests \
  --set camunda-platform.enabled=false \
  --set saas.enabled=true \
  --set saas.credentials.clientId="$ZEEBE_CLIENT_ID" \
  --set saas.credentials.clientSecret="$ZEEBE_CLIENT_SECRET" \
  --set saas.credentials.zeebeRestAddress="$ZEEBE_REST_ADDRESS" \
  --set saas.credentials.authServer="$ZEEBE_AUTHORIZATION_SERVER_URL" \
  --set saas.credentials.authType="OAUTH" \
  --set saas.credentials.zeebeGrpcAddress="$ZEEBE_GRPC_ADDRESS" \
  --set saas.credentials.authorizationAudience="$ZEEBE_TOKEN_AUDIENCE"
```

This will run the load test against the SaaS Camunda cluster.

## Uninstalling Charts

You can remove these charts by running:

```sh
helm uninstall YOUR_RELEASE_NAME
```

## Configuration

The following sections contain the configuration values for the chart and each sub-chart. All of them can be overwritten
via a separate `values.yaml` file.

Check out the default [values.yaml](values.yaml) file, which contains the same content and documentation.

> **Note**
> For more details about deploying Camunda Platform 8 on Kubernetes, please visit the
> [Helm/Kubernetes installation instructions docs](https://docs.camunda.io/docs/self-managed/platform-deployment/helm-kubernetes/overview/).

### Global

| Section | Parameter | Description | Default |
|-|-|-|-|
| `global` | | Global variables which can be accessed by all sub charts | |
| | `image.repository` | Defines the repository from which to fetch the images. | "gcr.io/zeebe-io" |
| | `image.tag` | Defines the tag / version which should be used in the chart. | "SNAPSHOT" |
| | `image.pullPolicy` | Defines the [image pull policy](https://kubernetes.io/docs/concepts/containers/images/#image-pull-policy) which should be used. | Always |


### Worker

Allows to configure the workers which can be deployed along the Zeebe Cluster. The worker code can be found [here](https://github.com/camunda/camunda/blob/main/load-tests/load-tester/src/main/java/io/camunda/zeebe/Worker.java).

| Section | Parameter | Description | Default |
|-|-|-|-|
| `worker` | | Configuration for the to be deployed worker application | |
| | `replicas` | Defines how many replicas of the application should be deployed. | `3` |
| | `capacity` | Defines how many jobs the worker should activate and work on. | `60` |


### Starter

Allows to configure the starter which can be deployed along the Zeebe Cluster. The start code can be found [here](https://github.com/camunda/camunda/blob/main/load-tests/load-tester/src/main/java/io/camunda/zeebe/Starter.java).

| Section | Parameter | Description | Default |
|-|-|-|-|
| `starter` | | Configuration for the to be deployed starter application | |
| | `replicas` | Defines how many replicas of the application should be deployed. | `1` |
| | `rate` | Defines with which rate process instances should be created by the starter. | `150` |

## Development

For development purposes, you might want to deploy and test the charts without creating a new helm chart release.
To do this you can run the following:

```sh
 helm install YOUR_RELEASE_NAME --atomic --debug ./charts/camunda-platform
```

* `--atomic if set, the installation process deletes the installation on failure. The --wait flag will be set automatically if --atomic is used`

* `--debug enable verbose output`

To generate the resources/manifests without really installing them, you can use:

* `--dry-run simulate an install`

If you see errors like:

```sh
Error: found in Chart.yaml, but missing in charts/ directory: elasticsearch
```

Then you need to download the dependencies first.

Run the following to add resolve the dependencies:

```sh
make helm.repos-add
```

After this, you can run: `make helm.dependency-update`, which will update and download the dependencies for all charts.

The execution should look like this:
```
$ make helm.dependency-update
helm dependency update charts/camunda-platform
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "camunda-platform" chart repository
...Successfully got an update from the "elastic" chart repository
...Successfully got an update from the "bitnami" chart repository
Update Complete. ⎈Happy Helming!⎈
Saving 6 charts
Dependency zeebe did not declare a repository. Assuming it exists in the charts directory
Dependency zeebe-gateway did not declare a repository. Assuming it exists in the charts directory
Dependency operate did not declare a repository. Assuming it exists in the charts directory
Dependency tasklist did not declare a repository. Assuming it exists in the charts directory
Dependency identity did not declare a repository. Assuming it exists in the charts directory
Downloading elasticsearch from repo https://helm.elastic.co
Deleting outdated charts
helm dependency update charts/camunda-platform/charts/identity
Hang tight while we grab the latest from your chart repositories...
...Successfully got an update from the "camunda-platform" chart repository
...Successfully got an update from the "elastic" chart repository
...Successfully got an update from the "bitnami" chart repository
Update Complete. ⎈Happy Helming!⎈
Saving 2 charts
Downloading keycloak from repo https://charts.bitnami.com/bitnami
Downloading common from repo https://charts.bitnami.com/bitnami
```

## Releasing the Charts

Please see the corresponding [release guide](../../RELEASE.md) to find out how to release the chart.
