package usecase

import (
	"context"
	"fmt"
	"legacy-messenger-control-plane/configs"
	"legacy-messenger-control-plane/internal/domain"
	"legacy-messenger-control-plane/internal/ports"
	"math"
)

type sessionAutoScalingUsecase struct {
	taskSessionPort ports.TaskSessionPort
	// ecsPort         ports.ECSPort AWS 연결용
	ecsPort   ports.SessionAutoScalingECSPort // Mock
	registry  *configs.ServiceRegistry
	autoScale *configs.AutoScaleConfig
	ecsCfg    *configs.ECSConfig
}

type SessionAutoScalingUsecase interface {
	EvaluateAndScale(ctx context.Context, serviceName string) (domain.SessionAutoScalingResult, error)
}

func NewSessionAutoScalingUsecase(
	taskSessionPort ports.TaskSessionPort,
	// ecsPort ports.ECSPort,
	ecsPort ports.SessionAutoScalingECSPort,
	registry *configs.ServiceRegistry,
	ecsCfg *configs.ECSConfig,
	autoScale *configs.AutoScaleConfig,
) SessionAutoScalingUsecase {
	return &sessionAutoScalingUsecase{
		taskSessionPort: taskSessionPort,
		ecsPort:         ecsPort,
		registry:        registry,
		ecsCfg:          ecsCfg,
		autoScale:       autoScale,
	}
}

// 실제 객체를 생성하지 않고 구현체가 인터페이스 계약을 만족하는지 컴파일러에게 확인시키는 명시적 인터페이스 구현 검증문
// https://chatgpt.com/share/6a583671-bcc4-83ee-9722-5f21a67d7b99
var _ SessionAutoScalingUsecase = (*sessionAutoScalingUsecase)(nil)

func (u *sessionAutoScalingUsecase) EvaluateAndScale(ctx context.Context, serviceName string) (domain.SessionAutoScalingResult, error) {
	// [핵심 정책]
	// 1. expires가 지난 보고는 세션 합산에서 제외한다.
	// 2. 보고가 만료됐다고 해서 해당 Task의 sessionCount가 0이거나, ECS Task가 종료됐다고 판단하지 않는다.
	// 3. RUNNING Task의 보고가 만료됐거나 누락됐다면 전체 세션 집계는 실제보다 작을 수 있다.
	// 4. 보고가 불완전한 상태에서는 Scale-in을 수행하지 않는다.
	// 5. 유효한 보고만으로도 Scale-out 조건을 충족하면 Scale-out은 수행할 수 있다.
	// 6. 서비스별로 동시에 하나의 Scale-in Drain 작업만 진행한다.

	// 만료, 정상 처리 구분
	expiredTask := make([]string, 0)
	normalTask := make([]string, 0)

	// 1. 만료된 task조회
	expiredReport, err := u.taskSessionPort.GetExpiredReportTask(ctx, serviceName)
	if err != nil {
		return domain.SessionAutoScalingResult{}, fmt.Errorf("get task session report expired error")
	}

	// 2. 전체 report 조회
	reported, err := u.taskSessionPort.GetTaskSessionReport(ctx, serviceName)
	if err != nil {
		return domain.SessionAutoScalingResult{}, fmt.Errorf("failed to get task session reports: %w", err)
	}

	// 3. report 결과를 순회하면서 유효와 만료 task 분리
	for k, _ := range reported {
		taskID := k
		_, exists := expiredReport[taskID]
		if exists {
			expiredTask = append(expiredTask, taskID)
			continue
		}
		normalTask = append(normalTask, taskID)
	}

	// 유효, 만료 누락 로깅
	fmt.Printf("expired task : %s, normal task : %s\n", expiredTask, normalTask)

	// 4. 정상 보고 report 커버리지 계산
	totalSessionCount := calculateTotalSessionCount(reported, normalTask)

	requiredTaskCount := calculateRequiredTaskCount(
		totalSessionCount,
		u.autoScale.TargetSessionsPerTask,
		u.autoScale.TargetUtilization,
		u.autoScale.MinTaskCount,
		u.autoScale.MaxTaskCount,
	)

	fmt.Printf("totalSessionCount : %d, desiredCount : %d\n", totalSessionCount, requiredTaskCount)

	// 5. Redis report와 ECS Task 비교
	// yaml 파일의 서비스 정의 조회
	serviceDef, err := u.registry.Find(serviceName)
	if err != nil {
		return domain.SessionAutoScalingResult{}, fmt.Errorf("service not found: %s", serviceName)
	}

	// 현재 ECS의 상태 조회
	ecsState, err := u.ecsPort.GetServiceControlState(
		ctx,
		u.ecsCfg.ClusterName,
		serviceDef.ECSServiceName,
	)
	if err != nil {
		return domain.SessionAutoScalingResult{},
			fmt.Errorf("failed to get ECS service control state: %w", err)
	}

	// 6. 현재 desiredCount와 필요한 Task 수를 비교하여 스케일링 판단
	result := evaluateScaling(
		serviceName,
		totalSessionCount,
		int(ecsState.DesiredCount),
		requiredTaskCount,
	)

	// 7. 판단 결과에 따라 실제 작업 수행
	result, err = u.applyScalingDecision(
		ctx,
		serviceDef.ECSServiceName,
		result,
	)
	if err != nil {
		return domain.SessionAutoScalingResult{}, fmt.Errorf(
			"failed to apply scaling decision: serviceName=%s: %w",
			serviceName,
			err,
		)
	}

	return result, nil
}

