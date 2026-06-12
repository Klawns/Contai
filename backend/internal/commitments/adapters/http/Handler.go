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
	"contai/internal/commitments/app/ports"
	"contai/internal/commitments/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	commitmentService ports.CommitmentService
}

func NewHandler(commitmentService ports.CommitmentService) Handler {
	return Handler{commitmentService: commitmentService}
}

func (handler Handler) ListPayables(ctx *gin.Context) {
	handler.list(ctx, domain.CommitmentTypePayable)
}

func (handler Handler) CreatePayable(ctx *gin.Context) {
	handler.create(ctx, domain.CommitmentTypePayable)
}

func (handler Handler) UpdatePayable(ctx *gin.Context) {
	handler.update(ctx, domain.CommitmentTypePayable)
}

func (handler Handler) PayPayable(ctx *gin.Context) {
	handler.settle(ctx, domain.CommitmentTypePayable)
}

func (handler Handler) CancelPayable(ctx *gin.Context) {
	handler.cancel(ctx, domain.CommitmentTypePayable)
}

func (handler Handler) ListReceivables(ctx *gin.Context) {
	handler.list(ctx, domain.CommitmentTypeReceivable)
}

func (handler Handler) CreateReceivable(ctx *gin.Context) {
	handler.create(ctx, domain.CommitmentTypeReceivable)
}

func (handler Handler) UpdateReceivable(ctx *gin.Context) {
	handler.update(ctx, domain.CommitmentTypeReceivable)
}

func (handler Handler) ReceiveReceivable(ctx *gin.Context) {
	handler.settle(ctx, domain.CommitmentTypeReceivable)
}

func (handler Handler) CancelReceivable(ctx *gin.Context) {
	handler.cancel(ctx, domain.CommitmentTypeReceivable)
}

func (handler Handler) list(ctx *gin.Context, commitmentType domain.CommitmentType) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	input, err := parseListInput(ctx, commitmentType)
	if err != nil {
		writeError(ctx, err)
		return
	}
	input.UserID = authenticatedUser.UserID

	commitments, err := handler.commitmentService.ListCommitments(ctx.Request.Context(), input)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCommitmentResponses(commitments))
}

func (handler Handler) create(ctx *gin.Context, commitmentType domain.CommitmentType) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request commitmentRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	input, err := requestToCreateInput(request, string(authenticatedUser.UserID), commitmentType)
	if err != nil {
		writeError(ctx, err)
		return
	}
	commitment, err := handler.commitmentService.CreateCommitment(ctx.Request.Context(), input)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toCommitmentResponse(commitment))
}

func (handler Handler) update(ctx *gin.Context, commitmentType domain.CommitmentType) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request commitmentRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	input, err := requestToUpdateInput(
		request,
		string(authenticatedUser.UserID),
		ctx.Param("id"),
		commitmentType,
	)
	if err != nil {
		writeError(ctx, err)
		return
	}
	commitment, err := handler.commitmentService.UpdateCommitment(ctx.Request.Context(), input)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCommitmentResponse(commitment))
}

func (handler Handler) settle(ctx *gin.Context, commitmentType domain.CommitmentType) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request settlementRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	input, err := requestToSettleInput(
		request,
		string(authenticatedUser.UserID),
		ctx.Param("id"),
		commitmentType,
	)
	if err != nil {
		writeError(ctx, err)
		return
	}
	commitment, err := handler.commitmentService.SettleCommitment(ctx.Request.Context(), input)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCommitmentResponse(commitment))
}

func (handler Handler) cancel(ctx *gin.Context, commitmentType domain.CommitmentType) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	commitment, err := handler.commitmentService.CancelCommitment(ctx.Request.Context(), ports.CancelCommitmentInput{
		UserID:       authenticatedUser.UserID,
		CommitmentID: domain.CommitmentID(ctx.Param("id")),
		Type:         commitmentType,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCommitmentResponse(commitment))
}

func parseListInput(ctx *gin.Context, commitmentType domain.CommitmentType) (ports.ListCommitmentsInput, error) {
	input := ports.ListCommitmentsInput{Type: commitmentType}
	if value := ctx.Query("startAt"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return input, domain.ErrCommitmentDueAtRequired
		}
		input.StartAt = &parsed
	}
	if value := ctx.Query("endAt"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return input, domain.ErrCommitmentDueAtRequired
		}
		input.EndAt = &parsed
	}
	if value := ctx.Query("status"); value != "" {
		status := domain.CommitmentStatus(value)
		input.Status = &status
	}
	if value := ctx.Query("effectiveStatus"); value != "" {
		status := domain.EffectiveStatus(value)
		input.EffectiveStatus = &status
	}
	if value := ctx.Query("accountId"); value != "" {
		accountID := accountdomain.AccountID(value)
		input.AccountID = &accountID
	}
	if value := ctx.Query("categoryId"); value != "" {
		categoryID := categorydomain.CategoryID(value)
		input.CategoryID = &categoryID
	}
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
		return 0, domain.ErrCommitmentInvalidPagination
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
	case errors.Is(err, domain.ErrCommitmentNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "commitment not found"})
	case errors.Is(err, domain.ErrCommitmentAccountNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
	case errors.Is(err, domain.ErrCommitmentCategoryNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
	case errors.Is(err, domain.ErrCommitmentNotPending),
		errors.Is(err, domain.ErrCommitmentSettlementTypeMismatch):
		ctx.JSON(http.StatusConflict, gin.H{"error": "commitment lifecycle conflict"})
	case errors.Is(err, domain.ErrCommitmentCategoryTypeMismatch),
		errors.Is(err, domain.ErrCommitmentIDRequired),
		errors.Is(err, domain.ErrCommitmentUserIDRequired),
		errors.Is(err, domain.ErrCommitmentDescriptionRequired),
		errors.Is(err, domain.ErrCommitmentAmountInvalid),
		errors.Is(err, domain.ErrCommitmentDueAtRequired),
		errors.Is(err, domain.ErrCommitmentAccountIDRequired),
		errors.Is(err, domain.ErrCommitmentCategoryIDRequired),
		errors.Is(err, domain.ErrCommitmentInvalidType),
		errors.Is(err, domain.ErrCommitmentInvalidStatus),
		errors.Is(err, domain.ErrCommitmentInvalidRecurrence),
		errors.Is(err, domain.ErrCommitmentInvalidPagination):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid commitment"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isBodyTooLarge(err error) bool {
	return strings.Contains(err.Error(), "request body too large")
}
