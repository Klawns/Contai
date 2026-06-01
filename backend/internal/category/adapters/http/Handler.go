package http

import (
	"errors"
	"net/http"
	"strings"

	authhttp "contai/internal/auth/adapters/http"
	"contai/internal/category/app/ports"
	"contai/internal/category/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	categoryService ports.CategoryService
}

func NewHandler(categoryService ports.CategoryService) Handler {
	return Handler{categoryService: categoryService}
}

func (handler Handler) ListCategories(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categoryType, err := parseCategoryType(ctx.Query("type"))
	if err != nil {
		writeError(ctx, err)
		return
	}
	status, err := parseCategoryStatus(ctx.Query("status"))
	if err != nil {
		writeError(ctx, err)
		return
	}

	categories, err := handler.categoryService.ListCategories(ctx.Request.Context(), ports.ListCategoriesInput{
		UserID: authenticatedUser.UserID,
		Type:   categoryType,
		Status: status,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toCategoryResponses(categories))
}

func (handler Handler) CreateCategory(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request createCategoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		if isBodyTooLarge(err) {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	category, err := handler.categoryService.CreateCategory(ctx.Request.Context(), ports.CreateCategoryInput{
		UserID: authenticatedUser.UserID,
		Name:   request.Name,
		Type:   domain.CategoryType(request.Type),
		Color:  request.Color,
		Icon:   request.Icon,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toCategoryResponse(category))
}

func (handler Handler) UpdateCategory(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var request updateCategoryRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		if isBodyTooLarge(err) {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	category, err := handler.categoryService.UpdateCategory(ctx.Request.Context(), ports.UpdateCategoryInput{
		UserID:     authenticatedUser.UserID,
		CategoryID: domain.CategoryID(ctx.Param("categoryID")),
		Name:       request.Name,
		Color:      request.Color,
		Icon:       request.Icon,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toCategoryResponse(category))
}

func (handler Handler) DeleteCategory(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	err := handler.categoryService.InactivateCategory(ctx.Request.Context(), ports.InactivateCategoryInput{
		UserID:     authenticatedUser.UserID,
		CategoryID: domain.CategoryID(ctx.Param("categoryID")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCategoryNameAlreadyExists):
		ctx.JSON(http.StatusConflict, gin.H{"error": "category name already exists"})
	case errors.Is(err, domain.ErrCategoryNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
	case errors.Is(err, domain.ErrCategoryIDRequired),
		errors.Is(err, domain.ErrCategoryUserIDRequired),
		errors.Is(err, domain.ErrCategoryNameRequired),
		errors.Is(err, domain.ErrCategoryInvalidType),
		errors.Is(err, domain.ErrCategoryInvalidStatus),
		errors.Is(err, domain.ErrCategoryInvalidColor),
		errors.Is(err, domain.ErrCategoryInvalidIcon):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isBodyTooLarge(err error) bool {
	return strings.Contains(err.Error(), "request body too large")
}