func calculateTotalSessionCount(reported map[string]domain.SessionReport, normalTask []string) int {

	var totalCount int

	for _, k := range normalTask {
		report, exists := reported[k]
		if exists {
			totalCount += report.SessionCount
		}
	}

	return totalCount
}

func calculateRequiredTaskCount(
	totalSessionCount int, // task 전체에 대한 session의 수
	targetSessionsPerTask int, // task당 추구하는 session의 수
	targetUtilization float64, // 대응해야하는 비율
	minTaskCount int, // 최소 task 수
	maxTaskCount int, // 최대 task 수
) int {
	if targetSessionsPerTask <= 0 {
		return minTaskCount
	}

	effectiveCapacity :=
		float64(targetSessionsPerTask) * targetUtilization

	if effectiveCapacity <= 0 {
		return minTaskCount
	}

	// 요구되는 task의 수
	// 전체 연결 sessionCount / task당 scale out 대비 가능한 효과적인 session 수
	requiredCount := int(math.Ceil(
		float64(totalSessionCount) / effectiveCapacity,
	))

	if requiredCount < minTaskCount {
		return minTaskCount
	}

	if requiredCount > maxTaskCount {
		return maxTaskCount
	}

	return requiredCount
}

func evaluateScaling(
	serviceName string,
	totalSessionCount int,
	currentDesiredCount int,
	requiredTaskCount int,
) domain.SessionAutoScalingResult {

	switch {
	case requiredTaskCount > currentDesiredCount:
		return domain.SessionAutoScalingResult{
			ServiceName:             serviceName,
			TotalSessionCount:       totalSessionCount,
			CurrentDesiredCount:     currentDesiredCount,
			RecommendedDesiredCount: requiredTaskCount,
			Action:                  domain.ScalingActionScaleOut,
			Executed:                false,
			Reason:                  "required task count exceeds current desired count",
		}

	case requiredTaskCount < currentDesiredCount:
		return domain.SessionAutoScalingResult{
			ServiceName:             serviceName,
			TotalSessionCount:       totalSessionCount,
			CurrentDesiredCount:     currentDesiredCount,
			RecommendedDesiredCount: requiredTaskCount,
			Action:                  domain.ScalingActionScaleIn,
			Executed:                false,
			Reason:                  "scale-in requires task drain",
		}

	default:
		return domain.SessionAutoScalingResult{
			ServiceName:             serviceName,
			TotalSessionCount:       totalSessionCount,
			CurrentDesiredCount:     currentDesiredCount,
			RecommendedDesiredCount: requiredTaskCount,
			Action:                  domain.ScalingActionKeep,
			Executed:                false,
			Reason:                  "current desired count is appropriate",
		}
	}
}

func (u *sessionAutoScalingUsecase) applyScalingDecision(
	ctx context.Context,
	ecsServiceName string,
	result domain.SessionAutoScalingResult,
) (domain.SessionAutoScalingResult, error) {

	switch result.Action {
	case domain.ScalingActionScaleOut:
		updatedState, err := u.ecsPort.UpdateServiceDesiredCount(
			ctx,
			u.ecsCfg.ClusterName,
			ecsServiceName,
			result.RecommendedDesiredCount,
		)
		if err != nil {
			return domain.SessionAutoScalingResult{}, fmt.Errorf(
				"failed to scale out ECS service: "+
					"serviceName=%s, currentDesiredCount=%d, recommendedDesiredCount=%d: %w",
				result.ServiceName,
				result.CurrentDesiredCount,
				result.RecommendedDesiredCount,
				err,
			)
		}

		result.Executed = true
		result.ECSState = updatedState

	case domain.ScalingActionScaleIn:
		// Scale-in은 Task drain 절차가 필요하므로 여기서 실행하지 않는다.
		result.Executed = false

		// (scale in) 가장 적은수의 sessionCount를 갖는 task에 scale in 통보
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

	case domain.ScalingActionKeep:
		result.Executed = false

	case domain.ScalingActionNotScalable:
		result.Executed = false

	default:
		return domain.SessionAutoScalingResult{}, fmt.Errorf(
			"unsupported scaling action: %s",
			result.Action,
		)
	}

	return result, nil
}
