package servicescale

type ECSServiceControlState struct {
	ECSServiceName string
	DesiredCount   int32
	RunningCount   int32
	PendingCount   int32
	Status         string
}
