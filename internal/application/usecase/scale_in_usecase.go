package usecase

import (
	"context"
	"legacy-messenger-control-plane/internal/ports"
)

type scaleInUsecase struct {
	taskSessionPort ports.TaskSessionPort
	ecsPort         ports.ECSPort
	scalingPolicy   *ScalingPolicy
	coordinator     *ScaleInCoordinator
}

type ScaleInUsecase interface {
	Process(ctx context.Context) error
}

func NewScaleInUsecase(
	taskSessionPort ports.TaskSessionPort,
	ecsPort ports.ECSPort,
	scalingPolicy *ScalingPolicy,
	coordinator *ScaleInCoordinator,
) ScaleInUsecase {
	return &scaleInUsecase{
		taskSessionPort: taskSessionPort,
		ecsPort:         ecsPort,
		scalingPolicy:   scalingPolicy,
		coordinator:     coordinator,
	}
}

func (u *scaleInUsecase) Process(ctx context.Context) error {
	jobs := u.coordinator.GetActiveJobs()

	for _, job := range jobs {
		if err := u.processJob(ctx, job); err != nil {
			u.coordinator.MarkFailed(
				job.ServiceName,
				err,
			)
		}
	}

	return nil
}
