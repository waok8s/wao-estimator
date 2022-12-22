# wao-estimator

[![GitHub](https://img.shields.io/github/license/Nedopro2022/wao-estimator)](https://github.com/Nedopro2022/wao-estimator/blob/main/LICENSE)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/Nedopro2022/wao-estimator)](https://github.com/Nedopro2022/wao-estimator/releases/latest)
[![CI](https://github.com/Nedopro2022/wao-estimator/actions/workflows/ci.yaml/badge.svg)](https://github.com/Nedopro2022/wao-estimator/actions/workflows/ci.yaml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Nedopro2022/wao-estimator)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nedopro2022/wao-estimator)](https://goreportcard.com/report/github.com/Nedopro2022/wao-estimator)
[![codecov](https://codecov.io/gh/Nedopro2022/wao-estimator/branch/main/graph/badge.svg)](https://codecov.io/gh/Nedopro2022/wao-estimator)

WAO-Estimator provides power consumption estimation capabilities to help schedule Pods in power efficient way.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Overview](#overview)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Setup WAO-Estimator by deploying an Estimator resource.](#setup-wao-estimator-by-deploying-an-estimator-resource)
  - [Check operation with estimator-cli](#check-operation-with-estimator-cli)
  - [Detailed configuration of Estimator resource](#detailed-configuration-of-estimator-resource)
  - [Uninstallation](#uninstallation)
  - [💡 Demo using kind and FakeNodeMonitor / FakePCPredictor](#-demo-using-kind-and-fakenodemonitor--fakepcpredictor)
- [Technical Details](#technical-details)
  - [Prediction models](#prediction-models)
  - [Estimation algorithms](#estimation-algorithms)
  - [NodeMonitor implementations](#nodemonitor-implementations)
  - [PowerConsumptionPredictor implementations](#powerconsumptionpredictor-implementations)
  - [HTTP APIs](#http-apis)
- [Developing](#developing)
  - [Prerequisites](#prerequisites)
  - [Run a development cluster with kind](#run-a-development-cluster-with-kind)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Overview

WAO-Estimator provides APIs for estimating the increase in power consumption when deploying Pods in a cluster.

Specifically, it responds to queries like this:

**Q.** *Tell me the increase in power consumption when deploying 5 Pods each requiring 500 mCPU.*

```json
{ "cpu_milli": 500, "num_workloads": 5 }
```

**A.** *The power consumption of the entire cluster will **at least** increase by 5W with 1 Pod, 9W with 2 Pods, 13W with 3 Pods, 15W with 4 Pods, and 16W with 5 Pods.*

```json
{ "watt_increases": [ 5, 9, 13, 15, 16 ] }
```

Each value indicates the following information:

- `cpu_milli`: the amount of CPU consumed by a single workload
- `num_workloads`: the number of workloads
- `watt_increases`: the estimated increase in power consumption when the workloads are placed

For use cases, see [WAO-Scheduler-v2](https://github.com/Nedopro2022/wao-scheduler-v2), which uses WAO-Estimator to place pods on the cluster, and [WAOFed](https://github.com/Nedopro2022/waofed), which works with KubeFed to optimally place Pods in multi-cluster environments.


## Getting Started

To start using WAO-Estimator, you need a Kubernetes cluster that meets the following conditions:

- Each worker node is a **physical machine**
- Each worker node supports IPMI or [Redfish](https://www.dmtf.org/standards/redfish)
- Power consumption models for each type of physical machine
- Environmental information required by the model are obtainable
  - e.g. static pressure difference between front and back of a server

Supported Kubernetes versions: __1.19 or higher__

> 💡 Mainly tested with 1.25, may work with old versions (but may require some efforts).

### Installation

Make sure you have [cert-manager](https://cert-manager.io/) deployed on the cluster where KubeFed control plane is deployed, as it is used to generate webhook certificates.

```sh
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.0/cert-manager.yaml
```

> ⚠️ You may have to wait a second for cert-manager to be ready.

Deploy the Operator with the following command. It creates `wao-estimator-system` namespace and deploys CRDs, controllers and other resources.

```sh
kubectl apply -f https://github.com/Nedopro2022/wao-estimator/releases/download/v0.1.0/wao-estimator.yaml
```

> 💡 Please verify that 3 Service objects have been created. **webhook-service** and **metrics-service** are normal and **estimator-service** is for providing WAO-Estimator APIs.
> 
> ```
> $ kubectl get svc -n wao-estimator-system
> NAME                                                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
> wao-estimator-controller-manager-estimator-service   ClusterIP   10.96.222.246   <none>        5656/TCP   5m
> wao-estimator-controller-manager-metrics-service     ClusterIP   10.96.161.177   <none>        8443/TCP   5m
> wao-estimator-webhook-service                        ClusterIP   10.96.217.214   <none>        443/TCP    5m
> ```

### Setup WAO-Estimator by deploying an Estimator resource.

This is a simple example of an Estimator resource.

```yaml
apiVersion: waofed.bitmedia.co.jp/v1beta1
kind: Estimator
metadata:
  namespace: default
  name: default
spec:
  nodeMonitor:
    type:
      default: None
  powerConsumptionPredictor:
    type:
      default: None
```

The namespace and name will affect the URL of the API endpoint.

```
http://<yourhost>:5656/namespaces/<namespace>/estimators/<name>/values/powerconsumption
```

In this case, we specified `namespace: default` `name: default`, so it will be:

```
http://<yourhost>:5656/namespaces/default/estimators/default/values/powerconsumption
```

> 💡 This allows for running multiple Estimators. Using `default/default` is usually fine.

For simplicity, we specified `spec.nodeMonitor` and `spec.powerConsumptionPredictor` as `None`. They will not function correctly, but it will still allow us to verify the behavior of the API.

### Check operation with estimator-cli

`estimator-cli` is a CLI tool for using WAO-Estimator APIs. It is essentially an HTTP client.

First, execute the following command to allow local access to the **estimator-service** Service resource.

> 💡 You can also use LoadBalancer or Ingress in supported environments.

```sh
kubectl port-forward -n wao-estimator-system svc/wao-estimator-controller-manager-estimator-service 5656:5656
```

Next, run `estimator-cli` in another terminal.

```
$ go run github.com/Nedopro2022/wao-estimator/pkg/cmd/estimator-cli -p 500,5 pc
[+Inf +Inf +Inf +Inf +Inf]
```

Since `spec.nodeMonitor` and `spec.powerConsumptionPredictor` are specified as `None`, the response (increases in power consumption) is not correct, but we can confirm that WAO-Estimator APIs are working.

> 💡 You can see the actual HTTP request by adding `-v` option to `estimator-cli`. For example, you will see the following request by running `./estimator-cli -v -p 500,5 pc`. See the help `-h` for details.
>
> ```
> curl -X 'POST' -d '{"cpu_milli":500,"num_workloads":5}' -H 'Content-Type: application/json' 'http://localhost:5656/namespaces/default/estimators/default/values/powerconsumption'
> ```

### Detailed configuration of Estimator resource

// TODO

### Uninstallation

Delete the Operator and resources with the following command.

```sh
kubectl delete -f https://github.com/Nedopro2022/wao-estimator/releases/download/v0.1.0/wao-estimator.yaml
```

### 💡 Demo using kind and FakeNodeMonitor / FakePCPredictor

// TODO

## Technical Details

// TODO

### Prediction models

### Estimation algorithms

### NodeMonitor implementations

### PowerConsumptionPredictor implementations

### HTTP APIs

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
