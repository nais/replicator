//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigResource) DeepCopyInto(out *ConfigResource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigResource.
func (in *ConfigResource) DeepCopy() *ConfigResource {
	if in == nil {
		return nil
	}
	out := new(ConfigResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicatorConfiguration) DeepCopyInto(out *ReplicatorConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicatorConfiguration.
func (in *ReplicatorConfiguration) DeepCopy() *ReplicatorConfiguration {
	if in == nil {
		return nil
	}
	out := new(ReplicatorConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ReplicatorConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicatorConfigurationList) DeepCopyInto(out *ReplicatorConfigurationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ReplicatorConfiguration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicatorConfigurationList.
func (in *ReplicatorConfigurationList) DeepCopy() *ReplicatorConfigurationList {
	if in == nil {
		return nil
	}
	out := new(ReplicatorConfigurationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ReplicatorConfigurationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicatorConfigurationSpec) DeepCopyInto(out *ReplicatorConfigurationSpec) {
	*out = *in
	in.NamespaceSelector.DeepCopyInto(&out.NamespaceSelector)
	in.Values.DeepCopyInto(&out.Values)
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]Resource, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicatorConfigurationSpec.
func (in *ReplicatorConfigurationSpec) DeepCopy() *ReplicatorConfigurationSpec {
	if in == nil {
		return nil
	}
	out := new(ReplicatorConfigurationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicatorConfigurationStatus) DeepCopyInto(out *ReplicatorConfigurationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicatorConfigurationStatus.
func (in *ReplicatorConfigurationStatus) DeepCopy() *ReplicatorConfigurationStatus {
	if in == nil {
		return nil
	}
	out := new(ReplicatorConfigurationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resource) DeepCopyInto(out *Resource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resource.
func (in *Resource) DeepCopy() *Resource {
	if in == nil {
		return nil
	}
	out := new(Resource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Values) DeepCopyInto(out *Values) {
	*out = *in
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]ConfigResource, len(*in))
		copy(*out, *in)
	}
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]ConfigResource, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Values.
func (in *Values) DeepCopy() *Values {
	if in == nil {
		return nil
	}
	out := new(Values)
	in.DeepCopyInto(out)
	return out
}
