package http

import (
	"errors"
	"net/http"
	"strings"

	authhttp "contai/internal/auth/adapters/http"
	categorydomain "contai/internal/category/domain"
	"contai/internal/creditcards/app/ports"
	"contai/internal/creditcards/domain"
	financedomain "contai/internal/finance/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service ports.CreditCardService
}

func NewHandler(service ports.CreditCardService) Handler {
	return Handler{service: service}
}

func (handler Handler) ListCreditCards(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	cards, err := handler.service.ListCreditCards(ctx.Request.Context(), authenticatedUser.UserID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCardResponses(cards))
}

func (handler Handler) CreateCreditCard(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var request cardRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	card, err := handler.service.CreateCreditCard(ctx.Request.Context(), toCreateCardInput(request, string(authenticatedUser.UserID)))
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toCardResponse(card))
}

func (handler Handler) UpdateCreditCard(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var request cardRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	card, err := handler.service.UpdateCreditCard(ctx.Request.Context(), toUpdateCardInput(request, string(authenticatedUser.UserID), ctx.Param("cardID")))
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCardResponse(card))
}

func (handler Handler) InactivateCreditCard(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	card, err := handler.service.InactivateCreditCard(ctx.Request.Context(), ports.CardIDInput{
		UserID: authenticatedUser.UserID,
		CardID: domain.CreditCardID(ctx.Param("cardID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toCardResponse(card))
}

func (handler Handler) ListPurchases(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	purchases, err := handler.service.ListPurchases(ctx.Request.Context(), ports.CardIDInput{
		UserID: authenticatedUser.UserID,
		CardID: domain.CreditCardID(ctx.Param("cardID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toPurchaseResponses(purchases))
}

func (handler Handler) CreatePurchase(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var request purchaseRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	purchaseDate, err := parseRFC3339(request.PurchaseDate)
	if err != nil {
		writeError(ctx, domain.ErrPurchaseDateRequired)
		return
	}
	firstInvoiceMonth := domain.FirstDayOfMonth(purchaseDate)
	if request.FirstInvoiceMonth != "" {
		parsed, err := parseYearMonth(request.FirstInvoiceMonth)
		if err != nil {
			writeError(ctx, domain.ErrPurchaseFirstInvoiceMonthRequired)
			return
		}
		firstInvoiceMonth = parsed
	}
	purchaseType := domain.PurchaseType(request.PurchaseType)
	if purchaseType == "" {
		if request.InstallmentCount > 1 {
			purchaseType = domain.PurchaseTypeInstallment
		} else {
			purchaseType = domain.PurchaseTypeSingle
		}
	}
	purchase, err := handler.service.CreatePurchase(ctx.Request.Context(), ports.CreatePurchaseInput{
		UserID:            authenticatedUser.UserID,
		CardID:            domain.CreditCardID(ctx.Param("cardID")),
		CategoryID:        categorydomain.CategoryID(request.CategoryID),
		Description:       request.Description,
		TotalAmount:       financedomain.NewMoney(request.TotalAmount),
		PurchaseDate:      purchaseDate,
		PurchaseType:      purchaseType,
		InstallmentCount:  request.InstallmentCount,
		FirstInvoiceMonth: firstInvoiceMonth,
		Note:              request.Note,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, toPurchaseResponse(purchase))
}

func (handler Handler) CancelPurchase(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	purchase, err := handler.service.CancelPurchase(ctx.Request.Context(), ports.PurchaseIDInput{
		UserID:     authenticatedUser.UserID,
		PurchaseID: domain.PurchaseID(ctx.Param("purchaseID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toPurchaseResponse(purchase))
}

func (handler Handler) ListInvoices(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	invoices, err := handler.service.ListInvoices(ctx.Request.Context(), ports.CardIDInput{
		UserID: authenticatedUser.UserID,
		CardID: domain.CreditCardID(ctx.Param("cardID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toInvoiceResponses(invoices))
}

func (handler Handler) GetInvoice(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	invoice, err := handler.service.GetInvoice(ctx.Request.Context(), ports.InvoiceIDInput{
		UserID:    authenticatedUser.UserID,
		InvoiceID: domain.InvoiceID(ctx.Param("invoiceID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toInvoiceResponse(invoice))
}

func (handler Handler) CloseInvoice(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	invoice, err := handler.service.CloseInvoice(ctx.Request.Context(), ports.InvoiceIDInput{
		UserID:    authenticatedUser.UserID,
		InvoiceID: domain.InvoiceID(ctx.Param("invoiceID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toInvoiceResponse(invoice))
}

func (handler Handler) PayInvoice(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var request payInvoiceRequest
	if err := bindJSON(ctx, &request); err != nil {
		return
	}
	occurredAt, err := parseRFC3339(request.OccurredAt)
	if err != nil {
		writeError(ctx, domain.ErrPurchaseDateRequired)
		return
	}
	invoice, err := handler.service.PayInvoice(ctx.Request.Context(), ports.PayInvoiceInput{
		UserID:     authenticatedUser.UserID,
		InvoiceID:  domain.InvoiceID(ctx.Param("invoiceID")),
		OccurredAt: occurredAt,
		CategoryID: categorydomain.CategoryID(request.CategoryID),
		Note:       request.Note,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toInvoiceResponse(invoice))
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
	case errors.Is(err, domain.ErrCreditCardNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "credit card not found"})
	case errors.Is(err, domain.ErrCreditCardAccountNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
	case errors.Is(err, domain.ErrCreditCardCategoryNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
	case errors.Is(err, domain.ErrPurchaseNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
	case errors.Is(err, domain.ErrInvoiceNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
	case errors.Is(err, domain.ErrCreditCardLimitExceeded),
		errors.Is(err, domain.ErrPurchaseInvoiceAlreadyPaid),
		errors.Is(err, domain.ErrInvoiceAlreadyClosed),
		errors.Is(err, domain.ErrInvoiceAlreadyPaid),
		errors.Is(err, domain.ErrInvoiceNotPayable),
		errors.Is(err, domain.ErrInvoiceAlreadyCanceled):
		ctx.JSON(http.StatusConflict, gin.H{"error": "credit card lifecycle conflict"})
	case errors.Is(err, domain.ErrCreditCardCategoryTypeMismatch):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "category type mismatch"})
	case errors.Is(err, domain.ErrCreditCardIDRequired),
		errors.Is(err, domain.ErrCreditCardUserIDRequired),
		errors.Is(err, domain.ErrCreditCardNameRequired),
		errors.Is(err, domain.ErrCreditCardAccountIDRequired),
		errors.Is(err, domain.ErrCreditCardLimitInvalid),
		errors.Is(err, domain.ErrCreditCardClosingDayInvalid),
		errors.Is(err, domain.ErrCreditCardDueDayInvalid),
		errors.Is(err, domain.ErrCreditCardInvalidStatus),
		errors.Is(err, domain.ErrPurchaseIDRequired),
		errors.Is(err, domain.ErrPurchaseDescriptionRequired),
		errors.Is(err, domain.ErrPurchaseAmountInvalid),
		errors.Is(err, domain.ErrPurchaseDateRequired),
		errors.Is(err, domain.ErrPurchaseTypeInvalid),
		errors.Is(err, domain.ErrPurchaseInstallmentCountInvalid),
		errors.Is(err, domain.ErrPurchaseFirstInvoiceMonthRequired),
		errors.Is(err, domain.ErrPurchaseInvalidStatus),
		errors.Is(err, domain.ErrInstallmentInvalid),
		errors.Is(err, domain.ErrInvoiceIDRequired),
		errors.Is(err, domain.ErrInvoiceReferenceMonthRequired),
		errors.Is(err, domain.ErrInvoiceDueAtRequired),
		errors.Is(err, domain.ErrInvoiceInvalidStatus):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid credit card"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isBodyTooLarge(err error) bool {
	return strings.Contains(err.Error(), "request body too large")
}
