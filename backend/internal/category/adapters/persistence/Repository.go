package persistence

import (
	"context"
	"errors"

	"contai/internal/category/app/ports"
	"contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var _ ports.CategoryRepository = CategoryRepository{}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return CategoryRepository{db: db}
}

func (repository CategoryRepository) WithTx(tx databaseports.TxHandle) ports.CategoryRepository {
	if db, ok := tx.Value().(*gorm.DB); ok && db != nil {
		return CategoryRepository{db: db}
	}

	return repository
}

func (repository CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	entity := toCategoryEntity(*category)
	if err := repository.db.WithContext(ctx).Create(&entity).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrCategoryNameAlreadyExists
		}
		return nil, err
	}

	created, err := toDomainCategory(entity)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (repository CategoryRepository) CreateCategories(ctx context.Context, categories []domain.Category) ([]domain.Category, error) {
	if len(categories) == 0 {
		return []domain.Category{}, nil
	}

	entities := make([]CategoryEntity, 0, len(categories))
	for _, category := range categories {
		entities = append(entities, toCategoryEntity(category))
	}

	if err := repository.db.WithContext(ctx).Create(&entities).Error; err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrCategoryNameAlreadyExists
		}
		return nil, err
	}

	created := make([]domain.Category, 0, len(entities))
	for _, entity := range entities {
		category, err := toDomainCategory(entity)
		if err != nil {
			return nil, err
		}
		created = append(created, category)
	}

	return created, nil
}

func (repository CategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	entity := toCategoryEntity(*category)
	result := repository.db.WithContext(ctx).Save(&entity)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return nil, domain.ErrCategoryNameAlreadyExists
		}
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, domain.ErrCategoryNotFound
	}

	updated, err := toDomainCategory(entity)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (repository CategoryRepository) FindCategoryByID(ctx context.Context, categoryID domain.CategoryID, userID userdomain.UserID) (*domain.Category, error) {
	var entity CategoryEntity
	err := repository.db.WithContext(ctx).First(&entity, "id = ? AND user_id = ?", string(categoryID), string(userID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	category, err := toDomainCategory(entity)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (repository CategoryRepository) FindCategoriesByUserID(ctx context.Context, input ports.ListCategoriesInput) ([]domain.Category, error) {
	query := repository.db.WithContext(ctx).Where("user_id = ?", string(input.UserID))
	if input.Type != nil {
		query = query.Where("type = ?", string(*input.Type))
	}
	if input.Status != nil {
		query = query.Where("status = ?", string(*input.Status))
	}

	var entities []CategoryEntity
	if err := query.Order("type ASC, is_default DESC, name ASC").Find(&entities).Error; err != nil {
		return nil, err
	}

	categories := make([]domain.Category, 0, len(entities))
	for _, entity := range entities {
		category, err := toDomainCategory(entity)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (repository CategoryRepository) CategoryNameExistsByUserAndType(ctx context.Context, userID userdomain.UserID, categoryType domain.CategoryType, normalizedName string, excludingCategoryID *domain.CategoryID) (bool, error) {
	query := repository.db.WithContext(ctx).
		Model(&CategoryEntity{}).
		Where("user_id = ? AND type = ? AND normalized_name = ?", string(userID), string(categoryType), normalizedName)
	if excludingCategoryID != nil {
		query = query.Where("id <> ?", string(*excludingCategoryID))
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func isUniqueViolation(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
