package main

import (
	"bank-api/internal/config"
	"bank-api/internal/handlers"
	"bank-api/internal/repository"
	"bank-api/internal/service"
	"bank-api/pkg/crypto"
	"bank-api/pkg/database"
	"bank-api/pkg/logging"
	"bank-api/pkg/mail"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	logger := logging.New()
	logger.Info("Starting banking API application")

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация PGP
	if err := crypto.InitPGP(); err != nil {
		logger.Fatalf("Failed to initialize PGP: %v", err)
	}

	// Подключение к БД
	db, err := database.Connect(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Successfully connected to database")

	// Инициализация почтового сервиса
	mailer := mail.NewMailer(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.SMTPFrom,
	)

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	cardRepo := repository.NewCardRepository(db)
	creditRepo := repository.NewCreditRepository(db)

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	notificationService := service.NewNotificationService(mailer)
	// transactionService := service.NewTransactionService(
	// 	transactionRepo,
	// 	accountRepo,
	// 	notificationService,
	// 	db,
	// )
	accountService := service.NewAccountService(accountRepo, transactionRepo, db)
	cardService := service.NewCardService(cardRepo, accountRepo, cfg.HMACSecret)
	creditService := service.NewCreditService(
		creditRepo,
		accountRepo,
		accountService,
		notificationService,
		db,
	)

	// Запуск шедулера для обработки платежей
	go StartScheduler(creditService)

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(authService)
	accountHandler := handlers.NewAccountHandler(accountService)
	cardHandler := handlers.NewCardHandler(cardService)
	creditHandler := handlers.NewCreditHandler(
		creditService,
		accountRepo,
	)

	router := mux.NewRouter()

	// Публичные маршруты
	publicRouter := router.PathPrefix("/api").Subrouter()
	authHandler.RegisterRoutes(publicRouter)

	// Защищенные маршруты
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	accountHandler.RegisterRoutes(protectedRouter)
	cardHandler.RegisterRoutes(protectedRouter)
	creditHandler.RegisterRoutes(protectedRouter)

	logger.Infof("Server is running on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, router))
}

func StartScheduler(creditSvc *service.CreditService) {
	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if err := creditSvc.ProcessDuePayments(); err != nil {
			log.Printf("Error processing due payments: %v", err)
		}
	}
}
