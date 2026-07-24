package repository

import (
	"context"

	"github.com/Wrehat/ewallet-wallet/internal/domain"
	"gorm.io/gorm"
)

type walletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) domain.WalletRepository {
	return &walletRepo{
		db: db,
	}
}

func (r *walletRepo) CreateWallet(ctx context.Context, wallet *domain.Wallet) error {
	return r.db.WithContext(ctx).Create(wallet).Error
}
