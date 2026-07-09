package connectionpressure

type ConnectionPressureStatus string

const (
	ConnectionPressureStatusLow    ConnectionPressureStatus = "LOW"
	ConnectionPressureStatusNormal ConnectionPressureStatus = "NORMAL"
	ConnectionPressureStatusHigh   ConnectionPressureStatus = "HIGH"
)
