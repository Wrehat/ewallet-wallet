package repository

import (
	"context"
	"errors"

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

func (r *walletRepo) GetTransactionByReference(ctx context.Context, ref string) (*domain.WalletTransaction, error) {
	var trx domain.WalletTransaction
	err := r.db.WithContext(ctx).Where("reference = ?", ref).First(&trx).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return &trx, nil
}

func (r *walletRepo) CreateTransaction(ctx context.Context, tx *domain.WalletTransaction) error {
	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		return err
	}

	return nil
}

func (r *walletRepo) UpdateBalance(ctx context.Context, userID int, amount float64) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Raw("SELECT id, user_id, balance FROM wallets WHERE user_id = ? FOR UPDATE", userID).Scan(&wallet).Error

		if err != nil {
			return err
		}

		if wallet.ID == 0 {
			return domain.ErrRecordNotFound
		}

		if (wallet.Balance + amount) < 0 {
			return domain.ErrInsufficientBalance
		}

		if err := tx.Exec("UPDATE wallets SET balance = balance + ? WHERE user_id = ?", amount, userID).Error; err != nil {
			return err
		}

		wallet.Balance += amount
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepo) GetWalletByUserID(ctx context.Context, userID int) (*domain.Wallet, error) {
	var wallet domain.Wallet

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return &wallet, nil
}

func (r *walletRepo) GetTransactionHistory(ctx context.Context, walletID int, param domain.WalletHistoryParam) ([]domain.WalletTransaction, error) {
	var trx []domain.WalletTransaction

	offset := (param.Page - 1) * param.Limit

	query := r.db.WithContext(ctx).Where("wallet_id = ?", walletID)

	if param.Type != "" {
		query = query.Where("wallet_transaction_type = ?", param.Type)
	}

	err := query.Order("id DESC").Limit(param.Limit).Offset(offset).Find(&trx).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrRecordNotFound
		}
		return nil, err
	}

	return trx, nil
}
