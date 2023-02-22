package controllers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
)

const (
	ctxK8sClient  = "ctxK8sClient"
	ctxNodeObjKey = "ctxNodeObjKey"

	labelNodeStatusCPUUsage           = "waofed.bitmedia.co.jp/node-status.cpuusage"
	labelNodeStatusAmbientTemp        = "waofed.bitmedia.co.jp/node-status.ambienttemp"
	labelNodeStatusStaticPressureDiff = "waofed.bitmedia.co.jp/node-status.staticpressurediff"

	labelFakePCPredictorBaseWatts    = "waofed.bitmedia.co.jp/fakepcp.basewatts"
	labelFakePCPredictorWattsPerCore = "waofed.bitmedia.co.jp/fakepcp.wpc"
)

func setupFakeNodeMonitor(k8sClient client.Client, nodeObjKey client.ObjectKey) estimator.NodeMonitor {
	fn := func(ctx context.Context, base *estimator.NodeStatus) (*estimator.NodeStatus, error) {

		if base == nil {
			base = estimator.NewNodeStatus()
		}

		var node corev1.Node
		if err := k8sClient.Get(ctx, nodeObjKey, &node); err != nil {
			return nil, err
		}

		cu, err := labelValueFloat(node.Labels, labelNodeStatusCPUUsage)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return nil, err
		}
		estimator.NodeStatusSetCPUUsage(base, cu)

		at, err := labelValueFloat(node.Labels, labelNodeStatusAmbientTemp)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return nil, err
		}
		estimator.NodeStatusSetAmbientTemp(base, at)

		spd, err := labelValueFloat(node.Labels, labelNodeStatusStaticPressureDiff)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return nil, err
		}
		estimator.NodeStatusSetStaticPressureDiff(base, spd)

		return base, nil
	}

	nm := &estimator.FakeNodeMonitor{FetchFunc: fn}
	return nm
}

func setupFakePCPredictor(k8sClient client.Client, nodeObjKey client.ObjectKey) estimator.PowerConsumptionPredictor {
	fn := func(ctx context.Context, requestCPUMilli int, status *estimator.NodeStatus) (watt float64, err error) {

		var node corev1.Node
		if err := k8sClient.Get(ctx, nodeObjKey, &node); err != nil {
			return 0.0, err
		}

		bw, err := labelValueInt(node.Labels, labelFakePCPredictorBaseWatts)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return 0.0, err
		}

		wpc, err := labelValueInt(node.Labels, labelFakePCPredictorWattsPerCore)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return 0.0, err
		}

		time.Sleep(10 * time.Millisecond) // emulate the time to respond

		return float64(bw) + ((float64(requestCPUMilli) / 1000) * float64(wpc)), nil
	}

	pcp := &estimator.FakePCPredictor{PredictFunc: fn}
	return pcp
}

var errLabelNotFound = errors.New("errLabelNotFound")

func labelValueString(m map[string]string, label string) (string, error) {
	if label == "" {
		return "", fmt.Errorf("label could not be empty")
	}
	v, ok := m[label]
	if !ok {
		return "", fmt.Errorf("unable to get value for label=%v (%w)", label, errLabelNotFound)
	}
	return v, nil
}

func labelValueInt(m map[string]string, label string) (int, error) {
	v, err := labelValueString(m, label)
	if err != nil {
		return 0, fmt.Errorf("unable to get value for label=%v (%w)", label, errLabelNotFound)
	}
	vv, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("unable to convert value for label=%v err=%w", label, err)
	}
	return vv, nil
}

func labelValueFloat(m map[string]string, label string) (float64, error) {
	v, err := labelValueString(m, label)
	if err != nil {
		return 0, fmt.Errorf("unable to get value for label=%v (%w)", label, errLabelNotFound)
	}
	vv, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to convert value for label=%v err=%w", label, err)
	}
	return vv, nil
}
