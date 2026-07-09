package service

import (
	"legacy-messenger-control-plane/internal/domain"
	"math"
)

const (
	defaultTargetConnectionsPerTask = 1500
	defaultMetricPeriodSeconds      = 60
	defaultMetricLookbackMinutes    = 5
)

func evaluatePressureStatus(
	connectionPerTask float64,
	targetConnectionsPerTask int,
) domain.ConnectionPressureStatus {
	if targetConnectionsPerTask <= 0 {
		targetConnectionsPerTask = defaultTargetConnectionsPerTask
	}

	highThreshold := float64(targetConnectionsPerTask) * 0.8
	lowThreshold := float64(targetConnectionsPerTask) * 0.3

	if connectionPerTask >= highThreshold {
		return domain.ConnectionPressureStatusHigh
	}

	if connectionPerTask <= lowThreshold {
		return domain.ConnectionPressureStatusLow
	}

	return domain.ConnectionPressureStatusNormal
}

func roundFloat(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}
