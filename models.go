package main

import (
	"fmt"
	"time"
)

type TransactionType string

const (
	Deposit  TransactionType = "DEPOSIT"
	Withdraw TransactionType = "WITHDRAW"
	Transfer TransactionType = "TRANSFER"
)

type Transaction struct {
	ID              string
	Type            TransactionType
	Amount          float64
	Timestamp       time.Time
	Description     string
	TargetAccountID string // Для переводов
}

type Account struct {
	ID           string
	Owner        string
	Balance      float64
	Transactions []Transaction
	CreatedAt    time.Time
}

func NewAccount(owner string) *Account {
	return &Account{
		ID:        generateID(),
		Owner:     owner,
		Balance:   0,
		CreatedAt: time.Now(),
	}
}

func generateID() string {
	return fmt.Sprintf("ACC%d", time.Now().UnixNano())
}
