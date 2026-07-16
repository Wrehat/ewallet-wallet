package domain

import "time"

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
