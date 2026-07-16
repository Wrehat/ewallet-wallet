package usecase

import (
	"context"

	"github.com/Wrehat/ewallet-wallet/internal/domain"
)

type HealthUsecase struct {
	repo domain.HealthCheckRepo
}

func NewHealthUsecase(rp domain.HealthCheckRepo) domain.HealthCheckUsecase {
	return &HealthUsecase{
		repo: rp,
	}
}

func (uc *HealthUsecase) Check(ctx context.Context) error {
	if err := uc.repo.PingDB(ctx); err != nil {
		return err
	}
	return nil
}
