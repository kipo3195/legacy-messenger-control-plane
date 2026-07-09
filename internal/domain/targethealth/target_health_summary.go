package targethealth

type TargetHealthSummary struct {
    Total       int `json:"total"`
    Healthy     int `json:"healthy"`
    Unhealthy   int `json:"unhealthy"`
    Initial     int `json:"initial"`
    Draining    int `json:"draining"`
    Unused      int `json:"unused"`
    Unavailable int `json:"unavailable"`
}
