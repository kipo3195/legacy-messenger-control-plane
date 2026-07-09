package domain

import (
	cp "legacy-messenger-control-plane/internal/domain/connectionpressure"
	svc "legacy-messenger-control-plane/internal/domain/service"
	ss "legacy-messenger-control-plane/internal/domain/servicescale"
	th "legacy-messenger-control-plane/internal/domain/targethealth"
	tk "legacy-messenger-control-plane/internal/domain/task"
)

// domain 디렉터리 하위를 참조할 수 있도록 구성
// 실제로 해당 도메인 구조체를 사용하는 곳도 'domain' 으로 접근 가능함.

// Connection pressure
type ConnectionPressure = cp.ConnectionPressure
type ConnectionPressureMetric = cp.ConnectionPressureMetric
type ScalingRecommendation = cp.ScalingRecommendation

// Service observation
type ServiceList = svc.ServiceList
type DeploymentStatus = svc.DeploymentStatus
type ServiceStatus = svc.ServiceStatus
type ServiceEvent = svc.ServiceEvent
type ServiceTargetGroup = svc.ServiceTargetGroup

// Target health
type TargetGroupHealth = th.TargetGroupHealth
type TargetHealthEntry = th.TargetHealthEntry
type TargetHealthResponse = th.TargetHealthResponse
type TargetHealthOverallStatus = th.TargetHealthOverallStatus
type TargetHealthSummary = th.TargetHealthSummary

const (
	TargetHealthOverallHealthy       = th.TargetHealthOverallHealthy
	TargetHealthOverallDegraded      = th.TargetHealthOverallDegraded
	TargetHealthOverallTransitioning = th.TargetHealthOverallTransitioning
	TargetHealthOverallUnknown       = th.TargetHealthOverallUnknown
)

// Task observation
type TaskNetworkInfo = tk.TaskNetworkInfo
type NetworkBindingInfo = tk.NetworkBindingInfo
type NetworkInterfaceInfo = tk.NetworkInterfaceInfo
type ContainerStatus = tk.ContainerStatus
type TaskStatus = tk.TaskStatus

// service scale
type ServiceScaleResult = ss.ServiceScaleResult
type ECSServiceControlState = ss.ECSServiceControlState
type ServiceScaleCommand = ss.ServiceScaleCommand
