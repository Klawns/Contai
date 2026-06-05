package server

import (
	accounthttp "contai/internal/account/adapters/http"
	accountids "contai/internal/account/adapters/ids"
	accountpersistence "contai/internal/account/adapters/persistence"
	accountusers "contai/internal/account/adapters/users"
	accountports "contai/internal/account/app/ports"
	accountservices "contai/internal/account/app/services"
	authhttp "contai/internal/auth/adapters/http"
	authjwt "contai/internal/auth/adapters/jwt"
	authservices "contai/internal/auth/app/services"
	categoryhttp "contai/internal/category/adapters/http"
	categoryids "contai/internal/category/adapters/ids"
	categorypersistence "contai/internal/category/adapters/persistence"
	categoryports "contai/internal/category/app/ports"
	categoryservices "contai/internal/category/app/services"
	dashboardhttp "contai/internal/dashboard/adapters/http"
	dashboardpersistence "contai/internal/dashboard/adapters/persistence"
	dashboardports "contai/internal/dashboard/app/ports"
	dashboardservices "contai/internal/dashboard/app/services"
	"contai/internal/database"
	reporthttp "contai/internal/reports/adapters/http"
	reportpdf "contai/internal/reports/adapters/pdf"
	reportpersistence "contai/internal/reports/adapters/persistence"
	reportports "contai/internal/reports/app/ports"
	reportservices "contai/internal/reports/app/services"
	transactionhttp "contai/internal/transactions/adapters/http"
	transactionids "contai/internal/transactions/adapters/ids"
	transactionpersistence "contai/internal/transactions/adapters/persistence"
	transactionports "contai/internal/transactions/app/ports"
	transactionservices "contai/internal/transactions/app/services"
	"contai/internal/users/adapters/ids"
	"contai/internal/users/adapters/password"
	"contai/internal/users/adapters/persistence"
	userports "contai/internal/users/app/ports"
	userservices "contai/internal/users/app/services"
)

type dependencies struct {
	userRepository     userports.UserRepository
	userIDGenerator    userports.UserIDGenerator
	userService        userports.UserService
	authHandler        authhttp.Handler
	categoryService    categoryports.CategoryService
	categoryHandler    categoryhttp.Handler
	accountService     accountports.AccountService
	accountHandler     accounthttp.Handler
	transactionService transactionports.TransactionService
	transactionHandler transactionhttp.Handler
	dashboardService   dashboardports.DashboardService
	dashboardHandler   dashboardhttp.Handler
	reportService      reportports.ReportService
	reportHandler      reporthttp.Handler
}

func newDependencies(cfg config) (dependencies, error) {
	db, err := database.OpenPostgres(cfg.databaseDSN)
	if err != nil {
		return dependencies{}, err
	}

	if cfg.autoMigrate {
		if err := db.AutoMigrate(&persistence.UserEntity{}, &categorypersistence.CategoryEntity{}, &accountpersistence.AccountEntity{}, &transactionpersistence.TransactionEntity{}); err != nil {
			return dependencies{}, err
		}
	}

	userRepository := persistence.NewUserRepository(db)
	categoryRepository := categorypersistence.NewCategoryRepository(db)
	accountRepository := accountpersistence.NewAccountRepository(db)
	transactionRepository := transactionpersistence.NewTransactionRepository(db)
	dashboardRepository := dashboardpersistence.NewRepository(db)
	reportRepository := reportpersistence.NewReportRepository(db)
	unitOfWork := database.NewUnitOfWork(db)
	userIDGenerator := ids.NewUUIDUserIDGenerator()
	categoryIDGenerator := categoryids.NewUUIDCategoryIDGenerator()
	accountIDGenerator := accountids.NewUUIDAccountIDGenerator()
	transactionIDGenerator := transactionids.NewUUIDTransactionIDGenerator()
	passwordHasher := password.NewBcryptHasher()
	activeUserValidator := accountusers.NewActiveUserValidator(userRepository)
	categoryService := categoryservices.NewCategoryService(categoryRepository, categoryIDGenerator)
	accountService := accountservices.NewAccountService(accountRepository, accountIDGenerator, activeUserValidator)
	transactionService := transactionservices.NewTransactionService(transactionRepository, accountRepository, categoryRepository, transactionIDGenerator, unitOfWork)
	dashboardService := dashboardservices.NewDashboardService(dashboardRepository)
	reportRenderer, err := reportpdf.NewRenderer()
	if err != nil {
		return dependencies{}, err
	}
	reportService := reportservices.NewReportService(reportRepository, reportRenderer)
	defaultCategoryCreator := categoryservices.NewDefaultCategoryCreatorAdapter(categoryService)
	userService := userservices.NewUserService(userRepository, userIDGenerator, passwordHasher, defaultCategoryCreator, unitOfWork)
	jwtService := authjwt.NewService(cfg.jwtSecret, cfg.jwtAccessTTL)
	authService := authservices.NewAuthService(userRepository, passwordHasher, jwtService)
	cookieService := authhttp.NewCookieService(isProduction())

	return dependencies{
		userRepository:     userRepository,
		userIDGenerator:    userIDGenerator,
		userService:        userService,
		authHandler:        authhttp.NewHandler(authService, userService, cookieService),
		categoryService:    categoryService,
		categoryHandler:    categoryhttp.NewHandler(categoryService),
		accountService:     accountService,
		accountHandler:     accounthttp.NewHandler(accountService),
		transactionService: transactionService,
		transactionHandler: transactionhttp.NewHandler(transactionService),
		dashboardService:   dashboardService,
		dashboardHandler:   dashboardhttp.NewHandler(dashboardService),
		reportService:      reportService,
		reportHandler:      reporthttp.NewHandler(reportService),
	}, nil
}
