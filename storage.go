package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type FileStorage struct {
	filePath string
	accounts map[string]*Account
	mu       sync.RWMutex
}

func NewFileStorage(filePath string) (*FileStorage, error) {
	storage := &FileStorage{
		filePath: filePath,
		accounts: make(map[string]*Account),
	}

	if err := storage.loadFromFile(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return storage, nil
}

func (fs *FileStorage) SaveAccount(account *Account) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.accounts[account.ID] = account
	return fs.saveToFile()
}

func (fs *FileStorage) LoadAccount(accountID string) (*Account, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	account, exists := fs.accounts[accountID]
	if !exists {
		return nil, ErrAccountNotFound
	}

	// Создаем копию для безопасности
	return &Account{
		ID:           account.ID,
		Owner:        account.Owner,
		Balance:      account.Balance,
		Transactions: append([]Transaction{}, account.Transactions...),
		CreatedAt:    account.CreatedAt,
	}, nil
}

func (fs *FileStorage) GetAllAccounts() ([]*Account, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	accounts := make([]*Account, 0, len(fs.accounts))
	for _, account := range fs.accounts {
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (fs *FileStorage) saveToFile() error {
	data, err := json.MarshalIndent(fs.accounts, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}

	return os.WriteFile(fs.filePath, data, 0644)
}

func (fs *FileStorage) loadFromFile() error {
	data, err := os.ReadFile(fs.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &fs.accounts)
}
