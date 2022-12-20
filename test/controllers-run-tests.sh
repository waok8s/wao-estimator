#!/usr/bin/env bash

# scripts must be run from project root
. hack/0-env.sh || exit 1

# main

cluster=$PROJECT_NAME-test

"$KUBECTL" config use-context kind-"$cluster"

set +x
make test
KUBEBUILDER_ASSETS="$LOCALBIN"/k8s/1.25.0-linux-amd64 go test ./... -coverprofile cover.out -tags=testOnExistingCluster
