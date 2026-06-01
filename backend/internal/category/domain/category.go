package domain

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	userdomain "contai/internal/users/domain"

	"golang.org/x/text/unicode/norm"
)

type CategoryID string

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type CategoryStatus string

const (
	CategoryStatusActive   CategoryStatus = "active"
	CategoryStatusInactive CategoryStatus = "inactive"
)

type Category struct {
	ID             CategoryID
	UserID         userdomain.UserID
	Name           string
	NormalizedName string
	Type           CategoryType
	Color          string
	Icon           string
	IsDefault      bool
	Status         CategoryStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

var (
	colorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	iconPattern  = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
)

func NewCategory(id CategoryID, userID userdomain.UserID, name string, categoryType CategoryType, color, icon string, isDefault bool) (Category, error) {
	now := time.Now()
	category := Category{
		ID:             CategoryID(strings.TrimSpace(string(id))),
		UserID:         userdomain.UserID(strings.TrimSpace(string(userID))),
		Name:           strings.TrimSpace(name),
		NormalizedName: NormalizeName(name),
		Type:           categoryType,
		Color:          strings.TrimSpace(color),
		Icon:           strings.TrimSpace(icon),
		IsDefault:      isDefault,
		Status:         CategoryStatusActive,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := category.validate(); err != nil {
		return Category{}, err
	}

	return category, nil
}

func RehydrateCategory(id CategoryID, userID userdomain.UserID, name, normalizedName string, categoryType CategoryType, color, icon string, isDefault bool, status CategoryStatus, createdAt, updatedAt time.Time) (Category, error) {
	category := Category{
		ID:             CategoryID(strings.TrimSpace(string(id))),
		UserID:         userdomain.UserID(strings.TrimSpace(string(userID))),
		Name:           strings.TrimSpace(name),
		NormalizedName: strings.TrimSpace(normalizedName),
		Type:           categoryType,
		Color:          strings.TrimSpace(color),
		Icon:           strings.TrimSpace(icon),
		IsDefault:      isDefault,
		Status:         status,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
	if category.NormalizedName == "" {
		category.NormalizedName = NormalizeName(category.Name)
	}

	if err := category.validate(); err != nil {
		return Category{}, err
	}

	return category, nil
}

func (category *Category) Edit(name, color, icon string) error {
	category.Name = strings.TrimSpace(name)
	category.NormalizedName = NormalizeName(name)
	category.Color = strings.TrimSpace(color)
	category.Icon = strings.TrimSpace(icon)
	category.UpdatedAt = time.Now()

	return category.validate()
}

func (category *Category) Activate() error {
	category.Status = CategoryStatusActive
	category.UpdatedAt = time.Now()
	return category.validate()
}

func (category *Category) Inactivate() error {
	category.Status = CategoryStatusInactive
	category.UpdatedAt = time.Now()
	return category.validate()
}

func NormalizeName(name string) string {
	normalized := norm.NFD.String(strings.ToLower(strings.TrimSpace(name)))
	var builder strings.Builder
	previousSpace := false

	for _, r := range normalized {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		if unicode.IsSpace(r) {
			if !previousSpace {
				builder.WriteRune(' ')
				previousSpace = true
			}
			continue
		}
		builder.WriteRune(r)
		previousSpace = false
	}

	return strings.TrimSpace(builder.String())
}

func (category Category) validate() error {
	if strings.TrimSpace(string(category.ID)) == "" {
		return ErrCategoryIDRequired
	}
	if strings.TrimSpace(string(category.UserID)) == "" {
		return ErrCategoryUserIDRequired
	}
	if strings.TrimSpace(category.Name) == "" {
		return ErrCategoryNameRequired
	}
	if strings.TrimSpace(category.NormalizedName) == "" {
		return ErrCategoryNameRequired
	}
	if category.Type != CategoryTypeIncome && category.Type != CategoryTypeExpense {
		return ErrCategoryInvalidType
	}
	if category.Status != CategoryStatusActive && category.Status != CategoryStatusInactive {
		return ErrCategoryInvalidStatus
	}
	if !colorPattern.MatchString(category.Color) {
		return ErrCategoryInvalidColor
	}
	if !iconPattern.MatchString(category.Icon) {
		return ErrCategoryInvalidIcon
	}

	return nil
}
