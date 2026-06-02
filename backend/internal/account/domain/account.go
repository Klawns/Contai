package domain

import (
	"regexp"
	"strings"
	"time"

	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"
)

type AccountID string

type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeDigital    AccountType = "digital"
	AccountTypeCash       AccountType = "cash"
	AccountTypeSalary     AccountType = "salary"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeOther      AccountType = "other"
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
)

type Account struct {
	ID                      AccountID
	UserID                  userdomain.UserID
	Name                    string
	Type                    AccountType
	InitialBalance          financedomain.Money
	CurrentBalance          financedomain.Money
	BankIconID              string
	IncludeInDashboardTotal bool
	Status                  AccountStatus
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

var bankIconIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

func NewAccount(id AccountID, userID userdomain.UserID, name string, accountType AccountType, initialBalance financedomain.Money, bankIconID string) (Account, error) {
	now := time.Now()
	account := Account{
		ID:                      AccountID(strings.TrimSpace(string(id))),
		UserID:                  userdomain.UserID(strings.TrimSpace(string(userID))),
		Name:                    strings.TrimSpace(name),
		Type:                    accountType,
		InitialBalance:          initialBalance,
		CurrentBalance:          initialBalance,
		BankIconID:              strings.TrimSpace(bankIconID),
		IncludeInDashboardTotal: true,
		Status:                  AccountStatusActive,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	if err := account.validate(); err != nil {
		return Account{}, err
	}

	return account, nil
}

func RehydrateAccount(id AccountID, userID userdomain.UserID, name string, accountType AccountType, initialBalance, currentBalance financedomain.Money, bankIconID string, includeInDashboardTotal bool, status AccountStatus, createdAt, updatedAt time.Time) (Account, error) {
	account := Account{
		ID:                      AccountID(strings.TrimSpace(string(id))),
		UserID:                  userdomain.UserID(strings.TrimSpace(string(userID))),
		Name:                    strings.TrimSpace(name),
		Type:                    accountType,
		InitialBalance:          initialBalance,
		CurrentBalance:          currentBalance,
		BankIconID:              strings.TrimSpace(bankIconID),
		IncludeInDashboardTotal: includeInDashboardTotal,
		Status:                  status,
		CreatedAt:               createdAt,
		UpdatedAt:               updatedAt,
	}

	if err := account.validate(); err != nil {
		return Account{}, err
	}

	return account, nil
}

func (account *Account) Edit(name string, accountType AccountType, bankIconID string, includeInDashboardTotal bool) error {
	account.Name = strings.TrimSpace(name)
	account.Type = accountType
	account.BankIconID = strings.TrimSpace(bankIconID)
	account.IncludeInDashboardTotal = includeInDashboardTotal
	account.UpdatedAt = time.Now()

	return account.validate()
}

func (account *Account) Activate() error {
	account.Status = AccountStatusActive
	account.UpdatedAt = time.Now()
	return account.validate()
}

func (account *Account) Inactivate() error {
	account.Status = AccountStatusInactive
	account.UpdatedAt = time.Now()
	return account.validate()
}

func (account *Account) IncreaseBalance(amount financedomain.Money) error {
	if !amount.IsPositive() {
		return ErrAccountMutationAmountInvalid
	}
	account.CurrentBalance = account.CurrentBalance.Add(amount)
	account.UpdatedAt = time.Now()
	return account.validate()
}

func (account *Account) DecreaseBalance(amount financedomain.Money) error {
	if !amount.IsPositive() {
		return ErrAccountMutationAmountInvalid
	}
	account.CurrentBalance = account.CurrentBalance.Sub(amount)
	account.UpdatedAt = time.Now()
	return account.validate()
}

func (account Account) validate() error {
	if strings.TrimSpace(string(account.ID)) == "" {
		return ErrAccountIDRequired
	}
	if strings.TrimSpace(string(account.UserID)) == "" {
		return ErrAccountUserIDRequired
	}
	if strings.TrimSpace(account.Name) == "" {
		return ErrAccountNameRequired
	}
	if !IsValidAccountType(account.Type) {
		return ErrAccountInvalidType
	}
	if account.Status != AccountStatusActive && account.Status != AccountStatusInactive {
		return ErrAccountInvalidStatus
	}
	if strings.TrimSpace(account.BankIconID) == "" {
		return ErrAccountBankIconIDRequired
	}
	if !bankIconIDPattern.MatchString(account.BankIconID) {
		return ErrAccountInvalidBankIconID
	}

	return nil
}

func IsValidAccountType(accountType AccountType) bool {
	switch accountType {
	case AccountTypeChecking,
		AccountTypeSavings,
		AccountTypeDigital,
		AccountTypeCash,
		AccountTypeSalary,
		AccountTypeInvestment,
		AccountTypeOther:
		return true
	default:
		return false
	}
}
