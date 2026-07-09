package service

type ServiceTargetGroup struct {
	TargetGroupArn string `json:"targetGroupArn"`
	ContainerName  string `json:"containerName,omitempty"`
	ContainerPort  *int32 `json:"containerPort,omitempty"`
}
