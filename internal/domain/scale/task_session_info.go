package scale

type TaskSessionInfo struct {
	TaskSessionCount []TaskInfo

	AvgSessionCount int

	TotalSessionCount int

	ReportCoverage float64
}
