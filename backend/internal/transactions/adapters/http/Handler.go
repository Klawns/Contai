package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	accountdomain "contai/internal/account/domain"
	authhttp "contai/internal/auth/adapters/http"
	categorydomain "contai/internal/category/domain"
	"contai/internal/transactions/app/ports"
	"contai/internal/transactions/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	transactionService ports.TransactionService
}

func NewHandler(transactionService ports.TransactionService) Handler {
	return Handler{transactionService: transactionService}
}

func (handler Handler) ListTransactions(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	input, err := parseListInput(ctx)
	if err != nil {
		writeError(ctx, err)
		return
	}
	input.UserID = authenticatedUser.UserID

	transactions, err := handler.transactionService.ListTransactions(ctx.Request.Context(), input)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toTransactionResponses(transactions))
}

func (handler Handler) CreateIncome(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request createTransactionRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	occurredAt, err := parseOccurredAt(request.OccurredAt)
	if err != nil {
		writeError(ctx, err)
		return
	}
	settlementStatus, err := parseRequiredSettlementStatus(request.SettlementStatus)
	if err != nil {
		writeError(ctx, err)
		return
	}
	settledAt, err := parseOptionalTime(request.SettledAt)
	if err != nil {
		writeError(ctx, err)
		return
	}
	recurrenceType, recurrence, err := parseRecurrencePayload(request.RecurrenceType, request.Recurrence)
	if err != nil {
		writeError(ctx, err)
		return
	}

	transaction, err := handler.transactionService.CreateIncome(ctx.Request.Context(), ports.CreateIncomeInput{
		UserID:           authenticatedUser.UserID,
		Description:      request.Description,
		Amount:           moneyFromCents(request.Amount),
		OccurredAt:       occurredAt,
		AccountID:        parseOptionalAccountID(request.AccountID),
		CategoryID:       categorydomain.CategoryID(request.CategoryID),
		SettlementStatus: settlementStatus,
		SettledAt:        settledAt,
		RecurrenceType:   recurrenceType,
		Recurrence:       recurrence,
		Note:             request.Note,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toTransactionResponse(transaction))
}

func (handler Handler) CreateExpense(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request createTransactionRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	occurredAt, err := parseOccurredAt(request.OccurredAt)
	if err != nil {
		writeError(ctx, err)
		return
	}
	settlementStatus, err := parseRequiredSettlementStatus(request.SettlementStatus)
	if err != nil {
		writeError(ctx, err)
		return
	}
	settledAt, err := parseOptionalTime(request.SettledAt)
	if err != nil {
		writeError(ctx, err)
		return
	}
	recurrenceType, recurrence, err := parseRecurrencePayload(request.RecurrenceType, request.Recurrence)
	if err != nil {
		writeError(ctx, err)
		return
	}

	transaction, err := handler.transactionService.CreateExpense(ctx.Request.Context(), ports.CreateExpenseInput{
		UserID:           authenticatedUser.UserID,
		Description:      request.Description,
		Amount:           moneyFromCents(request.Amount),
		OccurredAt:       occurredAt,
		AccountID:        parseOptionalAccountID(request.AccountID),
		CategoryID:       categorydomain.CategoryID(request.CategoryID),
		SettlementStatus: settlementStatus,
		SettledAt:        settledAt,
		RecurrenceType:   recurrenceType,
		Recurrence:       recurrence,
		Note:             request.Note,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toTransactionResponse(transaction))
}

func (handler Handler) CreateTransfer(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request createTransferRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	occurredAt, err := parseOccurredAt(request.OccurredAt)
	if err != nil {
		writeError(ctx, err)
		return
	}

	transaction, err := handler.transactionService.CreateTransfer(ctx.Request.Context(), ports.CreateTransferInput{
		UserID:               authenticatedUser.UserID,
		Description:          request.Description,
		Amount:               moneyFromCents(request.Amount),
		OccurredAt:           occurredAt,
		SourceAccountID:      accountdomain.AccountID(request.SourceAccountID),
		DestinationAccountID: accountdomain.AccountID(request.DestinationAccountID),
		Note:                 request.Note,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toTransactionResponse(transaction))
}

func (handler Handler) UpdateTransaction(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request updateTransactionRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	occurredAt, err := parseOccurredAt(request.OccurredAt)
	if err != nil {
		writeError(ctx, err)
		return
	}
	var settlementStatus domain.SettlementStatus
	if request.SettlementStatus != nil {
		parsedStatus, err := parseRequiredSettlementStatus(request.SettlementStatus)
		if err != nil {
			writeError(ctx, err)
			return
		}
		settlementStatus = parsedStatus
	}
	settledAt, err := parseOptionalTime(request.SettledAt)
	if err != nil {
		writeError(ctx, err)
		return
	}
	recurrenceType, recurrence, err := parseRecurrencePayload(request.RecurrenceType, request.Recurrence)
	if err != nil {
		writeError(ctx, err)
		return
	}

	transaction, err := handler.transactionService.UpdateTransaction(ctx.Request.Context(), ports.UpdateTransactionInput{
		UserID:               authenticatedUser.UserID,
		TransactionID:        domain.TransactionID(ctx.Param("transactionID")),
		Description:          request.Description,
		Amount:               moneyFromCents(request.Amount),
		OccurredAt:           occurredAt,
		AccountID:            parseOptionalAccountID(request.AccountID),
		SourceAccountID:      accountdomain.AccountID(request.SourceAccountID),
		DestinationAccountID: accountdomain.AccountID(request.DestinationAccountID),
		CategoryID:           categorydomain.CategoryID(request.CategoryID),
		SettlementStatus:     settlementStatus,
		SettledAt:            settledAt,
		RecurrenceType:       recurrenceType,
		Recurrence:           recurrence,
		Note:                 request.Note,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toTransactionResponse(transaction))
}

func (handler Handler) DeleteTransaction(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := handler.transactionService.DeleteTransaction(ctx.Request.Context(), ports.DeleteTransactionInput{
		UserID:        authenticatedUser.UserID,
		TransactionID: domain.TransactionID(ctx.Param("transactionID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func parseListInput(ctx *gin.Context) (ports.ListTransactionsInput, error) {
	var input ports.ListTransactionsInput
	if value := ctx.Query("startAt"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return input, domain.ErrTransactionOccurredAtRequired
		}
		input.StartAt = &parsed
	}
	if value := ctx.Query("endAt"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return input, domain.ErrTransactionOccurredAtRequired
		}
		input.EndAt = &parsed
	}
	if value := ctx.Query("accountId"); value != "" {
		if value == "none" {
			input.AccountIDNone = true
		} else {
			accountID := accountdomain.AccountID(value)
			input.AccountID = &accountID
		}
	}
	if value := ctx.Query("categoryId"); value != "" {
		categoryID := categorydomain.CategoryID(value)
		input.CategoryID = &categoryID
	}
	transactionType, err := parseTransactionType(ctx.Query("type"))
	if err != nil {
		return input, err
	}
	input.Type = transactionType
	settlementStatus, err := parseSettlementStatus(ctx.Query("settlementStatus"))
	if err != nil {
		return input, err
	}
	input.SettlementStatus = settlementStatus
	limit, err := parseOptionalNonNegativeInt(ctx.Query("limit"))
	if err != nil {
		return input, err
	}
	offset, err := parseOptionalNonNegativeInt(ctx.Query("offset"))
	if err != nil {
		return input, err
	}
	input.Limit = limit
	input.Offset = offset
	return input, nil
}

func parseOptionalNonNegativeInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, domain.ErrTransactionInvalidPagination
	}
	return parsed, nil
}

func bindJSON(ctx *gin.Context, request any) error {
	if err := ctx.ShouldBindJSON(request); err != nil {
		if isBodyTooLarge(err) {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return err
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return err
	}
	return nil
}

func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrTransactionNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
	case errors.Is(err, domain.ErrTransactionAccountNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
	case errors.Is(err, domain.ErrTransactionCategoryNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
	case errors.Is(err, domain.ErrTransactionRemoved):
		ctx.JSON(http.StatusConflict, gin.H{"error": "transaction is removed"})
	case errors.Is(err, domain.ErrTransactionManagedOrigin):
		ctx.JSON(http.StatusConflict, gin.H{"error": "transaction has managed origin"})
	case errors.Is(err, domain.ErrTransactionCategoryTypeMismatch):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "category type mismatch"})
	case errors.Is(err, domain.ErrTransactionIDRequired),
		errors.Is(err, domain.ErrTransactionUserIDRequired),
		errors.Is(err, domain.ErrTransactionDescriptionRequired),
		errors.Is(err, domain.ErrTransactionAmountInvalid),
		errors.Is(err, domain.ErrTransactionOccurredAtRequired),
		errors.Is(err, domain.ErrTransactionAccountIDRequired),
		errors.Is(err, domain.ErrTransactionSourceAccountIDRequired),
		errors.Is(err, domain.ErrTransactionDestinationAccountIDRequired),
		errors.Is(err, domain.ErrTransactionTransferAccountsMustBeDifferent),
		errors.Is(err, domain.ErrTransactionCategoryIDRequired),
		errors.Is(err, domain.ErrTransactionInvalidType),
		errors.Is(err, domain.ErrTransactionInvalidStatus),
		errors.Is(err, domain.ErrTransactionSettlementStatusRequired),
		errors.Is(err, domain.ErrTransactionInvalidSettlementStatus),
		errors.Is(err, domain.ErrTransactionInvalidSettledAt),
		errors.Is(err, domain.ErrTransactionInvalidRecurrence),
		errors.Is(err, domain.ErrTransactionInvalidPagination):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isBodyTooLarge(err error) bool {
	return strings.Contains(err.Error(), "request body too large")
}
