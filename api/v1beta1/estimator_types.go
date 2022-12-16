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
	NodeMonitorTypeNone         = "None"
	NodeMonitorTypeFake         = "Fake"
	NodeMonitorTypeIPMIExporter = "IPMIExporter"
	NodeMonitorTypeRedfish      = "Redfish"
)

type PowerConsumptionPredictorType string

const (
	PowerConsumptionPredictorTypeNone      = "None"
	PowerConsumptionPredictorTypeFake      = "Fake"
	PowerConsumptionPredictorTypeMLServer  = "MLServer"
	PowerConsumptionPredictorTypeTFServing = "TFServing"
)

type FieldRef struct {
	Label *string `json:"label,omitempty"`
}

type Field struct {
	Default  string    `json:"default"`
	Override *FieldRef `json:"override,omitempty"`
}

type NodeMonitor struct {
	Type            Field            `json:"type"`
	RefreshInterval *metav1.Duration `json:"refreshInterval,omitempty"`
	IPMIExporter    *IPMIExporter    `json:"ipmiExporter,omitempty"`
	Redfish         *Redfish         `json:"redfish,omitempty"`
}

type IPMIExporter struct {
	Endpoint Field `json:"endpoint"`
}

type Redfish struct {
	Endpoint Field `json:"endpoint"`
}

type PowerConsumptionPredictor struct {
	Type     Field     `json:"type"`
	MLServer *MLServer `json:"mlServer,omitempty"`
}

type MLServer struct {
	Endpoint Field `json:"endpoint"`
}

// EstimatorSpec defines the desired state of Estimator
type EstimatorSpec struct {
	NodeMonitor               *NodeMonitor               `json:"nodeMonitor,omitempty"`
	PowerConsumptionPredictor *PowerConsumptionPredictor `json:"powerConsumptionPredictor,omitempty"`
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
