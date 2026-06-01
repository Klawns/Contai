package server

import (
	authhttp "contai/internal/auth/adapters/http"
	authjwt "contai/internal/auth/adapters/jwt"
	authservices "contai/internal/auth/app/services"
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
}

func newDependencies(cfg config) (dependencies, error) {
	db, err := database.OpenPostgres(cfg.databaseDSN)
	if err != nil {
		return dependencies{}, err
	}

	if !isProduction() {
		if err := db.AutoMigrate(&persistence.UserEntity{}); err != nil {
			return dependencies{}, err
		}
	}

	userRepository := persistence.NewUserRepository(db)
	unitOfWork := database.NewUnitOfWork(db)
	userIDGenerator := ids.NewUUIDUserIDGenerator()
	passwordHasher := password.NewBcryptHasher()
	userService := userservices.NewUserService(userRepository, userIDGenerator, passwordHasher, nil, unitOfWork)
	jwtService := authjwt.NewService(cfg.jwtSecret, cfg.jwtAccessTTL)
	authService := authservices.NewAuthService(userRepository, passwordHasher, jwtService)
	cookieService := authhttp.NewCookieService(isProduction())

	return dependencies{
		userRepository:  userRepository,
		userIDGenerator: userIDGenerator,
		userService:     userService,
		authHandler:     authhttp.NewHandler(authService, userService, cookieService),
	}, nil
}
