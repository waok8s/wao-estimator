package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	OperatorName = "wao-estimator"
)

type NodeMonitorType string

const (
	NodeMonitorTypeNone         = "None"
	NodeMonitorTypeIPMIExporter = "IPMIExporter"
	NodeMonitorTypeRedfish      = "Redfish"
)

var nodeMonitorTypeCollection = map[NodeMonitorType]struct{}{
	NodeMonitorTypeNone:         {},
	NodeMonitorTypeIPMIExporter: {},
	NodeMonitorTypeRedfish:      {},
}

type PowerConsumptionPredictorType string

const (
	PowerConsumptionPredictorTypeNone      = "None"
	PowerConsumptionPredictorTypeMLServer  = "MLServer"
	PowerConsumptionPredictorTypeTFServing = "TFServing"
)

var powerConsumptionPredictorTypeCollection = map[PowerConsumptionPredictorType]struct{}{
	PowerConsumptionPredictorTypeNone:     {},
	PowerConsumptionPredictorTypeMLServer: {},
	// PowerConsumptionPredictorTypeTFServing: {}, // not implemented
}

type FieldRef struct {
	Label *string `json:"label,omitempty"`
}

type Field struct {
	Default  string    `json:"default,omitempty"`
	Override *FieldRef `json:"override,omitempty"`
}

type NodeMonitor struct {
	Type         Field         `json:"type"`
	IPMIExporter *IPMIExporter `json:"ipmiExporter,omitempty"`
	Redfish      *Redfish      `json:"redfish,omitempty"`
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
	NodeMonitor               NodeMonitor               `json:"nodeMonitor,omitempty"`
	PowerConsumptionPredictor PowerConsumptionPredictor `json:"powerConsumptionPredictor,omitempty"`
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
