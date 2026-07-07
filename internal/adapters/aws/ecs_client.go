package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

var _ ports.ECSPort = (*ECSClient)(nil)

type ECSClient struct {
	client *ecs.Client
}

func NewECSClient(ctx context.Context, region string) (*ECSClient, error) {
	if region == "" {
		return nil, fmt.Errorf("aws region is required")
	}

	cfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	return &ECSClient{
		client: ecs.NewFromConfig(cfg),
	}, nil
}

func (c *ECSClient) DescribeService(ctx context.Context, clusterName string, ecsServiceName string) (*domain.ServiceStatus, error) {

	if clusterName == "" {
		return nil, fmt.Errorf("clusterName is required")
	}

	if ecsServiceName == "" {
		return nil, fmt.Errorf("ecsServiceName is required")
	}

	out, err := c.client.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  &clusterName,
		Services: []string{ecsServiceName},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe ecs service %s: %w", ecsServiceName, err)
	}

	if len(out.Services) == 0 {
		return nil, fmt.Errorf("ecs service not found: %s", ecsServiceName)
	}

	svc := out.Services[0]

	status := &domain.ServiceStatus{
		ServiceName:    ecsServiceName,
		ClusterName:    clusterName,
		Status:         ptrString(svc.Status),
		DesiredCount:   svc.DesiredCount,
		RunningCount:   svc.RunningCount,
		PendingCount:   svc.PendingCount,
		TaskDefinition: ptrString(svc.TaskDefinition),
	}

	for _, d := range svc.Deployments {
		status.Deployments = append(status.Deployments, domain.DeploymentStatus{
			Status:       ptrString(d.Status),
			DesiredCount: d.DesiredCount,
			RunningCount: d.RunningCount,
			PendingCount: d.PendingCount,
			RolloutState: string(d.RolloutState),
		})
	}

	for _, e := range svc.Events {
		status.Events = append(status.Events, domain.ServiceEvent{
			Message:   ptrString(e.Message),
			CreatedAt: ptrTime(e.CreatedAt),
		})
	}

	return status, nil
}

func ptrString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func ptrTime(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}
	return *v
}

// adpater에서 AWS라이브러리 참조 및 결과를 ->  []domain.TaskStatus 하지 않고 AWS라이브러리 호출 결과를 return하여 usecase에서 처리하게 되면
// usecase에서 외부 의존성 (AWS)를 알게 된다.
func (c *ECSClient) DescribeTask(ctx context.Context, clusterName string, ecsServiceName string, desiredStatus string) ([]domain.TaskStatus, error) {

	if clusterName == "" {
		return nil, fmt.Errorf("clusterName is required")
	}

	if ecsServiceName == "" {
		return nil, fmt.Errorf("ecsServiceName is required")
	}

	status := toECSDesiredStatus(desiredStatus)

	// 특정 서비스의 Task를 조회하려면 먼저 ListTasks로 해당 서비스에 속한 Task ARN 목록을 가져오고,
	// 그 결과를 DescribeTasks에 넘겨야 해. AWS ListTasks는 serviceName으로 필터링할 수 있고, DescribeTasks는 tasks 배열을 필수로 받는다.

	listOut, err := c.client.ListTasks(ctx, &ecs.ListTasksInput{
		Cluster:       &clusterName,
		ServiceName:   &ecsServiceName,
		DesiredStatus: status,
	})

	if err != nil {
		return nil, fmt.Errorf("[ListServiceTasks] ListTasks error: %w", err)
	}

	if len(listOut.TaskArns) == 0 {
		return nil, nil
	}

	describeOut, err := c.client.DescribeTasks(ctx, &ecs.DescribeTasksInput{
		Cluster: &clusterName,
		Tasks:   listOut.TaskArns,
	})
	if err != nil {
		return nil, fmt.Errorf("[ListServiceTasks] DescribeTasks error: %w", err)
	}

	result := make([]domain.TaskStatus, 0, len(describeOut.Tasks))

	for _, task := range describeOut.Tasks {
		result = append(result, toTaskStatus(task))
	}

	return result, nil
}

