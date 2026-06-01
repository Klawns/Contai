package domain

import (
	"strings"
	"time"
)

type UserID string

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	ID           UserID
	Name         string
	Email        string
	PasswordHash string
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(id UserID, name, email, passwordHash string) (User, error) {
	now := time.Now()

	user := User{
		ID:           UserID(strings.TrimSpace(string(id))),
		Name:         strings.TrimSpace(name),
		Email:        normalizeEmail(email),
		PasswordHash: strings.TrimSpace(passwordHash),
		Status:       UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := user.validate(); err != nil {
		return User{}, err
	}

	return user, nil
}

func RehydrateUser(id UserID, name, email, passwordHash string, status UserStatus, createdAt, updatedAt time.Time) (User, error) {
	user := User{
		ID:           UserID(strings.TrimSpace(string(id))),
		Name:         strings.TrimSpace(name),
		Email:        normalizeEmail(email),
		PasswordHash: strings.TrimSpace(passwordHash),
		Status:       status,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	if err := user.validate(); err != nil {
		return User{}, err
	}

	return user, nil
}

func (u *User) Rename(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrUserNameRequired
	}

	u.Name = name
	u.UpdatedAt = time.Now()

	return nil
}

func (u *User) ChangeEmail(email string) error {
	email = normalizeEmail(email)

	if email == "" {
		return ErrUserEmailRequired
	}

	if !strings.Contains(email, "@") {
		return ErrUserInvalidEmail
	}

	u.Email = email
	u.UpdatedAt = time.Now()

	return nil
}

func NormalizeEmail(email string) string {
	return normalizeEmail(email)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func (u *User) ChangePasswordHash(passwordHash string) error {
	passwordHash = strings.TrimSpace(passwordHash)

	if passwordHash == "" {
		return ErrUserPasswordHashRequired
	}

	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()

	return nil
}

func (u *User) Activate() error {
	if u.Status == UserStatusActive {
		return ErrUserAlreadyActive
	}
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Deactivate() error {
	if u.Status == UserStatusInactive {
		return ErrUserAlreadyInactive
	}
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
	return nil
}

func (u User) CanAuthenticate() (bool, error) {
	if u.Status == UserStatusInactive {
		return false, ErrUserInactive
	}
	return true, nil
}

func (u User) CanCreateFinancialData() (bool, error) {
	if u.Status == UserStatusInactive {
		return false, ErrUserInactive
	}
	return true, nil
}

func (u User) validate() error {
	if strings.TrimSpace(string(u.ID)) == "" {
		return ErrUserIDRequired
	}

	if strings.TrimSpace(u.Name) == "" {
		return ErrUserNameRequired
	}

	if strings.TrimSpace(u.Email) == "" {
		return ErrUserEmailRequired
	}

	if !strings.Contains(u.Email, "@") {
		return ErrUserInvalidEmail
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		return ErrUserPasswordHashRequired
	}

	if u.Status != UserStatusActive && u.Status != UserStatusInactive {
		return ErrUserInvalidStatus
	}

	return nil
}
