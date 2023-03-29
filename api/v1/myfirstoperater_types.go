/*
Copyright 2023.

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

// MyFirstOperaterSpec defines the desired state of MyFirstOperater
type MyFirstOperaterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of MyFirstOperater. Edit myfirstoperater_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// MyFirstOperaterStatus defines the observed state of MyFirstOperater
type MyFirstOperaterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyFirstOperater is the Schema for the myfirstoperaters API
type MyFirstOperater struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyFirstOperaterSpec   `json:"spec,omitempty"`
	Status MyFirstOperaterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyFirstOperaterList contains a list of MyFirstOperater
type MyFirstOperaterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyFirstOperater `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyFirstOperater{}, &MyFirstOperaterList{})
}
