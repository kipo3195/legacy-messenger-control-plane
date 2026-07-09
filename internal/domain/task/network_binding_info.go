package task

type NetworkBindingInfo struct {
	BindIP        string `json:"bindIp,omitempty"`
	ContainerPort int32  `json:"containerPort"`
	HostPort      int32  `json:"hostPort"`
	Protocol      string `json:"protocol"`
}
