package configs

import "fmt"

type ServiceRegistry struct {
	clusterName string
	services    map[string]ServiceDef
}

func NewServiceRegistry(cfg *Config) (*ServiceRegistry, error) {
	if cfg.ECS.ClusterName == "" {
		return nil, fmt.Errorf("ecs clusterName is required")
	}

	if len(cfg.Services) == 0 {
		return nil, fmt.Errorf("services config is empty")
	}

	for name, svc := range cfg.Services {
		if svc.ECSServiceName == "" {
			return nil, fmt.Errorf("service %s ecsServiceName is required", name)
		}

		if svc.MinCount < 0 {
			return nil, fmt.Errorf("service %s minCount must be >= 0", name)
		}

		if svc.MaxCount < svc.MinCount {
			return nil, fmt.Errorf("service %s maxCount must be >= minCount", name)
		}
	}

	return &ServiceRegistry{
		clusterName: cfg.ECS.ClusterName,
		services:    cfg.Services,
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
