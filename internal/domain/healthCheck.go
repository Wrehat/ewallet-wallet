package domain

import "context"

type HealthCheckRepo interface {
	PingDB(ctx context.Context) error
}

type HealthCheckUsecase interface {
	Check(ctx context.Context) error
}
