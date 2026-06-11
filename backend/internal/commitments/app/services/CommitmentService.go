package services

import (
	"context"
	"time"

	accountports "contai/internal/account/app/ports"
	accountdomain "contai/internal/account/domain"
	categoryports "contai/internal/category/app/ports"
	categorydomain "contai/internal/category/domain"
	"contai/internal/commitments/app/ports"
	"contai/internal/commitments/domain"
	databaseports "contai/internal/database/ports"
	financedomain "contai/internal/finance/domain"
	transactionports "contai/internal/transactions/app/ports"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

var _ ports.CommitmentService = CommitmentService{}

type CommitmentService struct {
	commitmentRepository   ports.CommitmentRepository
	transactionRepository  transactionports.TransactionRepository
	accountRepository      accountports.AccountRepository
	categoryRepository     categoryports.CategoryRepository
	commitmentIDGenerator  ports.CommitmentIDGenerator
	transactionIDGenerator transactionports.TransactionIDGenerator
	unitOfWork             databaseports.UnitOfWork
}

func NewCommitmentService(
	commitmentRepository ports.CommitmentRepository,
	transactionRepository transactionports.TransactionRepository,
	accountRepository accountports.AccountRepository,
	categoryRepository categoryports.CategoryRepository,
	commitmentIDGenerator ports.CommitmentIDGenerator,
	transactionIDGenerator transactionports.TransactionIDGenerator,
	unitOfWork databaseports.UnitOfWork,
) CommitmentService {
	return CommitmentService{
		commitmentRepository:   commitmentRepository,
		transactionRepository:  transactionRepository,
		accountRepository:      accountRepository,
		categoryRepository:     categoryRepository,
		commitmentIDGenerator:  commitmentIDGenerator,
		transactionIDGenerator: transactionIDGenerator,
		unitOfWork:             unitOfWork,
	}
}

func (service CommitmentService) ListCommitments(
	ctx context.Context,
	input ports.ListCommitmentsInput,
) ([]ports.CommitmentDTO, error) {
	if input.UserID == "" {
		return nil, domain.ErrCommitmentUserIDRequired
	}
	if !isValidCommitmentType(input.Type) {
		return nil, domain.ErrCommitmentInvalidType
	}
	if input.Status != nil && !isValidCommitmentStatus(*input.Status) {
		return nil, domain.ErrCommitmentInvalidStatus
	}
	if input.EffectiveStatus != nil && !isValidEffectiveStatus(*input.EffectiveStatus) {
		return nil, domain.ErrCommitmentInvalidStatus
	}
	if input.Limit < 0 || input.Offset < 0 {
		return nil, domain.ErrCommitmentInvalidStatus
	}

	commitments, err := service.commitmentRepository.FindCommitmentsByUserID(ctx, input)
	if err != nil {
		return nil, err
	}
	return toCommitmentDTOs(commitments, time.Now()), nil
}

func (service CommitmentService) CreateCommitment(
	ctx context.Context,
	input ports.CreateCommitmentInput,
) (ports.CommitmentDTO, error) {
	fields := editableFields(input.Description, input.Amount, input.DueAt, input.AccountID, input.CategoryID, input.Note, input.Recurrence)
	var commitment domain.Commitment
	var err error
	switch input.Type {
	case domain.CommitmentTypePayable:
		commitment, err = domain.NewPayable(service.commitmentIDGenerator.NewCommitmentID(), input.UserID, fields)
	case domain.CommitmentTypeReceivable:
		commitment, err = domain.NewReceivable(service.commitmentIDGenerator.NewCommitmentID(), input.UserID, fields)
	default:
		return ports.CommitmentDTO{}, domain.ErrCommitmentInvalidType
	}
	if err != nil {
		return ports.CommitmentDTO{}, err
	}

	if err := service.validateReferences(ctx, service.accountRepository, service.categoryRepository, commitment); err != nil {
		return ports.CommitmentDTO{}, err
	}
	created, err := service.commitmentRepository.CreateCommitment(ctx, &commitment)
	if err != nil {
		return ports.CommitmentDTO{}, err
	}
	return toCommitmentDTO(*created, time.Now()), nil
}

func (service CommitmentService) UpdateCommitment(
	ctx context.Context,
	input ports.UpdateCommitmentInput,
) (ports.CommitmentDTO, error) {
	var dto ports.CommitmentDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		commitmentRepository := service.commitmentRepository.WithTx(tx)
		accountRepository := service.accountRepository.WithTx(tx)
		categoryRepository := service.categoryRepository.WithTx(tx)

		commitment, err := findPendingCommitment(
			txCtx,
			commitmentRepository,
			input.CommitmentID,
			input.UserID,
			input.Type,
		)
		if err != nil {
			return err
		}
		fields := editableFields(input.Description, input.Amount, input.DueAt, input.AccountID, input.CategoryID, input.Note, input.Recurrence)
		if err := commitment.Edit(fields); err != nil {
			return err
		}
		if err := service.validateReferences(txCtx, accountRepository, categoryRepository, *commitment); err != nil {
			return err
		}
		updated, err := commitmentRepository.UpdateCommitment(txCtx, commitment)
		if err != nil {
			return err
		}
		dto = toCommitmentDTO(*updated, time.Now())
		return nil
	})
	if err != nil {
		return ports.CommitmentDTO{}, err
	}
	return dto, nil
}

