package server

import (
	authhttp "contai/internal/auth/adapters/http"
	authjwt "contai/internal/auth/adapters/jwt"
	authservices "contai/internal/auth/app/services"
	categoryhttp "contai/internal/category/adapters/http"
	categoryids "contai/internal/category/adapters/ids"
	categorypersistence "contai/internal/category/adapters/persistence"
	categoryports "contai/internal/category/app/ports"
	categoryservices "contai/internal/category/app/services"
	"contai/internal/database"
	"contai/internal/users/adapters/ids"
	"contai/internal/users/adapters/password"
	"contai/internal/users/adapters/persistence"
	userports "contai/internal/users/app/ports"
	userservices "contai/internal/users/app/services"
)

type dependencies struct {
	userRepository  userports.UserRepository
	userIDGenerator userports.UserIDGenerator
	userService     userports.UserService
	authHandler     authhttp.Handler
	categoryService categoryports.CategoryService
	categoryHandler categoryhttp.Handler
}

func newDependencies(cfg config) (dependencies, error) {
	db, err := database.OpenPostgres(cfg.databaseDSN)
	if err != nil {
		return dependencies{}, err
	}

	if !isProduction() {
		if err := db.AutoMigrate(&persistence.UserEntity{}, &categorypersistence.CategoryEntity{}); err != nil {
			return dependencies{}, err
		}
	}

	userRepository := persistence.NewUserRepository(db)
	categoryRepository := categorypersistence.NewCategoryRepository(db)
	unitOfWork := database.NewUnitOfWork(db)
	userIDGenerator := ids.NewUUIDUserIDGenerator()
	categoryIDGenerator := categoryids.NewUUIDCategoryIDGenerator()
	passwordHasher := password.NewBcryptHasher()
	categoryService := categoryservices.NewCategoryService(categoryRepository, categoryIDGenerator)
	defaultCategoryCreator := categoryservices.NewDefaultCategoryCreatorAdapter(categoryService)
	userService := userservices.NewUserService(userRepository, userIDGenerator, passwordHasher, defaultCategoryCreator, unitOfWork)
	jwtService := authjwt.NewService(cfg.jwtSecret, cfg.jwtAccessTTL)
	authService := authservices.NewAuthService(userRepository, passwordHasher, jwtService)
	cookieService := authhttp.NewCookieService(isProduction())

	return dependencies{
		userRepository:  userRepository,
		userIDGenerator: userIDGenerator,
		userService:     userService,
		authHandler:     authhttp.NewHandler(authService, userService, cookieService),
		categoryService: categoryService,
		categoryHandler: categoryhttp.NewHandler(categoryService),
	}, nil
}
