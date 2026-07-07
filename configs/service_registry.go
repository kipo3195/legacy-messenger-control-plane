package configs

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// ECS에 올라가는 service 메타데이터 저장소
type ServiceRegistry struct {
	clusterName string
	services    map[string]ServiceDef
}

type serviceRegistryFile struct {
	Services map[string]ServiceDef `yaml:"services"`
}

func NewServiceRegistry(cfg *Config) (*ServiceRegistry, error) {
	if cfg.ECS.ClusterName == "" {
		return nil, fmt.Errorf("ecs clusterName is required")
	}

	fmt.Println("Service Registry Path:", cfg.ServiceRegistry.Path)
	bytes, err := os.ReadFile(cfg.ServiceRegistry.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read service registry file: %w", err)
	}

	var file serviceRegistryFile
	if err := yaml.Unmarshal(bytes, &file); err != nil {
		return nil, fmt.Errorf("failed to parse service registry file: %w", err)
	}

	if len(file.Services) == 0 {
		return nil, fmt.Errorf("service registry is empty")
	}

	fmt.Printf("services: %+v\n", file.Services)

	return &ServiceRegistry{
		services: file.Services,
	}, nil

}

func (r *ServiceRegistry) ClusterName() string {
	return r.clusterName
}

func (r *ServiceRegistry) Find(serviceName string) (ServiceDef, error) {
	svc, ok := r.services[serviceName]
	if !ok {
		return ServiceDef{}, fmt.Errorf("service %s not found", serviceName)
	}

	return svc, nil
}

func (r *ServiceRegistry) List() map[string]ServiceDef {
	return r.services
}
