package repository

import (
	"context"

	"github.com/Wrehat/ewallet-wallet/internal/domain"
	"gorm.io/gorm"
)

type HealthRepo struct {
	db *gorm.DB
}

func NewHealthRepo(db *gorm.DB) domain.HealthCheckRepo {
	return &HealthRepo{
		db: db,
	}
}

func (r *HealthRepo) PingDB(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
