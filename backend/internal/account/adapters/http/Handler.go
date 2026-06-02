package http

import (
	"errors"
	"net/http"
	"strings"

	"contai/internal/account/app/ports"
	"contai/internal/account/domain"
	authhttp "contai/internal/auth/adapters/http"
	financedomain "contai/internal/finance/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	accountService ports.AccountService
}

func NewHandler(accountService ports.AccountService) Handler {
	return Handler{accountService: accountService}
}

func (handler Handler) ListAccounts(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	status, err := parseAccountStatus(ctx.Query("status"))
	if err != nil {
		writeError(ctx, err)
		return
	}

	accounts, err := handler.accountService.ListAccounts(ctx.Request.Context(), ports.ListAccountsInput{
		UserID: authenticatedUser.UserID,
		Status: status,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toAccountResponses(accounts))
}

func (handler Handler) CreateAccount(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request createAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		if isBodyTooLarge(err) {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	account, err := handler.accountService.CreateAccount(ctx.Request.Context(), ports.CreateAccountInput{
		UserID:                  authenticatedUser.UserID,
		Name:                    request.Name,
		Type:                    domain.AccountType(request.Type),
		InitialBalance:          financedomain.NewMoney(request.InitialBalance),
		BankIconID:              request.BankIconID,
		IncludeInDashboardTotal: request.IncludeInDashboardTotal,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toAccountResponse(account))
}

func (handler Handler) GetTotalBalance(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	total, err := handler.accountService.GetTotalBalance(ctx.Request.Context(), ports.GetTotalBalanceInput{UserID: authenticatedUser.UserID})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, totalBalanceResponse{TotalBalance: total.Cents()})
}

func (handler Handler) UpdateAccount(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request updateAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		if isBodyTooLarge(err) {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	account, err := handler.accountService.UpdateAccount(ctx.Request.Context(), ports.UpdateAccountInput{
		UserID:                  authenticatedUser.UserID,
		AccountID:               domain.AccountID(ctx.Param("accountID")),
		Name:                    request.Name,
		Type:                    domain.AccountType(request.Type),
		BankIconID:              request.BankIconID,
		IncludeInDashboardTotal: request.IncludeInDashboardTotal,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toAccountResponse(account))
}

func (handler Handler) DeleteAccount(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := handler.accountService.InactivateAccount(ctx.Request.Context(), ports.InactivateAccountInput{
		UserID:    authenticatedUser.UserID,
		AccountID: domain.AccountID(ctx.Param("accountID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
	case errors.Is(err, domain.ErrAccountUserInactive):
		ctx.JSON(http.StatusForbidden, gin.H{"error": "user is inactive"})
	case errors.Is(err, domain.ErrAccountIDRequired),
		errors.Is(err, domain.ErrAccountUserIDRequired),
		errors.Is(err, domain.ErrAccountNameRequired),
		errors.Is(err, domain.ErrAccountInvalidType),
		errors.Is(err, domain.ErrAccountInvalidStatus),
		errors.Is(err, domain.ErrAccountBankIconIDRequired),
		errors.Is(err, domain.ErrAccountInvalidBankIconID),
		errors.Is(err, domain.ErrAccountMutationAmountInvalid):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid account"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isBodyTooLarge(err error) bool {
	return strings.Contains(err.Error(), "request body too large")
}
