# Camunda Load Test Helm 

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Test - Unit](https://github.com/camunda/camunda-load-tests-helm/actions/workflows/build.yaml/badge.svg)](https://github.com/camunda/camunda-load-tests-helm/actions/workflows/build.yaml)

- [Overview](#overview)
- [Installation](#installation)
- [Contributing](#contributing)
- [Versioning](#versioning)
- [Releasing](#releasing)
- [License](#license)

## Overview

Camunda load test Helm charts repo. The load test Helm chart is a chart which allows the Camunda engineering team
to run reliability tests to cover stability and performance. 

It can be executed against any Camunda Platform 8 cluster, either self-managed or the Camunda SaaS offering.

## Installation

Find out more details about different installation and deployment options
on the [Camunda load test Helm chart README](./charts/camunda-load-test/README.md).

## Contributing

We value all feedback and contributions. To start contributing to this project, please:

- **Don't create a PR without opening [an issue](https://github.com/camunda/camunda-load-tests-helm/issues/new/choose)
  and/or discussing it first.**
- Familiarize yourself with the
  [contribution guide](./CONTRIBUTING.md).
- Find more information about configuring and deploying the Camunda load tests
  [Helm chart](charts/camunda-load-test/README.md).

## Versioning

Right now is the versioning not aligned with the general [Camunda Platform 8](https://github.com/camunda/camunda-platform).
Since it is mostly for internal use we will not release a 1.x version soon.

## Releasing

Please visit the [Camunda Load Tests release guide](./RELEASE.md) to find out how to release the charts.

## License

Camunda load tests Helm charts are licensed under the open-source Apache License 2.0.
Please see [LICENSE](LICENSE) for details.

For Camunda Platform 8 components, please visit
[licensing information page](https://docs.camunda.io/docs/reference/licenses).
