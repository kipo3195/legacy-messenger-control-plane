package command

type TaskSessionReportCommand struct {
	ServiceName  string
	TaskID       string
	SessionCount int
}
