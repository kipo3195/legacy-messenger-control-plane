package servicescale

type ServiceScaleCommand struct {
	ServiceName  string
	DesiredCount int
	Reason       string
}
