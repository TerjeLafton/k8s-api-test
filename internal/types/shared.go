package types

import "time"

type DeletionPolicy string

const (
	DeletionPolicyDelete DeletionPolicy = "Delete"
	DeletionPolicyOrphan DeletionPolicy = "Orphan"
)

type DependencyKind string

const (
	DependencyKindPrefixSet   DependencyKind = "PrefixSet"
	DependencyKindRoutePolicy DependencyKind = "RoutePolicy"
	DependencyKindVRF         DependencyKind = "VRF"
	DependencyKindBGPVRF      DependencyKind = "BGPVRF"
)

type LifecyclePolicy string

const (
	LifecyclePolicyCreate LifecyclePolicy = "Create"
	LifecyclePolicyUpdate LifecyclePolicy = "Update"
	LifecyclePolicyDelete LifecyclePolicy = "Delete"
)

type ConditionStatus string

const (
	ConditionStatusTrue    ConditionStatus = "True"
	ConditionStatusFalse   ConditionStatus = "False"
	ConditionStatusUnknown ConditionStatus = "Unknown"
)

type Dependency struct {
	Kind      DependencyKind `json:"kind"`
	Name      string         `json:"name"`
	Namespace *string        `json:"namespace,omitempty"`
}

type Condition struct {
	LastTransitionTime time.Time       `json:"lastTransitionTime"`
	Message            string          `json:"message"`
	ObservedGeneration *int64          `json:"observedGeneration,omitempty"`
	Reason             string          `json:"reason"`
	Status             ConditionStatus `json:"status"`
	Type               string          `json:"type"`
}
