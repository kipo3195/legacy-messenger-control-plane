package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server          *ServerConfig          `mapstructure:"server"`
	AWS             *AWSConfig             `mapstructure:"aws"`
	ECS             *ECSConfig             `mapstructure:"ecs"`
	Services        map[string]ServiceDef  `mapstructure:"services"`
	ServiceRegistry *ServiceRegistryConfig `mapstructure:"serviceRegistry"`
	SSH             *SSHConfig
	Redis           *RedisConfig
	AutoScale       *AutoScaleConfig
	Mock            bool
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
	ssh, err := initSsh()
	redis, err := initRedis()
	autoScaling := initAuthScaling()
	mock := initMock()

	if err != nil {
		return nil, err
	}

	return &Config{
		Server:          server,
		AWS:             aws,
		ECS:             ecs,
		ServiceRegistry: serviceRegistry,
		SSH:             ssh,
		Redis:           redis,
		AutoScale:       autoScaling,
		Mock:            mock,
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

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       int
}

func initRedis() (*RedisConfig, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	portNumber, err := strconv.Atoi(port)

	if err != nil {
		return nil, fmt.Errorf("redis portNumber data is invalid.")
	}

	return &RedisConfig{
		Host:     host,
		Port:     portNumber,
		Password: password,
	}, nil

}

type SSHConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Timeout  time.Duration
}

func initSsh() (*SSHConfig, error) {
	host := os.Getenv("SSH_HOST")
	port := os.Getenv("SSH_PORT")
	user := os.Getenv("SSH_USER")
	password := os.Getenv("SSH_PASSWORD")

	portNumber, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("ssh portNumber data is invalid.")
	}

	return &SSHConfig{
		Host:     host,
		Port:     portNumber,
		User:     user,
		Password: password,
	}, nil
}

type AutoScaleConfig struct {
	Interval            int
	SessionPerTask      int     // taske당 최대 session의 수
	ScaleOutUtilization float64 // scale out 대응해야하는 비율
	ExpiresPeriod       int     // 만료
	StopCandidatePeriod int
	MinTaskCount        int // 최소 task 수
	MaxTaskCount        int // 최대 task 수
}

func initAuthScaling() *AutoScaleConfig {
	intervalStr := os.Getenv("AUTO_SCALE_INTERVAL")

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		interval = 10
	}

	sessionPerTaskStr := os.Getenv("AUTO_SCALE_SESSION_PER_TASK")
	sessionPerTask, err := strconv.Atoi(sessionPerTaskStr)
	if err != nil {
		sessionPerTask = 1500
	}

	return &AutoScaleConfig{
		Interval:            interval,
		SessionPerTask:      sessionPerTask,
		ScaleOutUtilization: 0.8,
		ExpiresPeriod:       30,
		StopCandidatePeriod: 60,
		MinTaskCount:        1,
		MaxTaskCount:        5,
	}
}

func initMock() bool {
	return true
}
