/*
Copyright 2021 The KEDA Authors

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
	autoscalingv2beta2 "k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=scaledobjects,scope=Namespaced,shortName=so
// +kubebuilder:printcolumn:name="ScaleTargetKind",type="string",JSONPath=".status.scaleTargetKind"
// +kubebuilder:printcolumn:name="ScaleTargetName",type="string",JSONPath=".spec.scaleTargetRef.name"
// +kubebuilder:printcolumn:name="Min",type="integer",JSONPath=".spec.minReplicaCount"
// +kubebuilder:printcolumn:name="Max",type="integer",JSONPath=".spec.maxReplicaCount"
// +kubebuilder:printcolumn:name="Triggers",type="string",JSONPath=".spec.triggers[*].type"
// +kubebuilder:printcolumn:name="Authentication",type="string",JSONPath=".spec.triggers[*].authenticationRef.name"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status"
// +kubebuilder:printcolumn:name="Active",type="string",JSONPath=".status.conditions[?(@.type==\"Active\")].status"
// +kubebuilder:printcolumn:name="Fallback",type="string",JSONPath=".status.conditions[?(@.type==\"Fallback\")].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ScaledObject is a specification for a ScaledObject resource
type ScaledObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ScaledObjectSpec `json:"spec"`
	// +optional
	Status ScaledObjectStatus `json:"status,omitempty"`
}

// HealthStatus is the status for a ScaledObject's health
type HealthStatus struct {
	// +optional
	NumberOfFailures *int32 `json:"numberOfFailures,omitempty"`
	// +optional
	Status HealthStatusType `json:"status,omitempty"`
}

// HealthStatusType is an indication of whether the health status is happy or failing
type HealthStatusType string

const (
	// HealthStatusHappy means the status of the health object is happy
	HealthStatusHappy HealthStatusType = "Happy"

	// HealthStatusFailing means the status of the health object is failing
	HealthStatusFailing HealthStatusType = "Failing"
)

// ScaledObjectSpec is the spec for a ScaledObject resource
type ScaledObjectSpec struct {
	// 要弹性的workload，可以只指定名字
	ScaleTargetRef *ScaleTarget `json:"scaleTargetRef"`
	// 弹性的参数
	// +optional
	PollingInterval *int32 `json:"pollingInterval,omitempty"`
	// +optional
	CooldownPeriod *int32 `json:"cooldownPeriod,omitempty"`
	// +optional
	// 无触发源时的副本数量，可以小于min，如果没有配置，则最小为min
	IdleReplicaCount *int32 `json:"idleReplicaCount,omitempty"`
	// +optional
	// 如果没有配置，最小为0
	MinReplicaCount *int32 `json:"minReplicaCount,omitempty"`
	// +optional
	MaxReplicaCount *int32 `json:"maxReplicaCount,omitempty"`
	// +optional
	// 原生 HPA 的配置
	Advanced *AdvancedConfig `json:"advanced,omitempty"`

	Triggers []ScaleTriggers `json:"triggers"`
	// +optional
	Fallback *Fallback `json:"fallback,omitempty"`
}

// Fallback is the spec for fallback options
type Fallback struct {
	FailureThreshold int32 `json:"failureThreshold"`
	Replicas         int32 `json:"replicas"`
}

// AdvancedConfig specifies advance scaling options
type AdvancedConfig struct {
	// +optional
	HorizontalPodAutoscalerConfig *HorizontalPodAutoscalerConfig `json:"horizontalPodAutoscalerConfig,omitempty"`
	// +optional
	RestoreToOriginalReplicaCount bool `json:"restoreToOriginalReplicaCount,omitempty"`
}

// HorizontalPodAutoscalerConfig specifies horizontal scale config
type HorizontalPodAutoscalerConfig struct {
	// +optional
	Behavior *autoscalingv2beta2.HorizontalPodAutoscalerBehavior `json:"behavior,omitempty"`
}

// ScaleTarget holds the a reference to the scale target Object
type ScaleTarget struct {
	Name string `json:"name"`
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
	// +optional
	Kind string `json:"kind,omitempty"`
	// +optional
	EnvSourceContainerName string `json:"envSourceContainerName,omitempty"`
}

// ScaleTriggers reference the scaler that will be used
type ScaleTriggers struct {
	// Scaler 的类型，比如Kafka、MQ等
	Type string `json:"type"`
	// +optional
	Name string `json:"name,omitempty"`
	// Scaler 的配置
	Metadata map[string]string `json:"metadata"`
	// 访问 Scaler 需要的认证配置
	// +optional
	AuthenticationRef *ScaledObjectAuthRef `json:"authenticationRef,omitempty"`
	// +optional
	FallbackReplicas *int32 `json:"fallback,omitempty"`
}

// +k8s:openapi-gen=true

// ScaledObjectStatus is the status for a ScaledObject resource
// +optional
type ScaledObjectStatus struct {
	// +optional
	ScaleTargetKind string `json:"scaleTargetKind,omitempty"`
	// +optional
	ScaleTargetGVKR *GroupVersionKindResource `json:"scaleTargetGVKR,omitempty"`
	// +optional
	OriginalReplicaCount *int32 `json:"originalReplicaCount,omitempty"`
	// +optional
	LastActiveTime *metav1.Time `json:"lastActiveTime,omitempty"`
	// +optional
	// 对 HPA 暴露的 ExternalMetric
	ExternalMetricNames []string `json:"externalMetricNames,omitempty"`
	// +optional
	// 对 HPA 暴露的 ResourceMetric
	ResourceMetricNames []string `json:"resourceMetricNames,omitempty"`
	// +optional
	Conditions Conditions `json:"conditions,omitempty"`
	// +optional
	// 被弹性对象的？还是当前scale object的？
	Health map[string]HealthStatus `json:"health,omitempty"`
}

// +kubebuilder:object:root=true

// ScaledObjectList is a list of ScaledObject resources
type ScaledObjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []ScaledObject `json:"items"`
}

// ScaledObjectAuthRef points to the TriggerAuthentication or ClusterTriggerAuthentication object that
// is used to authenticate the scaler with the environment
type ScaledObjectAuthRef struct {
	Name string `json:"name"`
	// Kind of the resource being referred to. Defaults to TriggerAuthentication.
	// +optional
	Kind string `json:"kind,omitempty"`
}

func init() {
	SchemeBuilder.Register(&ScaledObject{}, &ScaledObjectList{})
}
