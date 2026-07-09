package dto

import "legacy-messenger-control-plane/internal/application/command"

type ServiceRedeployRequest struct {
	Reason string `json:"reason,omitempty"`
}

func (r ServiceRedeployRequest) ToCommand(serviceName string) command.ServiceRedeployCommand {
	return command.ServiceRedeployCommand{
		ServiceName: serviceName,
		Reason:      r.Reason,
	}
}
