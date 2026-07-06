package configs

import "os"

type Config struct {
	AWS      *AWSConfig            `mapstructure:"aws"`
	ECS      *ECSConfig            `mapstructure:"ecs"`
	Services map[string]ServiceDef `mapstructure:"services"`
}

type AWSConfig struct {
	Region string
}

type ECSConfig struct {
	ClusterName string
}

// Control Plane이 관리할 ECS 서비스 목록과 운영 정책을 메모리에 올려두는 객체
type ServiceDef struct {
	ECSServiceName   string `mapstructure:"ecsServiceName"`
	DisplayName      string `mapstructure:"displayName"`
	Scalable         bool   `mapstructure:"scalable"`
	MinCount         int    `mapstructure:"minCount"`
	MaxCount         int    `mapstructure:"maxCount"`
	LoadBalancerType string `mapstructure:"loadBalancerType"`
}

func Load() (*Config, error) {

	aws := initAws()

	return &Config{
		AWS: aws,
	}, nil
}

func initAws() *AWSConfig {

	region := os.Getenv("REGION")

	return &AWSConfig{
		Region: region,
	}
}
