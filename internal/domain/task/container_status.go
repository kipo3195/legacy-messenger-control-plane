package task

type ContainerStatus struct {
	Name         string `json:"name"`
	LastStatus   string `json:"lastStatus"`
	HealthStatus string `json:"healthStatus"`
	Image        string `json:"image"`
	RuntimeID    string `json:"runtimeId,omitempty"`
	Reason       string `json:"reason,omitempty"`
	ExitCode     *int32 `json:"exitCode,omitempty"`

	NetworkBindings   []NetworkBindingInfo   `json:"networkBindings"`
	NetworkInterfaces []NetworkInterfaceInfo `json:"networkInterfaces"`
}
