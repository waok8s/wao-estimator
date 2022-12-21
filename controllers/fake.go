package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Nedopro2022/wao-estimator/pkg/estimator"
)

const (
	ctxK8sClient  = "ctxK8sClient"
	ctxNodeObjKey = "ctxNodeObjKey"

	labelNodeStatusCPUSockets     = "waofed.bitmedia.co.jp/node-status.cpusockets"
	labelNodeStatusCPUCores       = "waofed.bitmedia.co.jp/node-status.cpucores"
	labelNodeStatusCPUUsages      = "waofed.bitmedia.co.jp/node-status.cpuusages"
	labelNodeStatusCPUTemps       = "waofed.bitmedia.co.jp/node-status.cputemps"
	labelNodeStatusAmbientSensors = "waofed.bitmedia.co.jp/node-status.ambientsensors"
	labelNodeStatusAmbientTemps   = "waofed.bitmedia.co.jp/node-status.ambienttemps"

	labelFakePCPredictorBaseWatts    = "waofed.bitmedia.co.jp/fakepcp.basewatts"
	labelFakePCPredictorWattsPerCore = "waofed.bitmedia.co.jp/fakepcp.wpc"
)

func setupFakeNodeMonitor(k8sClient client.Client, nodeObjKey client.ObjectKey) estimator.NodeMonitor {
	fn := func(ctx context.Context) (estimator.NodeStatus, error) {

		var node corev1.Node
		if err := k8sClient.Get(ctx, nodeObjKey, &node); err != nil {
			return estimator.NodeStatus{}, err
		}

		s := estimator.NodeStatus{
			Timestamp: time.Now(),
		}

		cs, err := labelValueInt(node.Labels, labelNodeStatusCPUSockets)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return estimator.NodeStatus{}, err
		}
		s.CPUSockets = cs

		cc, err := labelValueInt(node.Labels, labelNodeStatusCPUCores)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return estimator.NodeStatus{}, err
		}
		s.CPUCores = cc

		cu, err := labelValueJSON[[][]float64](node.Labels, labelNodeStatusCPUUsages)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return estimator.NodeStatus{}, err
		}
		s.CPUUsages = cu

		ct, err := labelValueJSON[[][]float64](node.Labels, labelNodeStatusCPUTemps)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return estimator.NodeStatus{}, err
		}
		s.CPUTemps = ct

		as, err := labelValueInt(node.Labels, labelNodeStatusAmbientSensors)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return estimator.NodeStatus{}, err
		}
		s.AmbientSensors = as

		at, err := labelValueJSON[[]float64](node.Labels, labelNodeStatusAmbientTemps)
		if err != nil && !errors.Is(err, errLabelNotFound) {
			return estimator.NodeStatus{}, err
		}
		s.AmbientTemps = at

		return s, nil
	}

	nm := &estimator.FakeNodeMonitor{FetchFunc: fn}
	return nm
}

func setupFakePCPredictor(k8sClient client.Client, nodeObjKey client.ObjectKey) estimator.PowerConsumptionPredictor {
	fn := func(ctx context.Context, requestCPUMilli int, status estimator.NodeStatus) (watt float64, err error) {

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

		time.Sleep(10 * time.Millisecond) // emulate response time

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

// labelValueJSON loads a label value, convert it to JSON and then decode it.
//
// Conversion:
//   - "_" -> ""
//   - "b" -> "["
//   - "d" -> "]"
//   - "p" -> ","
//   - e.g. "b__b50.0p30.0d__p__b50p30d__d" -> "[[50.0,30.0],[50,30]]"
func labelValueJSON[T []float64 | [][]float64](m map[string]string, label string) (T, error) {
	v, err := labelValueString(m, label)
	if err != nil {
		return nil, fmt.Errorf("unable to get value for label=%v (%w)", label, errLabelNotFound)
	}
	v = strings.ReplaceAll(v, "_", "")
	v = strings.ReplaceAll(v, "b", "[")
	v = strings.ReplaceAll(v, "d", "]")
	v = strings.ReplaceAll(v, "p", ",")
	var vv T
	if err := json.Unmarshal([]byte(v), &vv); err != nil {
		return nil, fmt.Errorf("unable to convert value for label=%v err=%w", label, err)
	}
	return vv, nil
}
