package services

import (
	"context"

	"contai/internal/account/app/ports"
	"contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"
)

var _ ports.AccountService = AccountService{}

type AccountService struct {
	repository    ports.AccountRepository
	idGenerator   ports.AccountIDGenerator
	userValidator ports.UserValidator
}

func NewAccountService(repository ports.AccountRepository, idGenerator ports.AccountIDGenerator, userValidator ports.UserValidator) AccountService {
	return AccountService{
		repository:    repository,
		idGenerator:   idGenerator,
		userValidator: userValidator,
	}
}

func (service AccountService) CreateAccount(ctx context.Context, input ports.CreateAccountInput) (ports.AccountDTO, error) {
	if err := service.userValidator.ValidateActiveUser(ctx, input.UserID); err != nil {
		return ports.AccountDTO{}, err
	}

	account, err := domain.NewAccount(service.idGenerator.NewAccountID(), input.UserID, input.Name, input.Type, input.InitialBalance, input.BankIconID)
	if err != nil {
		return ports.AccountDTO{}, err
	}
	if input.IncludeInDashboardTotal != nil {
		account.IncludeInDashboardTotal = *input.IncludeInDashboardTotal
	}

	created, err := service.repository.CreateAccount(ctx, &account)
	if err != nil {
		return ports.AccountDTO{}, err
	}

	return toAccountDTO(*created), nil
}

func (service AccountService) ListAccounts(ctx context.Context, input ports.ListAccountsInput) ([]ports.AccountDTO, error) {
	if input.UserID == "" {
		return nil, domain.ErrAccountUserIDRequired
	}
	if input.Status != nil && *input.Status != domain.AccountStatusActive && *input.Status != domain.AccountStatusInactive {
		return nil, domain.ErrAccountInvalidStatus
	}

	accounts, err := service.repository.FindAccountsByUserID(ctx, input)
	if err != nil {
		return nil, err
	}

	return toAccountDTOs(accounts), nil
}

func (service AccountService) FindActiveAccountsByUserID(ctx context.Context, userID userdomain.UserID) ([]ports.AccountDTO, error) {
	active := domain.AccountStatusActive
	return service.ListAccounts(ctx, ports.ListAccountsInput{UserID: userID, Status: &active})
}

func (service AccountService) UpdateAccount(ctx context.Context, input ports.UpdateAccountInput) (ports.AccountDTO, error) {
	account, err := service.repository.FindAccountByID(ctx, input.AccountID, input.UserID)
	if err != nil {
		return ports.AccountDTO{}, err
	}
	if account == nil {
		return ports.AccountDTO{}, domain.ErrAccountNotFound
	}

	includeInDashboardTotal := account.IncludeInDashboardTotal
	if input.IncludeInDashboardTotal != nil {
		includeInDashboardTotal = *input.IncludeInDashboardTotal
	}

	if err := account.Edit(input.Name, input.Type, input.BankIconID, includeInDashboardTotal); err != nil {
		return ports.AccountDTO{}, err
	}

	updated, err := service.repository.UpdateAccount(ctx, account)
	if err != nil {
		return ports.AccountDTO{}, err
	}

	return toAccountDTO(*updated), nil
}

func (service AccountService) InactivateAccount(ctx context.Context, input ports.InactivateAccountInput) error {
	account, err := service.repository.FindAccountByID(ctx, input.AccountID, input.UserID)
	if err != nil {
		return err
	}
	if account == nil {
		return domain.ErrAccountNotFound
	}
	if account.Status == domain.AccountStatusInactive {
		return nil
	}

	if err := account.Inactivate(); err != nil {
		return err
	}

	_, err = service.repository.UpdateAccount(ctx, account)
	return err
}

func (service AccountService) GetTotalBalance(ctx context.Context, input ports.GetTotalBalanceInput) (financedomain.Money, error) {
	if input.UserID == "" {
		return 0, domain.ErrAccountUserIDRequired
	}

	total, err := service.repository.SumActiveAccountBalances(ctx, input.UserID)
	if err != nil {
		return 0, err
	}

	return financedomain.NewMoney(total), nil
}

func toAccountDTO(account domain.Account) ports.AccountDTO {
	return ports.AccountDTO{
		ID:                      account.ID,
		UserID:                  account.UserID,
		Name:                    account.Name,
		Type:                    account.Type,
		InitialBalance:          account.InitialBalance,
		CurrentBalance:          account.CurrentBalance,
		BankIconID:              account.BankIconID,
		IncludeInDashboardTotal: account.IncludeInDashboardTotal,
		Status:                  account.Status,
		CreatedAt:               account.CreatedAt,
		UpdatedAt:               account.UpdatedAt,
	}
}

func toAccountDTOs(accounts []domain.Account) []ports.AccountDTO {
	dtos := make([]ports.AccountDTO, 0, len(accounts))
	for _, account := range accounts {
		dtos = append(dtos, toAccountDTO(account))
	}
	return dtos
}
