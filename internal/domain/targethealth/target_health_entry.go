package targethealth

type TargetHealthEntry struct {
    TargetID         string `json:"targetId"`
    Port             int32  `json:"port"`
    AvailabilityZone string `json:"availabilityZone,omitempty"`
    State            string `json:"state"`
    Reason           string `json:"reason,omitempty"`
    Description      string `json:"description,omitempty"`
}
