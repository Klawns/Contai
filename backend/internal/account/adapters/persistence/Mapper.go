package persistence

import (
	"contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"
)

func toAccountEntity(account domain.Account) AccountEntity {
	return AccountEntity{
		ID:                      string(account.ID),
		UserID:                  string(account.UserID),
		Name:                    account.Name,
		Type:                    string(account.Type),
		InitialBalance:          account.InitialBalance.Cents(),
		CurrentBalance:          account.CurrentBalance.Cents(),
		BankIconID:              account.BankIconID,
		IncludeInDashboardTotal: account.IncludeInDashboardTotal,
		Status:                  string(account.Status),
		CreatedAt:               account.CreatedAt,
		UpdatedAt:               account.UpdatedAt,
	}
}

func toDomainAccount(entity AccountEntity) (domain.Account, error) {
	return domain.RehydrateAccount(
		domain.AccountID(entity.ID),
		userdomain.UserID(entity.UserID),
		entity.Name,
		domain.AccountType(entity.Type),
		financedomain.NewMoney(entity.InitialBalance),
		financedomain.NewMoney(entity.CurrentBalance),
		entity.BankIconID,
		entity.IncludeInDashboardTotal,
		domain.AccountStatus(entity.Status),
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}
