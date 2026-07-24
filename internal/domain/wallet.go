package domain

import (
	"context"
	"time"
)

type Wallet struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"column:user_id;unique"`
	Balance   float64   `json:"balance" gorm:"column:balance;type:decimal(15,2)"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (*Wallet) TableName() string {
	return "wallets"
}

type WalletTransaction struct {
	ID                    int       `json:"-" gorm:"primaryKey;autoIncrement"`
	WalletID              int       `json:"wallet_id" gorm:"column:wallet_id"`
	Amount                float64   `json:"amount" gorm:"column:amount;type:decimal(15,2)"`
	WalletTransactionType string    `json:"wallet_transaction_type" gorm:"column:wallet_transaction_type;type:varchar(10)"`
	Reference             string    `json:"reference" gorm:"column:reference;type:varchar(100);uniqueIndex"`
	CreatedAt             time.Time `json:"date"`
	UpdatedAt             time.Time `json:"-"`
}

func (*WalletTransaction) TableName() string {
	return "wallet_transactions"
}

type TokenData struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type UMSGateway interface {
	ValidateToken(ctx context.Context, token string) (*TokenData, error)
}

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet *Wallet) error
	GetTransactionByReference(ctx context.Context, ref string) (*WalletTransaction, error)
	UpdateBalance(ctx context.Context, userID int, amount float64) (*Wallet, error)
	CreateTransaction(ctx context.Context, tx *WalletTransaction) error
}

type WalletUsecase interface {
	Create(ctx context.Context, wallet *Wallet) error
	Credit(ctx context.Context, userID int, amount float64, ref string) (*WalletTransaction, error)
	Debit(ctx context.Context, userID int, amount float64, ref string) (*WalletTransaction, error)
}
