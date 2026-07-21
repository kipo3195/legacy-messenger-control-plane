package ports

import "context"

// TaskDrainClient가 taskID만 가지고는 호출할 주소를 알 수 없기 때문에
// 별도로 endPoint (drain을 호출할)를 구하는 역할을 갖는 interface 구현
type TaskEndpointResolver interface {
	ResolveTaskEndpoint(
		ctx context.Context,
		serviceName string,
		taskID string,
	) (string, error)
}
