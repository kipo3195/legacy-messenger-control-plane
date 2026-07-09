package service

import "time"

type ServiceEvent struct {
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}
