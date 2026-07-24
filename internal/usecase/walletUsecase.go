package usecase

import (
	"context"

	"github.com/Wrehat/ewallet-wallet/internal/domain"
	"go.uber.org/zap"
)

type walletUsecase struct {
	repo domain.WalletRepository
	log  *zap.Logger
}

func NewWalletUsecase(rp domain.WalletRepository, log *zap.Logger) domain.WalletUsecase {
	return &walletUsecase{
		repo: rp,
		log:  log,
	}
}

func (u *walletUsecase) Create(ctx context.Context, wallet *domain.Wallet) error {
	if err := u.repo.CreateWallet(ctx, wallet); err != nil {
		u.log.Error("error create wallet to database", zap.Error(err))
		return err
	}
	return nil
}
