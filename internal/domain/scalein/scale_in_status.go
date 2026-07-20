package scalein

type ScaleInStatus string

const (
	ScaleInStatusRequested ScaleInStatus = "REQUESTED"
	ScaleInStatusDraining  ScaleInStatus = "DRAINING"
	ScaleInStatusApplied   ScaleInStatus = "APPLIED"
	ScaleInStatusCompleted ScaleInStatus = "COMPLETED"
	ScaleInStatusFailed    ScaleInStatus = "FAILED"
)
