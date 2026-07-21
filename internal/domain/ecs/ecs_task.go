package ecs

type ECSTask struct {
	TaskID        string
	TaskARN       string
	LastStatus    string
	DesiredStatus string

	PrivateIP string
	HostIP    string
	HostPort  int
}
