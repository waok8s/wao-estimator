#!/usr/bin/env bash

# scripts must be run from project root
. hack/2-lib.sh || exit 1

# consts

KIND_IMAGE=${KIND_IMAGE:-"kindest/node:v1.25.3@sha256:f52781bc0d7a19fb6c405c2af83abfeb311f130707a0e219175677e366cc45d1"}

# main

cluster=$PROJECT_NAME-test

lib::start-docker

lib::create-cluster "$cluster" "$KIND_IMAGE"

"$KUBECTL" get nodes

"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.cpuusage"="50.0"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.ambienttemp"="30.0"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.staticpressurediff"="10.0"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/fakepcp.basewatts"="100"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/fakepcp.wpc"="10"

"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.cpuusage"="10.0"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.ambienttemp"="20.0"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.staticpressurediff"="10.0"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/fakepcp.basewatts"="50"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/fakepcp.wpc"="20"
