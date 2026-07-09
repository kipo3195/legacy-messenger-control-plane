package aws

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

var _ ports.ELBPort = (*ELBClient)(nil)

type ELBClient struct {
	client *elbv2.Client
}

func NewELBV2Client(ctx context.Context, region string) (*ELBClient, error) {

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

	return &ELBClient{
		client: elbv2.NewFromConfig(cfg),
	}, nil
}

func (c *ELBClient) DescribeTargetHealth(
	ctx context.Context,
	targetGroupArn string,
	loadBalancerType string,
) (*domain.TargetGroupHealth, error) {
	if targetGroupArn == "" {
		return nil, fmt.Errorf("targetGroupArn is required")
	}

	out, err := c.client.DescribeTargetHealth(ctx, &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(targetGroupArn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe target health targetGroupArn=%s: %w", targetGroupArn, err)
	}

	result := &domain.TargetGroupHealth{
		TargetGroupArn:   targetGroupArn,
		TargetGroupName:  extractTargetGroupName(targetGroupArn),
		LoadBalancerType: loadBalancerType,
		Targets:          make([]domain.TargetHealthEntry, 0),
	}

	for _, desc := range out.TargetHealthDescriptions {
		entry := domain.TargetHealthEntry{
			TargetID:         aws.ToString(desc.Target.Id),
			Port:             aws.ToInt32(desc.Target.Port),
			AvailabilityZone: aws.ToString(desc.Target.AvailabilityZone),
		}

		if desc.TargetHealth != nil {
			entry.State = string(desc.TargetHealth.State)
			entry.Reason = string(desc.TargetHealth.Reason)
			entry.Description = aws.ToString(desc.TargetHealth.Description)
		}

		result.Targets = append(result.Targets, entry)
		increaseTargetGroupHealthCount(result, entry.State)
	}

	return result, nil
}

func increaseTargetGroupHealthCount(
	tg *domain.TargetGroupHealth,
	state string,
) {
	tg.Total++

	switch strings.ToLower(state) {
	case "healthy":
		tg.Healthy++
	case "unhealthy":
		tg.Unhealthy++
	case "initial":
		tg.Initial++
	case "draining":
		tg.Draining++
	case "unused":
		tg.Unused++
	case "unavailable":
		tg.Unavailable++
	}
}

func extractTargetGroupName(targetGroupArn string) string {
	parts := strings.Split(targetGroupArn, ":")
	if len(parts) == 0 {
		return ""
	}

	resource := parts[len(parts)-1]
	resourceParts := strings.Split(resource, "/")

	if len(resourceParts) >= 2 {
		return resourceParts[1]
	}

	return ""
}

func (c *ELBClient) GetLoadBalancerArnByTargetGroupArn(ctx context.Context, targetGroupArn string) (string, error) {

	if targetGroupArn == "" {
		return "", fmt.Errorf("targetGroupArn is required")
	}

	out, err := c.client.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{
		TargetGroupArns: []string{targetGroupArn},
	})
	if err != nil {
		return "", fmt.Errorf("failed to describe target group: %w", err)
	}

	if len(out.TargetGroups) == 0 {
		return "", fmt.Errorf("target group not found: %s", targetGroupArn)
	}

	tg := out.TargetGroups[0]

	if len(tg.LoadBalancerArns) == 0 {
		return "", fmt.Errorf("load balancer arn not found for target group: %s", targetGroupArn)
	}

	return tg.LoadBalancerArns[0], nil
}
