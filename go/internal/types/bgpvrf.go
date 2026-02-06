package types

type BGPVRFConfig struct {
	ASNumber        string  `json:"asNumber"`
	IPv4Address     string  `json:"ipv4Address"`
	MaximumPathEBGP int     `json:"maximumPathEBGP"`
	MaximumPathIBGP int     `json:"maximumPathIBGP"`
	RD              int     `json:"rd"`
	VRFName         *string `json:"vrfName,omitempty"`
}

type BGPVRFSpec struct {
	Config            BGPVRFConfig      `json:"config"`
	DeletionPolicy    DeletionPolicy    `json:"deletionPolicy"`
	Dependencies      []Dependency      `json:"dependencies,omitempty"`
	DeviceAddress     string            `json:"deviceAddress"`
	LifecyclePolicies []LifecyclePolicy `json:"lifecyclePolicies"`
	Paused            bool              `json:"paused"`
}

type BGPVRFStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

type BGPVRF struct {
	APIVersion string         `json:"apiVersion,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	Spec       *BGPVRFSpec    `json:"spec,omitempty"`
	Status     *BGPVRFStatus  `json:"status,omitempty"`
}
