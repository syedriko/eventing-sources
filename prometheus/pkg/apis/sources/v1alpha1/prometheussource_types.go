/*
Copyright 2019 The Knative Authors

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis/duck"
	duckv1beta1 "knative.dev/pkg/apis/duck/v1beta1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PrometheusSource is the Schema for the prometheussources API
// +k8s:openapi-gen=true
type PrometheusSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PrometheusSourceSpec   `json:"spec,omitempty"`
	Status PrometheusSourceStatus `json:"status,omitempty"`
}

// Check that Prometheus source can be validated and can be defaulted.
var _ runtime.Object = (*PrometheusSource)(nil)

// Check that we can create OwnerReferences to a Configuration.
var _ kmeta.OwnerRefable = (*PrometheusSource)(nil)

// Check that PrometheusSource implements the Conditions duck type.
var _ = duck.VerifyType(&PrometheusSource{}, &duckv1beta1.Conditions{})

const (
	// PrometheusSourceChangesEventType is the PrometheusSource CloudEvent type for changes.
	PrometheusSourceChangesEventType = "dev.knative.prometheus.changes"
)

// PrometheusSourceSpec defines the desired state of PrometheusSource
type PrometheusSourceSpec struct {
	// ServerURL is the URL of the Prometheus server
	ServerURL string `json:"serverURL"`

	// PromQL is the Prometheus query for this source
	PromQL string `json:"promQL"`

	// Sink is a reference to an object that will resolve to a domain
	// name to use as the sink.
	// +optional
	Sink *corev1.ObjectReference `json:"sink,omitempty"`
}

// GetGroupVersionKind returns the GroupVersionKind.
func (s *PrometheusSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("PrometheusSource")
}

// PrometheusSourceStatus defines the observed state of PrometheusSource
type PrometheusSourceStatus struct {
	// inherits duck/v1alpha1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1beta1.Status `json:",inline"`

	// SinkURI is the current active sink URI that has been configured
	// for the PrometheusSource.
	// +optional
	SinkURI string `json:"sinkUri,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PrometheusSourceList contains a list of PrometheusSource
type PrometheusSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PrometheusSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PrometheusSource{}, &PrometheusSourceList{})
}