func (service CommitmentService) CancelCommitment(
	ctx context.Context,
	input ports.CancelCommitmentInput,
) (ports.CommitmentDTO, error) {
	var dto ports.CommitmentDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		commitmentRepository := service.commitmentRepository.WithTx(tx)
		commitment, err := findPendingCommitment(
			txCtx,
			commitmentRepository,
			input.CommitmentID,
			input.UserID,
			input.Type,
		)
		if err != nil {
			return err
		}
		if err := commitment.Cancel(); err != nil {
			return err
		}
		updated, err := commitmentRepository.UpdateCommitment(txCtx, commitment)
		if err != nil {
			return err
		}
		dto = toCommitmentDTO(*updated, time.Now())
		return nil
	})
	if err != nil {
		return ports.CommitmentDTO{}, err
	}
	return dto, nil
}

func (service CommitmentService) SettleCommitment(
	ctx context.Context,
	input ports.SettleCommitmentInput,
) (ports.CommitmentDTO, error) {
	var dto ports.CommitmentDTO
	err := service.withinTx(ctx, func(txCtx context.Context, tx databaseports.TxHandle) error {
		commitmentRepository := service.commitmentRepository.WithTx(tx)
		transactionRepository := service.transactionRepository.WithTx(tx)
		accountRepository := service.accountRepository.WithTx(tx)
		categoryRepository := service.categoryRepository.WithTx(tx)

		commitment, err := findPendingCommitment(
			txCtx,
			commitmentRepository,
			input.CommitmentID,
			input.UserID,
			input.Type,
		)
		if err != nil {
			return err
		}

		transaction, err := service.newSettlementTransaction(*commitment, input)
		if err != nil {
			return err
		}
		if err := service.validateTransactionReferences(txCtx, accountRepository, categoryRepository, transaction); err != nil {
			return err
		}
		if err := service.applyEffects(txCtx, accountRepository, transaction.UserID, transaction.BalanceEffects()); err != nil {
			return err
		}
		createdTransaction, err := transactionRepository.CreateTransaction(txCtx, &transaction)
		if err != nil {
			return err
		}
		if input.Type == domain.CommitmentTypePayable {
			err = commitment.MarkPaid(createdTransaction.ID, input.SettledAt)
		} else {
			err = commitment.MarkReceived(createdTransaction.ID, input.SettledAt)
		}
		if err != nil {
			return err
		}
		updated, err := commitmentRepository.UpdateCommitment(txCtx, commitment)
		if err != nil {
			return err
		}
		dto = toCommitmentDTO(*updated, time.Now())
		return nil
	})
	if err != nil {
		return ports.CommitmentDTO{}, err
	}
	return dto, nil
}