func toTaskStatus(task types.Task) domain.TaskStatus {
	return domain.TaskStatus{
		TaskID:               extractID(aws.ToString(task.TaskArn)),
		TaskArn:              aws.ToString(task.TaskArn),
		TaskDefinition:       extractID(aws.ToString(task.TaskDefinitionArn)),
		LastStatus:           aws.ToString(task.LastStatus),
		DesiredStatus:        aws.ToString(task.DesiredStatus),
		HealthStatus:         string(task.HealthStatus),
		LaunchType:           string(task.LaunchType),
		AvailabilityZone:     aws.ToString(task.AvailabilityZone),
		ContainerInstanceArn: aws.ToString(task.ContainerInstanceArn),
		CapacityProviderName: aws.ToString(task.CapacityProviderName),

		CreatedAt:  task.CreatedAt,
		StartedAt:  task.StartedAt,
		StoppingAt: task.StoppingAt,
		StoppedAt:  task.StoppedAt,

		StopCode:      string(task.StopCode),
		StoppedReason: aws.ToString(task.StoppedReason),

		Containers: toContainerStatuses(task.Containers),
		Network:    toTaskNetworkInfo(task),
	}
}

func toContainerStatuses(containers []types.Container) []domain.ContainerStatus {
	result := make([]domain.ContainerStatus, 0, len(containers))

	for _, c := range containers {
		result = append(result, domain.ContainerStatus{
			Name:              aws.ToString(c.Name),
			LastStatus:        aws.ToString(c.LastStatus),
			HealthStatus:      string(c.HealthStatus),
			Image:             aws.ToString(c.Image),
			RuntimeID:         aws.ToString(c.RuntimeId),
			Reason:            aws.ToString(c.Reason),
			ExitCode:          c.ExitCode,
			NetworkBindings:   toNetworkBindings(c.NetworkBindings),
			NetworkInterfaces: toNetworkInterfaces(c.NetworkInterfaces),
		})
	}

	return result
}

func toNetworkBindings(bindings []types.NetworkBinding) []domain.NetworkBindingInfo {
	result := make([]domain.NetworkBindingInfo, 0, len(bindings))

	for _, b := range bindings {
		result = append(result, domain.NetworkBindingInfo{
			BindIP:        aws.ToString(b.BindIP),
			ContainerPort: aws.ToInt32(b.ContainerPort),
			HostPort:      aws.ToInt32(b.HostPort),
			Protocol:      string(b.Protocol),
		})
	}

	return result
}

func toNetworkInterfaces(interfaces []types.NetworkInterface) []domain.NetworkInterfaceInfo {
	result := make([]domain.NetworkInterfaceInfo, 0, len(interfaces))

	for _, ni := range interfaces {
		result = append(result, domain.NetworkInterfaceInfo{
			AttachmentID:       aws.ToString(ni.AttachmentId),
			PrivateIPv4Address: aws.ToString(ni.PrivateIpv4Address),
		})
	}

	return result
}

func toTaskNetworkInfo(task types.Task) domain.TaskNetworkInfo {
	var bindings []domain.NetworkBindingInfo
	var interfaces []domain.NetworkInterfaceInfo
	var privateIPv4 string

	for _, c := range task.Containers {
		bindings = append(bindings, toNetworkBindings(c.NetworkBindings)...)

		containerInterfaces := toNetworkInterfaces(c.NetworkInterfaces)
		interfaces = append(interfaces, containerInterfaces...)

		for _, ni := range containerInterfaces {
			if ni.PrivateIPv4Address != "" && privateIPv4 == "" {
				privateIPv4 = ni.PrivateIPv4Address
			}
		}
	}

	modeHint := "unknown"

	if len(bindings) > 0 {
		modeHint = "bridge"
	}

	if privateIPv4 != "" {
		modeHint = "awsvpc"
	}

	return domain.TaskNetworkInfo{
		ModeHint:    modeHint,
		Bindings:    bindings,
		Interfaces:  interfaces,
		PrivateIPv4: privateIPv4,
	}
}

func extractID(arn string) string {
	if arn == "" {
		return ""
	}

	parts := strings.Split(arn, "/")
	return parts[len(parts)-1]
}

func toECSDesiredStatus(status string) types.DesiredStatus {
	switch status {
	case "STOPPED":
		return types.DesiredStatusStopped
	case "PENDING":
		return types.DesiredStatusPending
	default:
		return types.DesiredStatusRunning
	}
}
