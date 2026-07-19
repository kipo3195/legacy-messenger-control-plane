package scale

type TaskSessionInfo struct {
	SessionCountPerTask []TaskInfo

	AvgSessionCount int

	TotalSessionCount int

	ReportCoverage float64
}