func (service CommitmentService) newSettlementTransaction(
	commitment domain.Commitment,
	input ports.SettleCommitmentInput,
) (transactiondomain.Transaction, error) {
	switch input.Type {
	case domain.CommitmentTypePayable:
		accountID := input.AccountID
		transaction, err := transactiondomain.NewExpense(
			service.transactionIDGenerator.NewTransactionID(),
			input.UserID,
			commitment.Description,
			input.Amount,
			input.SettledAt,
			&accountID,
			input.CategoryID,
			transactiondomain.SettlementStatusSettled,
			&input.SettledAt,
			transactiondomain.RecurrenceTypeNone,
			nil,
			input.Note,
		)
		if err != nil {
			return transactiondomain.Transaction{}, err
		}
		err = transaction.SetOrigin(transactiondomain.TransactionOriginTypePayable, string(commitment.ID))
		return transaction, err
	case domain.CommitmentTypeReceivable:
		accountID := input.AccountID
		transaction, err := transactiondomain.NewIncome(
			service.transactionIDGenerator.NewTransactionID(),
			input.UserID,
			commitment.Description,
			input.Amount,
			input.SettledAt,
			&accountID,
			input.CategoryID,
			transactiondomain.SettlementStatusSettled,
			&input.SettledAt,
			transactiondomain.RecurrenceTypeNone,
			nil,
			input.Note,
		)
		if err != nil {
			return transactiondomain.Transaction{}, err
		}
		err = transaction.SetOrigin(transactiondomain.TransactionOriginTypeReceivable, string(commitment.ID))
		return transaction, err
	default:
		return transactiondomain.Transaction{}, domain.ErrCommitmentInvalidType
	}
}

func (service CommitmentService) validateReferences(
	ctx context.Context,
	accountRepository accountports.AccountRepository,
	categoryRepository categoryports.CategoryRepository,
	commitment domain.Commitment,
) error {
	account, err := accountRepository.FindAccountByID(ctx, commitment.AccountID, commitment.UserID)
	if err != nil {
		return err
	}
	if account == nil || account.Status != accountdomain.AccountStatusActive {
		return domain.ErrCommitmentAccountNotFound
	}

	category, err := categoryRepository.FindCategoryByID(ctx, commitment.CategoryID, commitment.UserID)
	if err != nil {
		return err
	}
	if category == nil || category.Status != categorydomain.CategoryStatusActive {
		return domain.ErrCommitmentCategoryNotFound
	}
	if commitment.Type == domain.CommitmentTypePayable && category.Type != categorydomain.CategoryTypeExpense {
		return domain.ErrCommitmentCategoryTypeMismatch
	}
	if commitment.Type == domain.CommitmentTypeReceivable && category.Type != categorydomain.CategoryTypeIncome {
		return domain.ErrCommitmentCategoryTypeMismatch
	}
	return nil
}

func (service CommitmentService) validateTransactionReferences(
	ctx context.Context,
	accountRepository accountports.AccountRepository,
	categoryRepository categoryports.CategoryRepository,
	transaction transactiondomain.Transaction,
) error {
	account, err := accountRepository.FindAccountByID(ctx, *transaction.AccountID, transaction.UserID)
	if err != nil {
		return err
	}
	if account == nil || account.Status != accountdomain.AccountStatusActive {
		return domain.ErrCommitmentAccountNotFound
	}

	category, err := categoryRepository.FindCategoryByID(ctx, *transaction.CategoryID, transaction.UserID)
	if err != nil {
		return err
	}
	if category == nil || category.Status != categorydomain.CategoryStatusActive {
		return domain.ErrCommitmentCategoryNotFound
	}
	if transaction.Type == transactiondomain.TransactionTypeExpense && category.Type != categorydomain.CategoryTypeExpense {
		return domain.ErrCommitmentCategoryTypeMismatch
	}
	if transaction.Type == transactiondomain.TransactionTypeIncome && category.Type != categorydomain.CategoryTypeIncome {
		return domain.ErrCommitmentCategoryTypeMismatch
	}
	return nil
}

