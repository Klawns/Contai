package services

import (
	"context"
	"strings"

	databaseports "contai/internal/database/ports"
	"contai/internal/users/app/ports"
	"contai/internal/users/domain"
)

var _ ports.UserService = UserService{}

type UserService struct {
	repository             ports.UserRepository
	idGenerator            ports.UserIDGenerator
	passwordHasher         ports.PasswordHasher
	defaultCategoryCreator ports.DefaultCategoryCreator
	unitOfWork             databaseports.UnitOfWork
}

func NewUserService(
	repository ports.UserRepository,
	idGenerator ports.UserIDGenerator,
	passwordHasher ports.PasswordHasher,
	defaultCategoryCreator ports.DefaultCategoryCreator,
	unitOfWork databaseports.UnitOfWork,
) UserService {
	return UserService{
		repository:             repository,
		idGenerator:            idGenerator,
		passwordHasher:         passwordHasher,
		defaultCategoryCreator: defaultCategoryCreator,
		unitOfWork:             unitOfWork,
	}
}

func (service UserService) CreateUser(ctx context.Context, input ports.CreateUserInput) (ports.UserDTO, error) {
	if len(strings.TrimSpace(input.PlainPassword)) < 8 {
		return ports.UserDTO{}, domain.ErrUserPasswordTooWeak
	}

	if service.unitOfWork != nil {
		var user ports.UserDTO
		err := service.unitOfWork.WithinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
			txUser, err := service.createUser(txCtx, service.repository.WithTx(tx), service.defaultCategoryCreatorWithTx(tx), input)
			if err != nil {
				return err
			}
			user = txUser
			return nil
		})
		if err != nil {
			return ports.UserDTO{}, err
		}
		return user, nil
	}

	return service.createUser(ctx, service.repository, service.defaultCategoryCreator, input)
}

func (service UserService) createUser(ctx context.Context, repository ports.UserRepository, defaultCategoryCreator ports.DefaultCategoryCreator, input ports.CreateUserInput) (ports.UserDTO, error) {
	email := domain.NormalizeEmail(input.Email)
	exists, err := repository.EmailExists(ctx, email)
	if err != nil {
		return ports.UserDTO{}, err
	}
	if exists {
		return ports.UserDTO{}, domain.ErrUserEmailAlreadyExists
	}

	passwordHash, err := service.passwordHasher.HashPassword(ctx, input.PlainPassword)
	if err != nil {
		return ports.UserDTO{}, err
	}

	user, err := domain.NewUser(service.idGenerator.NewUserID(), input.Name, email, passwordHash)
	if err != nil {
		return ports.UserDTO{}, err
	}

	created, err := repository.CreateUser(ctx, &user)
	if err != nil {
		return ports.UserDTO{}, err
	}

	if defaultCategoryCreator != nil {
		if err := defaultCategoryCreator.EnsureDefaultCategories(ctx, created.ID); err != nil {
			return ports.UserDTO{}, err
		}
	}

	return toUserDTO(*created), nil
}

func (service UserService) defaultCategoryCreatorWithTx(tx databaseports.TxHandle) ports.DefaultCategoryCreator {
	if service.defaultCategoryCreator == nil {
		return nil
	}

	return service.defaultCategoryCreator.WithTx(tx)
}

func (service UserService) GetUserByID(ctx context.Context, userID domain.UserID) (ports.UserDTO, error) {
	user, err := service.repository.FindUserById(ctx, userID)
	if err != nil {
		return ports.UserDTO{}, err
	}
	if user == nil {
		return ports.UserDTO{}, nil
	}

	return toUserDTO(*user), nil
}

func toUserDTO(user domain.User) ports.UserDTO {
	return ports.UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
}
