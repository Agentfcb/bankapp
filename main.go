package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Application struct {
	storage        Storage
	currentAccount *Account
	accountService AccountService
}

func NewApplication() *Application {
	storage, err := NewFileStorage("accounts.json")
	if err != nil {
		fmt.Printf("Ошибка инициализации хранилища: %v\n", err)
		os.Exit(1)
	}

	return &Application{
		storage: storage,
	}
}

func (app *Application) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if app.currentAccount == nil {
			app.showMainMenu()
		} else {
			app.showAccountMenu()
		}

		fmt.Print("\nВыберите действие: ")
		scanner.Scan()
		choice := scanner.Text()

		if app.currentAccount == nil {
			app.handleMainMenuChoice(choice, scanner)
		} else {
			app.handleAccountMenuChoice(choice, scanner)
		}
	}
}

func (app *Application) showMainMenu() {
	fmt.Println("\n=== Банковское приложение ===")
	fmt.Println("1. Создать новый счет")
	fmt.Println("2. Выбрать существующий счет")
	fmt.Println("3. Просмотреть все счета")
	fmt.Println("4. Выйти")
}

func (app *Application) showAccountMenu() {
	fmt.Printf("\n=== Счет: %s (%s) ===\n", app.currentAccount.ID, app.currentAccount.Owner)
	fmt.Println("1. Пополнить счет")
	fmt.Println("2. Снять средства")
	fmt.Println("3. Перевести другому счету")
	fmt.Println("4. Просмотреть баланс")
	fmt.Println("5. Получить выписку")
	fmt.Println("6. Вернуться к выбору счета")
}

func (app *Application) handleMainMenuChoice(choice string, scanner *bufio.Scanner) {
	switch choice {
	case "1":
		app.createAccount(scanner)
	case "2":
		app.selectAccount(scanner)
	case "3":
		app.showAllAccounts()
	case "4":
		fmt.Println("До свидания!")
		os.Exit(0)
	default:
		fmt.Println("Неверный выбор. Попробуйте снова.")
	}
}

func (app *Application) handleAccountMenuChoice(choice string, scanner *bufio.Scanner) {
	switch choice {
	case "1":
		app.deposit(scanner)
	case "2":
		app.withdraw(scanner)
	case "3":
		app.transfer(scanner)
	case "4":
		app.showBalance()
	case "5":
		app.showStatement()
	case "6":
		app.currentAccount = nil
		app.accountService = nil
	default:
		fmt.Println("Неверный выбор. Попробуйте снова.")
	}
}

func (app *Application) createAccount(scanner *bufio.Scanner) {
	fmt.Print("Введите имя владельца счета: ")
	scanner.Scan()
	owner := strings.TrimSpace(scanner.Text())

	if owner == "" {
		fmt.Println("Имя владельца не может быть пустым")
		return
	}

	account := NewAccount(owner)

	if err := app.storage.SaveAccount(account); err != nil {
		fmt.Printf("Ошибка создания счета: %v\n", err)
		return
	}

	fmt.Printf("Счет успешно создан! ID: %s\n", account.ID)
	app.currentAccount = account
	app.accountService = NewBankAccountService(account, app.storage)
}

func (app *Application) selectAccount(scanner *bufio.Scanner) {
	fmt.Print("Введите ID счета: ")
	scanner.Scan()
	accountID := strings.TrimSpace(scanner.Text())

	account, err := app.storage.LoadAccount(accountID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	app.currentAccount = account
	app.accountService = NewBankAccountService(account, app.storage)
	fmt.Printf("Выбран счет: %s (%s)\n", account.ID, account.Owner)
}

func (app *Application) showAllAccounts() {
	accounts, err := app.storage.GetAllAccounts()
	if err != nil {
		fmt.Printf("Ошибка получения счетов: %v\n", err)
		return
	}

	if len(accounts) == 0 {
		fmt.Println("Счета не найдены")
		return
	}

	fmt.Println("\n=== Все счета ===")
	for _, acc := range accounts {
		fmt.Printf("ID: %s | Владелец: %s | Баланс: %.2f\n",
			acc.ID, acc.Owner, acc.Balance)
	}
}

func (app *Application) deposit(scanner *bufio.Scanner) {
	amount, err := app.readAmount(scanner, "Введите сумму для пополнения: ")
	if err != nil {
		return
	}

	if err := app.accountService.Deposit(amount); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Счет успешно пополнен!")
	}
}

func (app *Application) withdraw(scanner *bufio.Scanner) {
	amount, err := app.readAmount(scanner, "Введите сумму для снятия: ")
	if err != nil {
		return
	}

	if err := app.accountService.Withdraw(amount); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Средства успешно сняты!")
	}
}

func (app *Application) transfer(scanner *bufio.Scanner) {
	amount, err := app.readAmount(scanner, "Введите сумму для перевода: ")
	if err != nil {
		return
	}

	fmt.Print("Введите ID счета получателя: ")
	scanner.Scan()
	targetAccountID := strings.TrimSpace(scanner.Text())

	targetAccount, err := app.storage.LoadAccount(targetAccountID)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	if err := app.accountService.Transfer(targetAccount, amount); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Println("Перевод выполнен успешно!")
	}
}

func (app *Application) showBalance() {
	balance := app.accountService.GetBalance()
	fmt.Printf("Текущий баланс: %.2f\n", balance)
}

func (app *Application) showStatement() {
	statement := app.accountService.GetStatement()
	fmt.Println(statement)
}

func (app *Application) readAmount(scanner *bufio.Scanner, prompt string) (float64, error) {
	fmt.Print(prompt)
	scanner.Scan()

	amountStr := strings.TrimSpace(scanner.Text())
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Ошибка: введите корректное число")
		return 0, err
	}

	return amount, nil
}

func main() {
	app := NewApplication()
	app.Run()
}
