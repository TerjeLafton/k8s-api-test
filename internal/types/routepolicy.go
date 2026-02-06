package types

type RoutePolicyConfig struct {
	Entries []string `json:"entries"`
	Name    string   `json:"name"`
}

type RoutePolicySpec struct {
	Config            RoutePolicyConfig `json:"config"`
	DeletionPolicy    DeletionPolicy    `json:"deletionPolicy"`
	Dependencies      []Dependency      `json:"dependencies,omitempty"`
	DeviceAddress     string            `json:"deviceAddress"`
	LifecyclePolicies []LifecyclePolicy `json:"lifecyclePolicies"`
	Paused            bool              `json:"paused"`
}

type RoutePolicyStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

type RoutePolicy struct {
	APIVersion string         `json:"apiVersion,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	Spec       *RoutePolicySpec   `json:"spec,omitempty"`
	Status     *RoutePolicyStatus `json:"status,omitempty"`
}
