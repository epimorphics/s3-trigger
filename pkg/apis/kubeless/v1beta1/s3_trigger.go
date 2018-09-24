package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// S3Trigger is Kubeless resource representing S3 event source
type S3Trigger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              S3TriggerSpec   `json:"spec"`
	Status            S3TriggerStatus `json:"status"`
}

// KafkaTriggerSpec defines specification for KafkaTrigger
type S3TriggerSpec struct {
	Bucket           string               `json:"bucket"`        // Trigger bucket
	SubDir           string               `json:"subDir"`        // Trigger subdirectory
	PollFrequency    int64                `json:"pollFrequency"` // S3 Polling Frequency (seconds)
	FunctionSelector metav1.LabelSelector `json:"functionSelector"`
}

type S3TriggerStatus struct {
	LastPolled string `json:"lastPolled"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KafkaTriggerList is list of KafkaTrigger's
type S3TriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	// Items is a list of third party objects
	Items []*S3Trigger `json:"items"`
}
