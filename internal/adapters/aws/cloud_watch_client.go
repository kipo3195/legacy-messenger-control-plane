package aws

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cwtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

var ErrCloudWatchMetricNotFound = errors.New("cloudwatch metric datapoint not found")

type CloudWatchClient struct {
	client *cloudwatch.Client
}

func NewCloudWatchClient(ctx context.Context, region string) (*CloudWatchClient, error) {

	// configs.Config 안의 region 값을 이용해서 AWS SDK용 config를 한 번 생성
	awsCfg, err := awsconfig.LoadDefaultConfig(
		ctx,
		awsconfig.WithRegion(region),
	)

	if err != nil {
		return nil, err
	}

	cloudWatchClient := &CloudWatchClient{
		client: cloudwatch.NewFromConfig(awsCfg),
	}

	return cloudWatchClient, nil
}

func (c *CloudWatchClient) GetALBActiveConnectionCount(ctx context.Context, loadBalancerArn string, periodSeconds int32, lookback time.Duration) (float64, error) {
	if loadBalancerArn == "" {
		return 0, fmt.Errorf("loadBalancerArn is required")
	}

	if periodSeconds <= 0 {
		periodSeconds = 60
	}

	if lookback <= 0 {
		lookback = 5 * time.Minute
	}

	loadBalancerDimension, err := extractLoadBalancerDimension(loadBalancerArn)
	if err != nil {
		return 0, err
	}

	endTime := time.Now().UTC()
	startTime := endTime.Add(-lookback)

	// AWS SDK for Go v2의 CloudWatch client는 GetMetricData도 지원하지만,
	// 지금처럼 단일 지표 하나만 조회할 때는 GetMetricStatistics로도 충분
	out, err := c.client.GetMetricStatistics(ctx, &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String("AWS/ApplicationELB"),
		MetricName: aws.String("ActiveConnectionCount"),
		Dimensions: []cwtypes.Dimension{
			{
				Name:  aws.String("LoadBalancer"),
				Value: aws.String(loadBalancerDimension),
			},
		},
		StartTime:  aws.Time(startTime),
		EndTime:    aws.Time(endTime),
		Period:     aws.Int32(periodSeconds),
		Statistics: []cwtypes.Statistic{cwtypes.StatisticAverage},
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get ALB ActiveConnectionCount metric: %w", err)
	}

	if len(out.Datapoints) == 0 {
		return 0, ErrCloudWatchMetricNotFound
	}

	sort.Slice(out.Datapoints, func(i, j int) bool {
		return out.Datapoints[i].Timestamp.Before(*out.Datapoints[j].Timestamp)
	})

	latest := out.Datapoints[len(out.Datapoints)-1]

	if latest.Average == nil {
		return 0, ErrCloudWatchMetricNotFound
	}

	return *latest.Average, nil
}

func extractLoadBalancerDimension(loadBalancerArn string) (string, error) {
	const marker = ":loadbalancer/"

	idx := strings.Index(loadBalancerArn, marker)
	if idx == -1 {
		return "", fmt.Errorf("invalid loadBalancerArn: %s", loadBalancerArn)
	}

	return loadBalancerArn[idx+len(marker):], nil
}
