package usecase

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
)

type sessionAutoScalingUsecase struct {
	taskSessionPort ports.TaskSessionPort
	ecsPort         ports.ECSPort
	registry        *configs.ServiceRegistry
}

type SessionAutoScalingUsecase interface {
	EvaluateAndScale(ctx context.Context, serviceName string) (domain.SessionAutoScalingResult, error)
}

func NewSessionAutoScalingUsecase(
	taskSessionPort ports.TaskSessionPort,
	ecsPort ports.ECSPort,
	registry *configs.ServiceRegistry,
) SessionAutoScalingUsecase {
	return &sessionAutoScalingUsecase{
		taskSessionPort: taskSessionPort,
		ecsPort:         ecsPort,
		registry:        registry,
	}
}

func (u *sessionAutoScalingUsecase) EvaluateAndScale(ctx context.Context, serviceName string) (domain.SessionAutoScalingResult, error) {
	// [핵심 정책]
	// 1. expires가 지난 보고는 세션 합산에서 제외한다.
	// 2. 보고가 만료됐다고 해서 해당 Task의 sessionCount가 0이거나, ECS Task가 종료됐다고 판단하지 않는다.
	// 3. RUNNING Task의 보고가 만료됐거나 누락됐다면 전체 세션 집계는 실제보다 작을 수 있다.
	// 4. 보고가 불완전한 상태에서는 Scale-in을 수행하지 않는다.
	// 5. 유효한 보고만으로도 Scale-out 조건을 충족하면 Scale-out은 수행할 수 있다.
	// 6. 서비스별로 동시에 하나의 Scale-in Drain 작업만 진행한다.

	// [프로세스]
	// 1. 만료된 task조회 (map형태)
	expiredMap, err := u.taskSessionPort.GetExpiredReportTask(ctx, serviceName)

	fmt.Println("여기여기여기 : ", expiredMap)

	if err != nil {
		return domain.SessionAutoScalingResult{}, fmt.Errorf("get task session report expired error")
	}
	// 2. report 조회
	reportedList, err := u.taskSessionPort.GetTaskSessionReport(ctx, serviceName)

	// 3. report 결과를 순회하면서 유효 report 만료 report 분리
	for _, value := range reportedList {
		taskID := value.TaskID
		_, exists := expiredMap[taskID]
		if exists {
			fmt.Printf("serviceName : %s task : %s is expired\n", serviceName, taskID)
			continue
		}
		fmt.Printf("exists task :  %s, data : %s \n", taskID, value)
	}

	// 4. Redis report와 ECS Task 비교
	// 정상 보고, 보고 누락

	// 5. 정상 보고 report 커버리지 계산 및 scale out, in 판단

	// 6. (scale in) 가장 적은수의 sessionCount를 갖는 task에 scale in 통보
	// desiredCount만 변경한다고해서 선정한 Task가 종료된다고 보장되지 않습니다.
	// ECS Service는 desiredCount를 유지하는 역할을 하며, Scale-in 시 어떤 Task가 종료될지는 ECS 스케줄러가 결정합니다.
	// 또한 Service Task를 직접 중지해도 desiredCount가 그대로라면 ECS는 대체 Task를 시작됩니다.
	// 그러므로 Task draining 절차를 통해 안전하게 Scale-in 필요.

	// [Task Drain 프로세스]
	// a. Scale-in 대상 Task 선정
	// b. 전체 RUNNING Task 보호 (protection=true, Expire을 정상 범주 내로 선정함 Drain이 진행되는 동안 별도의 Scale-in이나 배포가 발생하면 대상 Task 또는 다른 Task가 먼저 종료되는 것을 방지)
	// c. 대상 Task에 Drain 요청 (신규 요청 차단, 기존 session drain)
	// d. 대상 Task의 sessionCount == 0 확인 (2~3회 확인), Redis 보고가 만료되지 않아야함, ECS Task 상태가 RUNNING -> 이 조건에 모두 부합했을때
	// e. 대상 Task의 protection만 해제
	// f. 나머지 Task가 모두 protected인지 확인
	// g. desiredCount를 1 감소 (변경 실패시 어떻게 처리할 것인가?)
	// h. 대상 Task가 STOPPED인지 확인
	// i. 서비스가 안정 상태인지 확인
	// j. 남은 Task의 protection 해제

	// redis 에서 값 조회
	u.taskSessionPort.GetTaskSessionReport(ctx, serviceName)

	// task 값을 기준으로 scaling 판단

	// desired count 변경 처리

	return domain.SessionAutoScalingResult{}, nil
}
