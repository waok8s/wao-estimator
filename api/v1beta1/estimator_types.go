package v1beta1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	OperatorName = "wao-estimator"

	DefaultNodeMonitorRefreshInterval = 30 * time.Second
)

type NodeMonitorType string

const (
	NodeMonitorTypeNone                    = "None"
	NodeMonitorTypeFake                    = "Fake"
	NodeMonitorTypeIPMIExporter            = "IPMIExporter"
	NodeMonitorTypeRedfish                 = "Redfish"
	NodeMonitorTypeDifferentialPressureAPI = "DifferentialPressureAPI"
)

type NodeMonitorAgent struct {
	Type     NodeMonitorType `json:"type"`
	Endpoint string          `json:"endpoint,omitempty"`
}

type NodeMonitor struct {
	RefreshInterval *metav1.Duration   `json:"refreshInterval,omitempty"`
	Agents          []NodeMonitorAgent `json:"agents"`
}

type PowerConsumptionPredictorType string

const (
	PowerConsumptionPredictorTypeNone     = "None"
	PowerConsumptionPredictorTypeFake     = "Fake"
	PowerConsumptionPredictorTypeMLServer = "MLServer"
	// PowerConsumptionPredictorTypeTFServing = "TFServing"
)

type PowerConsumptionPredictor struct {
	Type     PowerConsumptionPredictorType `json:"type"`
	Endpoint string                        `json:"endpoint,omitempty"`
}

type NodeConfig struct {
	NodeMonitor               *NodeMonitor               `json:"nodeMonitor,omitempty"`
	PowerConsumptionPredictor *PowerConsumptionPredictor `json:"powerConsumptionPredictor,omitempty"`
}

// EstimatorSpec defines the desired state of Estimator
type EstimatorSpec struct {
	DefaultNodeConfig   *NodeConfig            `json:"defaultNodeConfig,omitempty"`
	NodeConfigOverrides map[string]*NodeConfig `json:"nodeConfigOverrides,omitempty"`
}

func (r *Estimator) MergeNodeConfig(nodeName string) *NodeConfig {

	merged := r.Spec.DefaultNodeConfig.DeepCopy()

	v, ok := r.Spec.NodeConfigOverrides[nodeName]

	// no overrides
	if !ok || v == nil {
		return merged
	}

	overrides := v.DeepCopy()

	// override NodeMonitor
	if overrides.NodeMonitor != nil {
		if overrides.NodeMonitor.RefreshInterval != nil {
			merged.NodeMonitor.RefreshInterval = overrides.NodeMonitor.RefreshInterval
		}
		if len(overrides.NodeMonitor.Agents) != 0 {
			merged.NodeMonitor.Agents = overrides.NodeMonitor.Agents
		}
	}
	// override PowerConsumptionPredictor
	if overrides.PowerConsumptionPredictor != nil {
		if overrides.PowerConsumptionPredictor.Type != "" {
			merged.PowerConsumptionPredictor.Type = overrides.PowerConsumptionPredictor.Type
		}
		if overrides.PowerConsumptionPredictor.Endpoint != "" {
			merged.PowerConsumptionPredictor.Endpoint = overrides.PowerConsumptionPredictor.Endpoint
		}
	}

	return merged
}

// EstimatorStatus defines the observed state of Estimator
type EstimatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=est;estm

// Estimator is the Schema for the estimators API
type Estimator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EstimatorSpec   `json:"spec,omitempty"`
	Status EstimatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EstimatorList contains a list of Estimator
type EstimatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Estimator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Estimator{}, &EstimatorList{})
}
