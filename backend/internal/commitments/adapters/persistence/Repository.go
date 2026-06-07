package persistence

import (
	"context"
	"errors"
	"time"

	"contai/internal/commitments/app/ports"
	"contai/internal/commitments/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.CommitmentRepository = CommitmentRepository{}

type CommitmentRepository struct {
	db *gorm.DB
}

func NewCommitmentRepository(db *gorm.DB) CommitmentRepository {
	return CommitmentRepository{db: db}
}

func (repository CommitmentRepository) WithTx(tx databaseports.TxHandle) ports.CommitmentRepository {
	if db, ok := tx.Value().(*gorm.DB); ok && db != nil {
		return CommitmentRepository{db: db}
	}
	return repository
}

func (repository CommitmentRepository) CreateCommitment(
	ctx context.Context,
	commitment *domain.Commitment,
) (*domain.Commitment, error) {
	entity := toCommitmentEntity(*commitment)
	if err := repository.db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}
	created, err := toDomainCommitment(entity)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (repository CommitmentRepository) UpdateCommitment(
	ctx context.Context,
	commitment *domain.Commitment,
) (*domain.Commitment, error) {
	entity := toCommitmentEntity(*commitment)
	result := repository.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrCommitmentNotFound
	}
	updated, err := toDomainCommitment(entity)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (repository CommitmentRepository) FindCommitmentByIDForUpdate(
	ctx context.Context,
	commitmentID domain.CommitmentID,
	userID userdomain.UserID,
) (*domain.Commitment, error) {
	var entity CommitmentEntity
	err := repository.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&entity, "id = ? AND user_id = ?", string(commitmentID), string(userID)).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	commitment, err := toDomainCommitment(entity)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}

func (repository CommitmentRepository) FindCommitmentsByUserID(
	ctx context.Context,
	input ports.ListCommitmentsInput,
) ([]domain.Commitment, error) {
	query := repository.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", string(input.UserID), string(input.Type))
	if input.StartAt != nil {
		query = query.Where("due_at >= ?", *input.StartAt)
	}
	if input.EndAt != nil {
		query = query.Where("due_at <= ?", *input.EndAt)
	}
	if input.Status != nil {
		query = query.Where("status = ?", string(*input.Status))
	}
	if input.EffectiveStatus != nil {
		query = applyEffectiveStatusFilter(query, *input.EffectiveStatus, time.Now())
	}
	if input.AccountID != nil {
		query = query.Where("account_id = ?", string(*input.AccountID))
	}
	if input.CategoryID != nil {
		query = query.Where("category_id = ?", string(*input.CategoryID))
	}
	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}
	if input.Offset > 0 {
		query = query.Offset(input.Offset)
	}

	var entities []CommitmentEntity
	if err := query.Order("due_at ASC, created_at ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	commitments := make([]domain.Commitment, 0, len(entities))
	for _, entity := range entities {
		commitment, err := toDomainCommitment(entity)
		if err != nil {
			return nil, err
		}
		commitments = append(commitments, commitment)
	}
	return commitments, nil
}

func applyEffectiveStatusFilter(
	query *gorm.DB,
	status domain.EffectiveStatus,
	now time.Time,
) *gorm.DB {
	switch status {
	case domain.EffectiveStatusOverdue:
		return query.Where("status = ? AND due_at < ?", string(domain.CommitmentStatusPending), now)
	case domain.EffectiveStatusPending:
		return query.Where("status = ? AND due_at >= ?", string(domain.CommitmentStatusPending), now)
	case domain.EffectiveStatusPaid:
		return query.Where("status = ?", string(domain.CommitmentStatusPaid))
	case domain.EffectiveStatusReceived:
		return query.Where("status = ?", string(domain.CommitmentStatusReceived))
	case domain.EffectiveStatusCanceled:
		return query.Where("status = ?", string(domain.CommitmentStatusCanceled))
	default:
		return query
	}
}
