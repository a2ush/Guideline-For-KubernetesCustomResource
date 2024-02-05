/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SampleObjectSpec defines the desired state of SampleObject
type SampleObjectSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SampleObject. Edit sampleobject_types.go to remove/update
	Filename string `json:"filename"`
	// +kubebuilder:default=BLANK
	// +optional
	Reason string `json:"reason"`
}

// SampleObjectStatus defines the observed state of SampleObject
type SampleObjectStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Filename string `json:"filename"`
	// +optional
	Reason string `json:"reason"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=sample
//+kubebuilder:printcolumn:name="FILENAME",type="string",JSONPath=".status.filename"
//+kubebuilder:printcolumn:name="REASON",type="string",JSONPath=".status.reason"

// SampleObject is the Schema for the sampleobjects API
type SampleObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SampleObjectSpec   `json:"spec,omitempty"`
	Status SampleObjectStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SampleObjectList contains a list of SampleObject
type SampleObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SampleObject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SampleObject{}, &SampleObjectList{})
}
