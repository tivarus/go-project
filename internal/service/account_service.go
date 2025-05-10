package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"database/sql"
)

type AccountService struct {
	accountRepo     *repository.AccountRepository
	transactionRepo *repository.TransactionRepository
	db              *sql.DB
}

func NewAccountService(
	accountRepo *repository.AccountRepository,
	transactionRepo *repository.TransactionRepository,
	db *sql.DB,
) *AccountService {
	return &AccountService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		db:              db,
	}
}

func (s *AccountService) ProcessCreditDeposit(accountID int, amount float64, userEmail string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Обновляем баланс
	if err := s.accountRepo.UpdateBalance(accountID, amount); err != nil {
		return err
	}

	// Создаем запись о транзакции
	transaction := &models.Transaction{
		AccountID:   accountID,
		Amount:      amount,
		Type:        models.TransactionDeposit,
		Description: "Credit deposit",
	}

	if err := s.transactionRepo.CreateTransaction(tx, transaction); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *AccountService) CreateAccount(userID int, req *models.CreateAccountRequest) (*models.Account, error) {
	account := &models.Account{
		UserID:   userID,
		Balance:  0,
		Currency: req.Currency,
	}

	if err := s.accountRepo.CreateAccount(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) GetAccountByID(id int) (*models.Account, error) {
	return s.accountRepo.GetAccountByID(id)
}

func (s *AccountService) UpdateBalance(accountID int, amount float64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем существование счета
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	if account == nil {
		return ErrAccountNotFound
	}

	// Проверяем достаточность средств при снятии
	if amount < 0 && account.Balance < -amount {
		return ErrInsufficientFunds
	}

	// Обновляем баланс
	if err := s.accountRepo.UpdateBalance(accountID, amount); err != nil {
		return err
	}

	// Определяем тип транзакции
	var transactionType models.TransactionType
	if amount > 0 {
		transactionType = models.TransactionDeposit
	} else {
		transactionType = models.TransactionWithdrawal
	}

	// Создаем запись о транзакции
	transaction := &models.Transaction{
		AccountID:   accountID,
		Amount:      amount,
		Type:        transactionType,
		Description: "Balance update",
	}

	if err := s.transactionRepo.CreateTransaction(tx, transaction); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AccountService) Transfer(req *models.TransferRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем счета
	fromAccount, err := s.accountRepo.GetAccountByID(req.FromAccountID)
	if err != nil {
		return err
	}
	if fromAccount == nil {
		return ErrAccountNotFound
	}

	toAccount, err := s.accountRepo.GetAccountByID(req.ToAccountID)
	if err != nil {
		return err
	}
	if toAccount == nil {
		return ErrAccountNotFound
	}

	// Проверяем достаточность средств
	if fromAccount.Balance < req.Amount {
		return ErrInsufficientFunds
	}

	// Выполняем перевод
	if err := s.accountRepo.UpdateBalance(req.FromAccountID, -req.Amount); err != nil {
		return err
	}
	if err := s.accountRepo.UpdateBalance(req.ToAccountID, req.Amount); err != nil {
		return err
	}

	// Создаем записи о транзакциях
	fromTransaction := &models.Transaction{
		AccountID:   req.FromAccountID,
		Amount:      req.Amount,
		Type:        models.TransactionTransfer,
		Description: req.Description,
	}

	toTransaction := &models.Transaction{
		AccountID:   req.ToAccountID,
		Amount:      req.Amount,
		Type:        models.TransactionTransfer,
		Description: req.Description,
	}

	if err := s.transactionRepo.CreateTransaction(tx, fromTransaction); err != nil {
		return err
	}
	if err := s.transactionRepo.CreateTransaction(tx, toTransaction); err != nil {
		return err
	}

	return tx.Commit()
}
