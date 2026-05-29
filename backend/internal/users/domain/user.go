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

type Usuario struct {
	ID           UserID
	Name         string
	Email        string
	PasswordHash string
	Status       UserStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(name, email, passwordHash string) (Usuario, error) {
	now := time.Now()

	user := Usuario{
		Name:         strings.TrimSpace(name),
		Email:        strings.TrimSpace(email),
		PasswordHash: strings.TrimSpace(passwordHash),
		Status:       UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := user.validate(); err != nil {
		return Usuario{}, err
	}

	return user, nil
}

func (u *Usuario) Rename(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return ErrUserNameRequired
	}

	u.Name = name
	u.UpdatedAt = time.Now()

	return nil
}

func (u *Usuario) ChangeEmail(email string) error {
	email = strings.TrimSpace(email)

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

func (u *Usuario) ChangePasswordHash(passwordHash string) error {
	passwordHash = strings.TrimSpace(passwordHash)

	if passwordHash == "" {
		return ErrUserPasswordHashRequired
	}

	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()

	return nil
}

func (u *Usuario) Activate() error {
	if u.Status == UserStatusActive {
		return ErrUserAlreadyActive
	}
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
	return nil
}

func (u *Usuario) Deactivate() error {
	if u.Status == UserStatusInactive {
		return ErrUserAlreadyInactive
	}
	u.Status = UserStatusInactive
	u.UpdatedAt = time.Now()
	return nil
}

func (u Usuario) CanAuthenticate() (bool, error) {
	if u.Status == UserStatusInactive {
		return false, ErrUserInactive
	}
	return true, nil
}

func (u Usuario) CanCreateFinancialData() (bool, error) {
	if u.Status == UserStatusInactive {
		return false, ErrUserInactive
	}
	return true, nil
}

func (u Usuario) validate() error {
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

	return nil
}
