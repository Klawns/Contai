package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	"contai/internal/commitments/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"
)

type ListCommitmentsInput struct {
	UserID          userdomain.UserID
	Type            domain.CommitmentType
	StartAt         *time.Time
	EndAt           *time.Time
	Status          *domain.CommitmentStatus
	EffectiveStatus *domain.EffectiveStatus
	AccountID       *accountdomain.AccountID
	CategoryID      *categorydomain.CategoryID
	Limit           int
	Offset          int
}

type CommitmentRepository interface {
	WithTx(tx databaseports.TxHandle) CommitmentRepository
	CreateCommitment(ctx context.Context, commitment *domain.Commitment) (*domain.Commitment, error)
	UpdateCommitment(ctx context.Context, commitment *domain.Commitment) (*domain.Commitment, error)
	FindCommitmentByIDForUpdate(
		ctx context.Context,
		commitmentID domain.CommitmentID,
		userID userdomain.UserID,
	) (*domain.Commitment, error)
	FindCommitmentsByUserID(ctx context.Context, input ListCommitmentsInput) ([]domain.Commitment, error)
}
