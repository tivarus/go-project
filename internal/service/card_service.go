package service

import (
	"bank-api/internal/models"
	"bank-api/internal/repository"
	"bank-api/pkg/crypto"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type CardService struct {
	cardRepo    *repository.CardRepository
	accountRepo *repository.AccountRepository
	hmacSecret  string
}

func NewCardService(
	cardRepo *repository.CardRepository,
	accountRepo *repository.AccountRepository,
	hmacSecret string,
) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
		hmacSecret:  hmacSecret,
	}
}

func (s *CardService) CreateCard(userID int, req *models.CreateCardRequest) (*models.Card, error) {
	// Проверка что счет принадлежит пользователю
	account, err := s.accountRepo.GetAccountByID(req.AccountID)
	if err != nil {
		return nil, err
	}
	if account == nil || account.UserID != userID {
		return nil, errors.New("account not found or access denied")
	}

	// Генерация данных карты
	cardNumber := generateCardNumber()
	expiryDate := generateExpiryDate()
	cvv := generateCVV()

	// Шифрование данных
	encryptedNumber, err := crypto.EncryptPGP(cardNumber)
	if err != nil {
		return nil, err
	}

	encryptedExpiry, err := crypto.EncryptPGP(expiryDate)
	if err != nil {
		return nil, err
	}

	cvvHash := crypto.GenerateHMAC(cvv, s.hmacSecret)

	card := &models.Card{
		AccountID:  req.AccountID,
		Number:     encryptedNumber,
		ExpiryDate: encryptedExpiry,
		CVV:        cvvHash,
	}

	if err := s.cardRepo.CreateCard(card); err != nil {
		return nil, err
	}

	// Возвращаем карту с частично скрытыми данными
	card.Number = fmt.Sprintf("**** **** **** %s", cardNumber[len(cardNumber)-4:])
	card.ExpiryDate = expiryDate
	card.CVV = "***"

	return card, nil
}

func (s *CardService) GetCard(userID, cardID int) (*models.Card, error) {
	card, err := s.cardRepo.GetCardByID(cardID)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, errors.New("card not found")
	}

	// Проверка прав доступа
	account, err := s.accountRepo.GetAccountByID(card.AccountID)
	if err != nil {
		return nil, err
	}
	if account.UserID != userID {
		return nil, errors.New("access denied")
	}

	// Расшифровка данных для владельца
	number, err := crypto.DecryptPGP(card.Number)
	if err != nil {
		return nil, err
	}

	expiry, err := crypto.DecryptPGP(card.ExpiryDate)
	if err != nil {
		return nil, err
	}

	card.Number = number
	card.ExpiryDate = expiry

	return card, nil
}

// Вспомогательные функции генерации данных карты
func generateCardNumber() string {
	rand.Seed(time.Now().UnixNano())
	number := "4" // Visa-подобные карты
	for i := 0; i < 15; i++ {
		number += fmt.Sprintf("%d", rand.Intn(10))
	}
	return number
}

func generateExpiryDate() string {
	now := time.Now()
	return fmt.Sprintf("%02d/%02d", now.Month(), now.Year()+3)
}

func generateCVV() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%03d", rand.Intn(1000))
}

// Добавляем эти методы в конец файла credit_service.go
func (s *CreditService) GetCreditByID(id int) (*models.Credit, error) {
	credit, err := s.creditRepo.GetCreditByID(id)
	if err != nil {
		return nil, err
	}
	return credit, nil
}

func (s *CreditService) GetPaymentSchedule(creditID int) ([]*models.PaymentSchedule, error) {
	payments, err := s.creditRepo.GetPaymentSchedule(creditID)
	if err != nil {
		return nil, err
	}
	return payments, nil
}
