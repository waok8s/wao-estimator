#!/usr/bin/env bash

# scripts must be run from project root
. hack/2-lib.sh || exit 1

# consts

KIND_IMAGE=${KIND_IMAGE:-"kindest/node:v1.25.3@sha256:f52781bc0d7a19fb6c405c2af83abfeb311f130707a0e219175677e366cc45d1"}

# main

cluster=$PROJECT_NAME-test

lib::start-docker

lib::create-cluster "$cluster" "$KIND_IMAGE"

sleep 15

"$KUBECTL" get nodes


"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-monitor"="Fake"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.cpusockets"="2"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.cpucores"="4"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.cpuusages"="b_b50p50p50p50d_p_b50p50p50p50d_d" # [[50,50,50,50],[50,50,50,50]]
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.cputemps"="b_b50p50p50p50d_p_b50p50p50p50d_d" # [[50,50,50,50],[50,50,50,50]]
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.ambientsensors"="2"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/node-status.ambienttemps"="b30p30d" # [30,30]
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/power-consumption-predictor"="Fake"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/fakepcp.basewatts"="100"
"$KUBECTL" label nodes "$cluster"-worker  "waofed.bitmedia.co.jp/fakepcp.wpc"="10"

"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-monitor"="Fake"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.cpusockets"="2"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.cpucores"="4"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.cpuusages"="b_b10p10p10p10d_p_b10p10p10p10d_d" # [[10,10,10,10],[10,10,10,10]]
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.cputemps"="b_b40p40p40p40d_p_b40p40p40p40d_d" # [[40,40,40,40],[40,40,40,40]]
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.ambientsensors"="2"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/node-status.ambienttemps"="b20p20d" # [20,20]
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/power-consumption-predictor"="Fake"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/fakepcp.basewatts"="50"
"$KUBECTL" label nodes "$cluster"-worker2 "waofed.bitmedia.co.jp/fakepcp.wpc"="20"
