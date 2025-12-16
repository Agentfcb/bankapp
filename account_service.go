package main

import (
	"fmt"
	"strings"
	"time"
)

type BankAccountService struct {
	account *Account
	storage Storage
}

func NewBankAccountService(account *Account, storage Storage) *BankAccountService {
	return &BankAccountService{
		account: account,
		storage: storage,
	}
}

func (s *BankAccountService) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	s.account.Balance += amount
	s.account.addTransaction(Transaction{
		Type:        Deposit,
		Amount:      amount,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Пополнение на сумму %.2f", amount),
	})

	return s.storage.SaveAccount(s.account)
}

func (s *BankAccountService) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if s.account.Balance < amount {
		return ErrInsufficientFunds
	}

	s.account.Balance -= amount
	s.account.addTransaction(Transaction{
		Type:        Withdraw,
		Amount:      amount,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Снятие на сумму %.2f", amount),
	})

	return s.storage.SaveAccount(s.account)
}

func (s *BankAccountService) Transfer(to *Account, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if s.account.Balance < amount {
		return ErrInsufficientFunds
	}

	if s.account.ID == to.ID {
		return ErrSameAccountTransfer
	}

	s.account.Balance -= amount
	s.account.addTransaction(Transaction{
		Type:            Transfer,
		Amount:          amount,
		Timestamp:       time.Now(),
		Description:     fmt.Sprintf("Перевод на счет %s", to.ID),
		TargetAccountID: to.ID,
	})

	
	to.Balance += amount
	to.addTransaction(Transaction{
		Type:        Transfer,
		Amount:      amount,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Перевод от счета %s", s.account.ID),
	})

	
	if err := s.storage.SaveAccount(s.account); err != nil {
		return err
	}

	return s.storage.SaveAccount(to)
}

func (s *BankAccountService) GetBalance() float64 {
	return s.account.Balance
}

func (s *BankAccountService) GetStatement() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Выписка по счету %s\n", s.account.ID))
	sb.WriteString(fmt.Sprintf("Владелец: %s\n", s.account.Owner))
	sb.WriteString(fmt.Sprintf("Дата открытия: %s\n", s.account.CreatedAt.Format("02.01.2006 15:04:05")))
	sb.WriteString(fmt.Sprintf("Текущий баланс: %.2f\n\n", s.account.Balance))
	sb.WriteString("История транзакций:\n")
	sb.WriteString(strings.Repeat("-", 80) + "\n")

	for _, t := range s.account.Transactions {
		sb.WriteString(fmt.Sprintf("%s | %s | Сумма: %.2f | %s\n",
			t.Timestamp.Format("02.01.2006 15:04:05"),
			t.Type,
			t.Amount,
			t.Description))
	}

	if len(s.account.Transactions) == 0 {
		sb.WriteString("Транзакции отсутствуют\n")
	}

	return sb.String()
}

func (a *Account) addTransaction(t Transaction) {
	t.ID = fmt.Sprintf("TX%d", time.Now().UnixNano())
	a.Transactions = append(a.Transactions, t)
}

