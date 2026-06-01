package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authdomain "contai/internal/auth/domain"
	"contai/internal/category/app/ports"
	"contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"

	"github.com/gin-gonic/gin"
)

func TestHandlerRequiresAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(&fakeCategoryService{}).RegisterForTest(router)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/categories", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestHandlerCreateCategoryUsesAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeCategoryService{}
	router := authenticatedRouter(service)
	body := bytes.NewBufferString(`{"name":"Alimentação","type":"expense","color":"#EA580C","icon":"utensils"}`)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/categories", body)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.createInput.UserID != "authenticated-user" {
		t.Fatalf("expected authenticated user id, got %s", service.createInput.UserID)
	}
}

func TestHandlerMapsDuplicateNameToConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeCategoryService{err: domain.ErrCategoryNameAlreadyExists}
	router := authenticatedRouter(service)
	body := bytes.NewBufferString(`{"name":"Alimentação","type":"expense","color":"#EA580C","icon":"utensils"}`)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/categories", body)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", recorder.Code)
	}
}

func TestHandlerDeleteCategoryReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeCategoryService{}
	router := authenticatedRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/categories/category-id", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}
	if service.inactivateInput.CategoryID != "category-id" {
		t.Fatalf("expected category id from path, got %s", service.inactivateInput.CategoryID)
	}
}

func authenticatedRouter(service *fakeCategoryService) *gin.Engine {
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("authenticated_user", authdomain.AuthenticatedUser{UserID: "authenticated-user"})
		ctx.Next()
	})
	NewHandler(service).RegisterForTest(router)
	return router
}

func (handler Handler) RegisterForTest(router *gin.Engine) {
	router.GET("/categories", handler.ListCategories)
	router.POST("/categories", handler.CreateCategory)
	router.PATCH("/categories/:categoryID", handler.UpdateCategory)
	router.DELETE("/categories/:categoryID", handler.DeleteCategory)
}

type fakeCategoryService struct {
	err             error
	createInput     ports.CreateCategoryInput
	inactivateInput ports.InactivateCategoryInput
}

func (service *fakeCategoryService) WithTx(tx databaseports.TxHandle) ports.CategoryService {
	return service
}

func (service *fakeCategoryService) CreateCategory(ctx context.Context, input ports.CreateCategoryInput) (ports.CategoryDTO, error) {
	service.createInput = input
	if service.err != nil {
		return ports.CategoryDTO{}, service.err
	}
	return fakeCategoryDTO(input.UserID), nil
}

func (service *fakeCategoryService) CreateDefaultCategories(ctx context.Context, userID userdomain.UserID) error {
	return service.err
}

func (service *fakeCategoryService) ListCategories(ctx context.Context, input ports.ListCategoriesInput) ([]ports.CategoryDTO, error) {
	if service.err != nil {
		return nil, service.err
	}
	return []ports.CategoryDTO{fakeCategoryDTO(input.UserID)}, nil
}

func (service *fakeCategoryService) UpdateCategory(ctx context.Context, input ports.UpdateCategoryInput) (ports.CategoryDTO, error) {
	if service.err != nil {
		return ports.CategoryDTO{}, service.err
	}
	return fakeCategoryDTO(input.UserID), nil
}

func (service *fakeCategoryService) InactivateCategory(ctx context.Context, input ports.InactivateCategoryInput) error {
	service.inactivateInput = input
	if errors.Is(service.err, domain.ErrCategoryNotFound) {
		return service.err
	}
	return service.err
}

func fakeCategoryDTO(userID userdomain.UserID) ports.CategoryDTO {
	now := time.Now()
	return ports.CategoryDTO{
		ID:             "category-id",
		UserID:         userID,
		Name:           "Alimentação",
		NormalizedName: "alimentacao",
		Type:           domain.CategoryTypeExpense,
		Color:          "#EA580C",
		Icon:           "utensils",
		Status:         domain.CategoryStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
