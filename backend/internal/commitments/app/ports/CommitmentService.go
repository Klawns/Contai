package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	"contai/internal/commitments/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type CommitmentDTO struct {
	ID                      domain.CommitmentID
	UserID                  userdomain.UserID
	Type                    domain.CommitmentType
	Description             string
	Amount                  financedomain.Money
	DueAt                   time.Time
	AccountID               accountdomain.AccountID
	CategoryID              categorydomain.CategoryID
	Note                    string
	Status                  domain.CommitmentStatus
	EffectiveStatus         domain.EffectiveStatus
	Recurrence              *domain.Recurrence
	SettledAt               *time.Time
	SettlementTransactionID *transactiondomain.TransactionID
	CanceledAt              *time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type CreateCommitmentInput struct {
	UserID      userdomain.UserID
	Type        domain.CommitmentType
	Description string
	Amount      financedomain.Money
	DueAt       time.Time
	AccountID   accountdomain.AccountID
	CategoryID  categorydomain.CategoryID
	Note        string
	Recurrence  *domain.Recurrence
}

type UpdateCommitmentInput struct {
	UserID       userdomain.UserID
	CommitmentID domain.CommitmentID
	Type         domain.CommitmentType
	Description  string
	Amount       financedomain.Money
	DueAt        time.Time
	AccountID    accountdomain.AccountID
	CategoryID   categorydomain.CategoryID
	Note         string
	Recurrence   *domain.Recurrence
}

type CancelCommitmentInput struct {
	UserID       userdomain.UserID
	CommitmentID domain.CommitmentID
	Type         domain.CommitmentType
}

type SettleCommitmentInput struct {
	UserID       userdomain.UserID
	CommitmentID domain.CommitmentID
	Type         domain.CommitmentType
	Amount       financedomain.Money
	SettledAt    time.Time
	AccountID    accountdomain.AccountID
	CategoryID   categorydomain.CategoryID
	Note         string
}

type CommitmentService interface {
	ListCommitments(ctx context.Context, input ListCommitmentsInput) ([]CommitmentDTO, error)
	CreateCommitment(ctx context.Context, input CreateCommitmentInput) (CommitmentDTO, error)
	UpdateCommitment(ctx context.Context, input UpdateCommitmentInput) (CommitmentDTO, error)
	CancelCommitment(ctx context.Context, input CancelCommitmentInput) (CommitmentDTO, error)
	SettleCommitment(ctx context.Context, input SettleCommitmentInput) (CommitmentDTO, error)
}
