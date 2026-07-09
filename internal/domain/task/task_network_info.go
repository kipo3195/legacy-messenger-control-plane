package task

type TaskNetworkInfo struct {
    ModeHint    string                 `json:"modeHint"`
    Bindings    []NetworkBindingInfo   `json:"bindings,omitempty"`
    Interfaces  []NetworkInterfaceInfo `json:"interfaces,omitempty"`
    PrivateIPv4 string                 `json:"privateIpv4,omitempty"`
}
