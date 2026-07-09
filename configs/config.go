package configs

import (
	"os"
)

type Config struct {
	Server          *ServerConfig          `mapstructure:"server"`
	AWS             *AWSConfig             `mapstructure:"aws"`
	ECS             *ECSConfig             `mapstructure:"ecs"`
	Services        map[string]ServiceDef  `mapstructure:"services"`
	ServiceRegistry *ServiceRegistryConfig `mapstructure:"serviceRegistry"`
}

type ServerConfig struct {
	Port string
}

type AWSConfig struct {
	Region string
}

type ECSConfig struct {
	ClusterName string
}

// Control Plane이 관리할 ECS 서비스 목록과 운영 정책을 메모리에 올려두는 객체
type ServiceDef struct {
	ECSServiceName           string `yaml:"ecsServiceName" mapstructure:"ecsServiceName"`
	DisplayName              string `yaml:"displayName" mapstructure:"displayName"`
	Scalable                 bool   `yaml:"scalable" mapstructure:"scalable"`
	MinCount                 int    `yaml:"minCount" mapstructure:"minCount"`
	MaxCount                 int    `yaml:"maxCount" mapstructure:"maxCount"`
	LoadBalancerType         string `yaml:"loadBalancerType" mapstructure:"loadBalancerType"`
	TargetConnectionsPerTask int    `yaml:"targetConnectionsPerTask" mapstructure:"targetConnectionsPerTask"`
}

type ServiceRegistryConfig struct {
	Path string `mapstructure:"path"`
}

func Load() (*Config, error) {

	server := initServer()
	aws := initAws()
	ecs := initECS()
	serviceRegistry := initServiceRegistry()

	return &Config{
		Server:          server,
		AWS:             aws,
		ECS:             ecs,
		ServiceRegistry: serviceRegistry,
	}, nil
}

func initServer() *ServerConfig {
	port := os.Getenv("PORT")
	return &ServerConfig{
		Port: port,
	}
}

func initAws() *AWSConfig {

	region := os.Getenv("REGION")

	return &AWSConfig{
		Region: region,
	}
}

func initECS() *ECSConfig {
	clusterName := os.Getenv("ECS_CLUSTER_NAME")

	return &ECSConfig{
		ClusterName: clusterName,
	}
}

func initServiceRegistry() *ServiceRegistryConfig {

	path := os.Getenv("SERVICE_REGISTRY_PATH")

	return &ServiceRegistryConfig{
		Path: path,
	}
}
