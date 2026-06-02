package ports

import (
	"context"
	"time"

	"contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"
)

type AccountDTO struct {
	ID                      domain.AccountID
	UserID                  userdomain.UserID
	Name                    string
	Type                    domain.AccountType
	InitialBalance          financedomain.Money
	CurrentBalance          financedomain.Money
	BankIconID              string
	IncludeInDashboardTotal bool
	Status                  domain.AccountStatus
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type CreateAccountInput struct {
	UserID                  userdomain.UserID
	Name                    string
	Type                    domain.AccountType
	InitialBalance          financedomain.Money
	BankIconID              string
	IncludeInDashboardTotal *bool
}

type UpdateAccountInput struct {
	UserID                  userdomain.UserID
	AccountID               domain.AccountID
	Name                    string
	Type                    domain.AccountType
	BankIconID              string
	IncludeInDashboardTotal *bool
}

type ListAccountsInput struct {
	UserID userdomain.UserID
	Status *domain.AccountStatus
}

type InactivateAccountInput struct {
	UserID    userdomain.UserID
	AccountID domain.AccountID
}

type GetTotalBalanceInput struct {
	UserID userdomain.UserID
}

type AccountService interface {
	CreateAccount(ctx context.Context, input CreateAccountInput) (AccountDTO, error)
	ListAccounts(ctx context.Context, input ListAccountsInput) ([]AccountDTO, error)
	FindActiveAccountsByUserID(ctx context.Context, userID userdomain.UserID) ([]AccountDTO, error)
	UpdateAccount(ctx context.Context, input UpdateAccountInput) (AccountDTO, error)
	InactivateAccount(ctx context.Context, input InactivateAccountInput) error
	GetTotalBalance(ctx context.Context, input GetTotalBalanceInput) (financedomain.Money, error)
}
