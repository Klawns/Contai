package services

import (
	"context"

	"contai/internal/category/app/ports"
	"contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	userports "contai/internal/users/app/ports"
	userdomain "contai/internal/users/domain"
)

var _ ports.CategoryService = CategoryService{}

type CategoryService struct {
	repository  ports.CategoryRepository
	idGenerator ports.CategoryIDGenerator
}

func NewCategoryService(repository ports.CategoryRepository, idGenerator ports.CategoryIDGenerator) CategoryService {
	return CategoryService{
		repository:  repository,
		idGenerator: idGenerator,
	}
}

func (service CategoryService) WithTx(tx databaseports.TxHandle) ports.CategoryService {
	return CategoryService{
		repository:  service.repository.WithTx(tx),
		idGenerator: service.idGenerator,
	}
}

func (service CategoryService) CreateCategory(ctx context.Context, input ports.CreateCategoryInput) (ports.CategoryDTO, error) {
	category, err := domain.NewCategory(service.idGenerator.NewCategoryID(), input.UserID, input.Name, input.Type, input.Color, input.Icon, false)
	if err != nil {
		return ports.CategoryDTO{}, err
	}

	exists, err := service.repository.CategoryNameExistsByUserAndType(ctx, category.UserID, category.Type, category.NormalizedName, nil)
	if err != nil {
		return ports.CategoryDTO{}, err
	}
	if exists {
		return ports.CategoryDTO{}, domain.ErrCategoryNameAlreadyExists
	}

	created, err := service.repository.CreateCategory(ctx, &category)
	if err != nil {
		return ports.CategoryDTO{}, err
	}

	return toCategoryDTO(*created), nil
}

func (service CategoryService) CreateDefaultCategories(ctx context.Context, userID userdomain.UserID) error {
	categories := make([]domain.Category, 0, len(defaultCategories))
	for _, preset := range defaultCategories {
		category, err := domain.NewCategory(service.idGenerator.NewCategoryID(), userID, preset.name, preset.categoryType, preset.color, preset.icon, true)
		if err != nil {
			return err
		}
		categories = append(categories, category)
	}

	_, err := service.repository.CreateCategories(ctx, categories)
	return err
}

func (service CategoryService) ListCategories(ctx context.Context, input ports.ListCategoriesInput) ([]ports.CategoryDTO, error) {
	if input.UserID == "" {
		return nil, domain.ErrCategoryUserIDRequired
	}
	if input.Type != nil && *input.Type != domain.CategoryTypeIncome && *input.Type != domain.CategoryTypeExpense {
		return nil, domain.ErrCategoryInvalidType
	}
	if input.Status != nil && *input.Status != domain.CategoryStatusActive && *input.Status != domain.CategoryStatusInactive {
		return nil, domain.ErrCategoryInvalidStatus
	}

	categories, err := service.repository.FindCategoriesByUserID(ctx, input)
	if err != nil {
		return nil, err
	}

	dtos := make([]ports.CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, toCategoryDTO(category))
	}

	return dtos, nil
}

func (service CategoryService) UpdateCategory(ctx context.Context, input ports.UpdateCategoryInput) (ports.CategoryDTO, error) {
	category, err := service.repository.FindCategoryByID(ctx, input.CategoryID, input.UserID)
	if err != nil {
		return ports.CategoryDTO{}, err
	}
	if category == nil {
		return ports.CategoryDTO{}, domain.ErrCategoryNotFound
	}

	if err := category.Edit(input.Name, input.Color, input.Icon); err != nil {
		return ports.CategoryDTO{}, err
	}

	exists, err := service.repository.CategoryNameExistsByUserAndType(ctx, category.UserID, category.Type, category.NormalizedName, &category.ID)
	if err != nil {
		return ports.CategoryDTO{}, err
	}
	if exists {
		return ports.CategoryDTO{}, domain.ErrCategoryNameAlreadyExists
	}

	updated, err := service.repository.UpdateCategory(ctx, category)
	if err != nil {
		return ports.CategoryDTO{}, err
	}

	return toCategoryDTO(*updated), nil
}

func (service CategoryService) InactivateCategory(ctx context.Context, input ports.InactivateCategoryInput) error {
	category, err := service.repository.FindCategoryByID(ctx, input.CategoryID, input.UserID)
	if err != nil {
		return err
	}
	if category == nil {
		return domain.ErrCategoryNotFound
	}
	if category.Status == domain.CategoryStatusInactive {
		return nil
	}

	if err := category.Inactivate(); err != nil {
		return err
	}

	_, err = service.repository.UpdateCategory(ctx, category)
	return err
}

type DefaultCategoryCreatorAdapter struct {
	service ports.CategoryService
}

var _ userports.DefaultCategoryCreator = DefaultCategoryCreatorAdapter{}

func NewDefaultCategoryCreatorAdapter(service ports.CategoryService) DefaultCategoryCreatorAdapter {
	return DefaultCategoryCreatorAdapter{service: service}
}

func (adapter DefaultCategoryCreatorAdapter) WithTx(tx databaseports.TxHandle) userports.DefaultCategoryCreator {
	return DefaultCategoryCreatorAdapter{service: adapter.service.WithTx(tx)}
}

func (adapter DefaultCategoryCreatorAdapter) EnsureDefaultCategories(ctx context.Context, userID userdomain.UserID) error {
	return adapter.service.CreateDefaultCategories(ctx, userID)
}

func toCategoryDTO(category domain.Category) ports.CategoryDTO {
	return ports.CategoryDTO{
		ID:             category.ID,
		UserID:         category.UserID,
		Name:           category.Name,
		NormalizedName: category.NormalizedName,
		Type:           category.Type,
		Color:          category.Color,
		Icon:           category.Icon,
		IsDefault:      category.IsDefault,
		Status:         category.Status,
		CreatedAt:      category.CreatedAt,
		UpdatedAt:      category.UpdatedAt,
	}
}

type defaultCategoryPreset struct {
	name         string
	categoryType domain.CategoryType
	color        string
	icon         string
}

var defaultCategories = []defaultCategoryPreset{
	{name: "Salário", categoryType: domain.CategoryTypeIncome, color: "#16A34A", icon: "briefcase-business"},
	{name: "Freelance", categoryType: domain.CategoryTypeIncome, color: "#0EA5E9", icon: "laptop"},
	{name: "Investimentos", categoryType: domain.CategoryTypeIncome, color: "#7C3AED", icon: "trending-up"},
	{name: "Outras receitas", categoryType: domain.CategoryTypeIncome, color: "#059669", icon: "circle-plus"},
	{name: "Moradia", categoryType: domain.CategoryTypeExpense, color: "#2563EB", icon: "house"},
	{name: "Alimentação", categoryType: domain.CategoryTypeExpense, color: "#EA580C", icon: "utensils"},
	{name: "Transporte", categoryType: domain.CategoryTypeExpense, color: "#0891B2", icon: "car"},
	{name: "Saúde", categoryType: domain.CategoryTypeExpense, color: "#DC2626", icon: "heart-pulse"},
	{name: "Educação", categoryType: domain.CategoryTypeExpense, color: "#9333EA", icon: "graduation-cap"},
	{name: "Lazer", categoryType: domain.CategoryTypeExpense, color: "#DB2777", icon: "party-popper"},
	{name: "Compras", categoryType: domain.CategoryTypeExpense, color: "#CA8A04", icon: "shopping-bag"},
	{name: "Outras despesas", categoryType: domain.CategoryTypeExpense, color: "#64748B", icon: "ellipsis"},
}
