package connectionpressure

type ConnectionPressureMetric struct {
    Namespace       string `json:"namespace"`
    MetricName      string `json:"metricName"`
    Stat            string `json:"stat"`
    PeriodSeconds   int32  `json:"periodSeconds"`
    LookbackMinutes int    `json:"lookbackMinutes"`
}
