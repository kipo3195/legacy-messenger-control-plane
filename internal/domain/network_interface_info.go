package domain

type NetworkInterfaceInfo struct {
	AttachmentID       string `json:"attachmentId,omitempty"`
	PrivateIPv4Address string `json:"privateIpv4Address,omitempty"`
}
