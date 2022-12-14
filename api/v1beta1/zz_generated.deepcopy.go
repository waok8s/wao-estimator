//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
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
	in.NodeMonitor.DeepCopyInto(&out.NodeMonitor)
	in.PowerConsumptionPredictor.DeepCopyInto(&out.PowerConsumptionPredictor)
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
func (in *Field) DeepCopyInto(out *Field) {
	*out = *in
	if in.Override != nil {
		in, out := &in.Override, &out.Override
		*out = new(FieldRef)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Field.
func (in *Field) DeepCopy() *Field {
	if in == nil {
		return nil
	}
	out := new(Field)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FieldRef) DeepCopyInto(out *FieldRef) {
	*out = *in
	if in.Label != nil {
		in, out := &in.Label, &out.Label
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FieldRef.
func (in *FieldRef) DeepCopy() *FieldRef {
	if in == nil {
		return nil
	}
	out := new(FieldRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IPMIExporter) DeepCopyInto(out *IPMIExporter) {
	*out = *in
	in.Endpoint.DeepCopyInto(&out.Endpoint)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IPMIExporter.
func (in *IPMIExporter) DeepCopy() *IPMIExporter {
	if in == nil {
		return nil
	}
	out := new(IPMIExporter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MLServer) DeepCopyInto(out *MLServer) {
	*out = *in
	in.Endpoint.DeepCopyInto(&out.Endpoint)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MLServer.
func (in *MLServer) DeepCopy() *MLServer {
	if in == nil {
		return nil
	}
	out := new(MLServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeMonitor) DeepCopyInto(out *NodeMonitor) {
	*out = *in
	in.Type.DeepCopyInto(&out.Type)
	if in.IPMIExporter != nil {
		in, out := &in.IPMIExporter, &out.IPMIExporter
		*out = new(IPMIExporter)
		(*in).DeepCopyInto(*out)
	}
	if in.Redfish != nil {
		in, out := &in.Redfish, &out.Redfish
		*out = new(Redfish)
		(*in).DeepCopyInto(*out)
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
func (in *PowerConsumptionPredictor) DeepCopyInto(out *PowerConsumptionPredictor) {
	*out = *in
	in.Type.DeepCopyInto(&out.Type)
	if in.MLServer != nil {
		in, out := &in.MLServer, &out.MLServer
		*out = new(MLServer)
		(*in).DeepCopyInto(*out)
	}
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

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Redfish) DeepCopyInto(out *Redfish) {
	*out = *in
	in.Endpoint.DeepCopyInto(&out.Endpoint)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Redfish.
func (in *Redfish) DeepCopy() *Redfish {
	if in == nil {
		return nil
	}
	out := new(Redfish)
	in.DeepCopyInto(out)
	return out
}
