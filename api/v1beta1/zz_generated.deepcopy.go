//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Estimator) DeepCopyInto(out *Estimator) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Estimator.
func (in *Estimator) DeepCopy() *Estimator {
	if in == nil {
		return nil
	}
	out := new(Estimator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Estimator) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EstimatorList) DeepCopyInto(out *EstimatorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Estimator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EstimatorList.
func (in *EstimatorList) DeepCopy() *EstimatorList {
	if in == nil {
		return nil
	}
	out := new(EstimatorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EstimatorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EstimatorSpec) DeepCopyInto(out *EstimatorSpec) {
	*out = *in
	if in.DefaultNodeConfig != nil {
		in, out := &in.DefaultNodeConfig, &out.DefaultNodeConfig
		*out = new(NodeConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeConfigOverrides != nil {
		in, out := &in.NodeConfigOverrides, &out.NodeConfigOverrides
		*out = make(map[string]*NodeConfig, len(*in))
		for key, val := range *in {
			var outVal *NodeConfig
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(NodeConfig)
				(*in).DeepCopyInto(*out)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EstimatorSpec.
func (in *EstimatorSpec) DeepCopy() *EstimatorSpec {
	if in == nil {
		return nil
	}
	out := new(EstimatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EstimatorStatus) DeepCopyInto(out *EstimatorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EstimatorStatus.
func (in *EstimatorStatus) DeepCopy() *EstimatorStatus {
	if in == nil {
		return nil
	}
	out := new(EstimatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeConfig) DeepCopyInto(out *NodeConfig) {
	*out = *in
	if in.NodeMonitor != nil {
		in, out := &in.NodeMonitor, &out.NodeMonitor
		*out = new(NodeMonitor)
		(*in).DeepCopyInto(*out)
	}
	if in.PowerConsumptionPredictor != nil {
		in, out := &in.PowerConsumptionPredictor, &out.PowerConsumptionPredictor
		*out = new(PowerConsumptionPredictor)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeConfig.
func (in *NodeConfig) DeepCopy() *NodeConfig {
	if in == nil {
		return nil
	}
	out := new(NodeConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeMonitor) DeepCopyInto(out *NodeMonitor) {
	*out = *in
	if in.RefreshInterval != nil {
		in, out := &in.RefreshInterval, &out.RefreshInterval
		*out = new(v1.Duration)
		**out = **in
	}
	if in.Agents != nil {
		in, out := &in.Agents, &out.Agents
		*out = make([]NodeMonitorAgent, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeMonitor.
func (in *NodeMonitor) DeepCopy() *NodeMonitor {
	if in == nil {
		return nil
	}
	out := new(NodeMonitor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeMonitorAgent) DeepCopyInto(out *NodeMonitorAgent) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeMonitorAgent.
func (in *NodeMonitorAgent) DeepCopy() *NodeMonitorAgent {
	if in == nil {
		return nil
	}
	out := new(NodeMonitorAgent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PowerConsumptionPredictor) DeepCopyInto(out *PowerConsumptionPredictor) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PowerConsumptionPredictor.
func (in *PowerConsumptionPredictor) DeepCopy() *PowerConsumptionPredictor {
	if in == nil {
		return nil
	}
	out := new(PowerConsumptionPredictor)
	in.DeepCopyInto(out)
	return out
}
