package v1beta1

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultNodeConf = &NodeConfig{
		NodeMonitor: &NodeMonitor{
			RefreshInterval: &metav1.Duration{
				Duration: time.Minute,
			},
			Agents: []NodeMonitorAgent{{
				Type:     NodeMonitorTypeFake,
				Endpoint: "foo",
			}, {
				Type:     NodeMonitorTypeIPMIExporter,
				Endpoint: "bar",
			}},
		},
		PowerConsumptionPredictor: &PowerConsumptionPredictor{
			Type:     PowerConsumptionPredictorTypeFake,
			Endpoint: "baz",
		},
	}
	node3NodeConf = &NodeConfig{
		NodeMonitor: &NodeMonitor{
			RefreshInterval: &metav1.Duration{
				Duration: time.Hour,
			},
			Agents: []NodeMonitorAgent{{
				Type:     NodeMonitorTypeNone,
				Endpoint: "hoge",
			}},
		},
		PowerConsumptionPredictor: &PowerConsumptionPredictor{
			Type:     PowerConsumptionPredictorTypeNone,
			Endpoint: "fuga",
		},
	}
	estConf = Estimator{
		Spec: EstimatorSpec{
			DefaultNodeConfig: defaultNodeConf,
			NodeConfigOverrides: map[string]*NodeConfig{
				"node0": nil,
				"node1": {
					NodeMonitor:               nil,
					PowerConsumptionPredictor: nil,
				},
				"node2": {
					NodeMonitor:               &NodeMonitor{},
					PowerConsumptionPredictor: &PowerConsumptionPredictor{},
				},
				"node3": node3NodeConf,
			},
		},
	}
)

func TestEstimator_MergeNodeConfig(t *testing.T) {
	type args struct {
		nodeName string
	}
	tests := []struct {
		name string
		obj  Estimator
		in   string
		want *NodeConfig
	}{
		{"nodeX", estConf, "nodeX", defaultNodeConf},
		{"node0", estConf, "node0", defaultNodeConf},
		{"node1", estConf, "node1", defaultNodeConf},
		{"node2", estConf, "node2", defaultNodeConf},
		{"node3", estConf, "node3", node3NodeConf},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.MergeNodeConfig(tt.in)
			if diff := cmp.Diff(got, tt.want); len(diff) != 0 {
				t.Errorf("Estimator.MergeNodeConfig() = %v, want %v, diff %v", got, tt.want, diff)
			}
		})
	}
}
