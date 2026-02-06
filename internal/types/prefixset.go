package types

type PrefixSetConfig struct {
	Entries []string `json:"entries"`
	Name    string   `json:"name"`
}

type PrefixSetSpec struct {
	Config            PrefixSetConfig   `json:"config"`
	DeletionPolicy    DeletionPolicy    `json:"deletionPolicy"`
	Dependencies      []Dependency      `json:"dependencies,omitempty"`
	DeviceAddress     string            `json:"deviceAddress"`
	LifecyclePolicies []LifecyclePolicy `json:"lifecyclePolicies"`
	Paused            bool              `json:"paused"`
}

type PrefixSetStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

type PrefixSet struct {
	APIVersion string         `json:"apiVersion,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	Spec       *PrefixSetSpec   `json:"spec,omitempty"`
	Status     *PrefixSetStatus `json:"status,omitempty"`
}
