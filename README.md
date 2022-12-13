# wao-estimator

[![GitHub](https://img.shields.io/github/license/Nedopro2022/wao-estimator)](https://github.com/Nedopro2022/wao-estimator/blob/main/LICENSE)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/Nedopro2022/wao-estimator)](https://github.com/Nedopro2022/wao-estimator/releases/latest)
[![CI](https://github.com/Nedopro2022/wao-estimator/actions/workflows/ci.yaml/badge.svg)](https://github.com/Nedopro2022/wao-estimator/actions/workflows/ci.yaml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Nedopro2022/wao-estimator)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nedopro2022/wao-estimator)](https://goreportcard.com/report/github.com/Nedopro2022/wao-estimator)
[![codecov](https://codecov.io/gh/Nedopro2022/wao-estimator/branch/main/graph/badge.svg)](https://codecov.io/gh/Nedopro2022/wao-estimator)


// TODO(user): Add simple overview of use/purpose

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Developing](#developing)
  - [Prerequisites](#prerequisites)
  - [Run a development cluster with kind](#run-a-development-cluster-with-kind)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Overview
// TODO(user): An in-depth paragraph about your project and overview of use

## Getting Started
// TODO

## Developing

This Operator uses [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder), so we basically follow the Kubebuilder way. See the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html) for details.

### Prerequisites

Make sure you have the following tools installed:

- Git
- Make
- Go
- Docker


### Run a development cluster with [kind](https://kind.sigs.k8s.io/)

```sh
./hack/dev-kind-reset-cluster.sh # create a K8s cluster `kind-wao-estimator`
./hack/dev-kind-deploy.sh # build and deploy the Operator
```

## License

Copyright 2022 Bitmedia Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
