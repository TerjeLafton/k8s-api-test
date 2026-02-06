package types

type RouteTarget struct {
	ID        int  `json:"id"`
	Name      int  `json:"name"`
	Stitching bool `json:"stitching"`
}

type Export struct {
	RoutePolicyName *string       `json:"routePolicyName,omitempty"`
	RouteTargets    []RouteTarget `json:"routeTargets,omitempty"`
	ToAllowImportedVPN bool      `json:"toAllowImportedVPN"`
}

type Import struct {
	FromVRFAdvertiseAsVPN bool          `json:"fromVRFAdvertiseAsVPN"`
	RoutePolicyName       *string       `json:"routePolicyName,omitempty"`
	RouteTargets          []RouteTarget `json:"routeTargets,omitempty"`
}

type IPUnicast struct {
	Export *Export `json:"export,omitempty"`
	Import *Import `json:"import,omitempty"`
}

type VRFConfig struct {
	Description string     `json:"description"`
	IPv4        *IPUnicast `json:"ipv4,omitempty"`
	IPv6        *IPUnicast `json:"ipv6,omitempty"`
	Name        string     `json:"name"`
}

type VRFSpec struct {
	Config            VRFConfig         `json:"config"`
	DeletionPolicy    DeletionPolicy    `json:"deletionPolicy"`
	Dependencies      []Dependency      `json:"dependencies,omitempty"`
	DeviceAddress     string            `json:"deviceAddress"`
	LifecyclePolicies []LifecyclePolicy `json:"lifecyclePolicies"`
	Paused            bool              `json:"paused"`
}

type VRFStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

type VRF struct {
	APIVersion string         `json:"apiVersion,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	Spec       *VRFSpec       `json:"spec,omitempty"`
	Status     *VRFStatus     `json:"status,omitempty"`
}
