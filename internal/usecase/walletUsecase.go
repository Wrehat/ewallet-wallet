package usecase

import (
	"context"
	"errors"

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

func (u *walletUsecase) Credit(ctx context.Context, userID int, amount float64, ref string) (*domain.WalletTransaction, error) {
	_, err := u.repo.GetTransactionByReference(ctx, ref)
	if err == nil {
		return nil, domain.ErrDuplicateReference
	}

	if !errors.Is(err, domain.ErrRecordNotFound) {
		return nil, err
	}

	updatedWallet, err := u.repo.UpdateBalance(ctx, userID, amount)
	if err != nil {
		return nil, err
	}

	tx := domain.WalletTransaction{
		WalletID:              updatedWallet.ID,
		Amount:                amount,
		WalletTransactionType: domain.WalletTransactionTypeCredit,
		Reference:             ref,
	}

	if err := u.repo.CreateTransaction(ctx, &tx); err != nil {
		return nil, err
	}

	return &tx, nil

}

func (u *walletUsecase) Debit(ctx context.Context, userID int, amount float64, ref string) (*domain.WalletTransaction, error) {
	_, err := u.repo.GetTransactionByReference(ctx, ref)
	if err == nil {
		return nil, domain.ErrDuplicateReference
	}

	if !errors.Is(err, domain.ErrRecordNotFound) {
		return nil, err
	}

	updatedWallet, err := u.repo.UpdateBalance(ctx, userID, -amount)
	if err != nil {
		return nil, err
	}

	tx := domain.WalletTransaction{
		WalletID:              updatedWallet.ID,
		Amount:                amount,
		WalletTransactionType: domain.WalletTransactionTypeDebit,
		Reference:             ref,
	}

	if err := u.repo.CreateTransaction(ctx, &tx); err != nil {
		return nil, err
	}

	return &tx, nil

}

func (u *walletUsecase) GetBalance(ctx context.Context, userID int) (*domain.Wallet, error) {
	wallet, err := u.repo.GetWalletByUserID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (u *walletUsecase) GetHistory(ctx context.Context, userID int, param domain.WalletHistoryParam) ([]domain.WalletTransaction, error) {
	wallet, err := u.repo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	listTrx, err := u.repo.GetTransactionHistory(ctx, wallet.ID, param)
	if err != nil {
		return nil, err
	}

	return listTrx, nil
}
