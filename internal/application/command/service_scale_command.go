package command

type ScaleServiceCommand struct {
	ServiceName  string
	DesiredCount int
	Reason       string
}
