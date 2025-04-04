package models

import "time"

type Transaction struct {
	ID        int64     `json:"id" db:"id"`
	WalletID  string    `json:"wallet_id" db:"wallet_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Type      string    `json:"type" db:"type"` // "credit" or "debit"
	Timestamp time.Time `json:"timestamp" db:"created_at"`
}

type Wallet struct {
	ID      string  `json:"id" db:"id"`
	Balance float64 `json:"balance" db:"balance"`
}