func (service CommitmentService) applyEffects(
	ctx context.Context,
	accountRepository accountports.AccountRepository,
	userID userdomain.UserID,
	effects []transactiondomain.BalanceEffect,
) error {
	for _, effect := range effects {
		account, err := accountRepository.FindAccountByIDForUpdate(ctx, effect.AccountID, userID)
		if err != nil {
			return err
		}
		if account == nil {
			return domain.ErrCommitmentAccountNotFound
		}
		if effect.Amount.IsPositive() {
			if err := account.IncreaseBalance(effect.Amount); err != nil {
				return err
			}
		} else {
			if err := account.DecreaseBalance(effect.Amount.Neg()); err != nil {
				return err
			}
		}
		if _, err := accountRepository.UpdateAccount(ctx, account); err != nil {
			return err
		}
	}
	return nil
}

func (service CommitmentService) withinTx(ctx context.Context, fn func(context.Context, databaseports.TxHandle) error) error {
	if service.unitOfWork == nil {
		return fn(ctx, databaseports.NewTxHandle(nil))
	}
	return service.unitOfWork.WithinTx(ctx, fn)
}

func findPendingCommitment(
	ctx context.Context,
	repository ports.CommitmentRepository,
	commitmentID domain.CommitmentID,
	userID userdomain.UserID,
	commitmentType domain.CommitmentType,
) (*domain.Commitment, error) {
	commitment, err := repository.FindCommitmentByIDForUpdate(ctx, commitmentID, userID)
	if err != nil {
		return nil, err
	}
	if commitment == nil {
		return nil, domain.ErrCommitmentNotFound
	}
	if commitment.Type != commitmentType {
		return nil, domain.ErrCommitmentNotFound
	}
	if commitment.Status != domain.CommitmentStatusPending {
		return nil, domain.ErrCommitmentNotPending
	}
	return commitment, nil
}

func editableFields(
	description string,
	amount financedomain.Money,
	dueAt time.Time,
	accountID accountdomain.AccountID,
	categoryID categorydomain.CategoryID,
	note string,
	recurrence *domain.Recurrence,
) domain.EditableFields {
	return domain.EditableFields{
		Description: description,
		Amount:      amount,
		DueAt:       dueAt,
		AccountID:   accountID,
		CategoryID:  categoryID,
		Note:        note,
		Recurrence:  recurrence,
	}
}

func toCommitmentDTO(commitment domain.Commitment, now time.Time) ports.CommitmentDTO {
	return ports.CommitmentDTO{
		ID:                      commitment.ID,
		UserID:                  commitment.UserID,
		Type:                    commitment.Type,
		Description:             commitment.Description,
		Amount:                  commitment.Amount,
		DueAt:                   commitment.DueAt,
		AccountID:               commitment.AccountID,
		CategoryID:              commitment.CategoryID,
		Note:                    commitment.Note,
		Status:                  commitment.Status,
		EffectiveStatus:         commitment.EffectiveStatus(now),
		Recurrence:              commitment.Recurrence,
		SettledAt:               commitment.SettledAt,
		SettlementTransactionID: commitment.SettlementTransactionID,
		CanceledAt:              commitment.CanceledAt,
		CreatedAt:               commitment.CreatedAt,
		UpdatedAt:               commitment.UpdatedAt,
	}
}

func toCommitmentDTOs(
	commitments []domain.Commitment,
	now time.Time,
) []ports.CommitmentDTO {
	dtos := make([]ports.CommitmentDTO, 0, len(commitments))
	for _, commitment := range commitments {
		dtos = append(dtos, toCommitmentDTO(commitment, now))
	}
	return dtos
}

func isValidCommitmentType(commitmentType domain.CommitmentType) bool {
	return commitmentType == domain.CommitmentTypePayable || commitmentType == domain.CommitmentTypeReceivable
}

func isValidCommitmentStatus(status domain.CommitmentStatus) bool {
	switch status {
	case domain.CommitmentStatusPending,
		domain.CommitmentStatusPaid,
		domain.CommitmentStatusReceived,
		domain.CommitmentStatusCanceled:
		return true
	default:
		return false
	}
}

func isValidEffectiveStatus(status domain.EffectiveStatus) bool {
	switch status {
	case domain.EffectiveStatusPending,
		domain.EffectiveStatusOverdue,
		domain.EffectiveStatusPaid,
		domain.EffectiveStatusReceived,
		domain.EffectiveStatusCanceled:
		return true
	default:
		return false
	}
}
