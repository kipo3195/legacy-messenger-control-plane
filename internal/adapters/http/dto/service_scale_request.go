package dto

import "legacy-messenger-control-plane/internal/application/command"

type ServiceScaleRequest struct {
	DesiredCount *int   `json:"desiredCount" binding:"required"`
	Reason       string `json:"reason,omitempty"`
}

func (r ServiceScaleRequest) ToCommand(serviceName string) command.ScaleServiceCommand {
	return command.ScaleServiceCommand{
		ServiceName:  serviceName,
		DesiredCount: *r.DesiredCount,
		Reason:       r.Reason,
	}
}
