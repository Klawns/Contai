package domain

import (
	"errors"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("should create a valid user", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.Name != "John Doe" {
			t.Errorf("expected name John Doe, got %s", user.Name)
		}

		if user.Email != "john@example.com" {
			t.Errorf("expected email john@example.com, got %s", user.Email)
		}

		if user.PasswordHash != "hashed-password" {
			t.Errorf("expected password hash hashed-password, got %s", user.PasswordHash)
		}

		if user.Status != UserStatusActive {
			t.Errorf("expected status active, got %s", user.Status)
		}

		if user.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}

		if user.UpdatedAt.IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		_, err := NewUser("", "john@example.com", "hashed-password")

		if !errors.Is(err, ErrUserNameRequired) {
			t.Fatalf("expected ErrUserNameRequired, got %v", err)
		}
	})

	t.Run("should return error when email is empty", func(t *testing.T) {
		_, err := NewUser("John Doe", "", "hashed-password")

		if !errors.Is(err, ErrUserEmailRequired) {
			t.Fatalf("expected ErrUserEmailRequired, got %v", err)
		}
	})

	t.Run("should return error when email is invalid", func(t *testing.T) {
		_, err := NewUser("John Doe", "invalid-email", "hashed-password")

		if !errors.Is(err, ErrUserInvalidEmail) {
			t.Fatalf("expected ErrUserInvalidEmail, got %v", err)
		}
	})

	t.Run("should return error when password hash is empty", func(t *testing.T) {
		_, err := NewUser("John Doe", "john@example.com", "")

		if !errors.Is(err, ErrUserPasswordHashRequired) {
			t.Fatalf("expected ErrUserPasswordHashRequired, got %v", err)
		}
	})
}

func TestUsuario_Rename(t *testing.T) {
	t.Run("should rename user", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Rename("Jane Doe")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.Name != "Jane Doe" {
			t.Errorf("expected name Jane Doe, got %s", user.Name)
		}
	})

	t.Run("should return error when name is empty", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Rename("")

		if !errors.Is(err, ErrUserNameRequired) {
			t.Fatalf("expected ErrUserNameRequired, got %v", err)
		}
	})
}

func TestUsuario_ChangePasswordHash(t *testing.T) {
	t.Run("should change password hash", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.ChangePasswordHash("new-hashed-password")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.PasswordHash != "new-hashed-password" {
			t.Errorf("expected new password hash, got %s", user.PasswordHash)
		}
	})

	t.Run("should return error when password hash is empty", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.ChangePasswordHash("")

		if !errors.Is(err, ErrUserPasswordHashRequired) {
			t.Fatalf("expected ErrUserPasswordHashRequired, got %v", err)
		}
	})
}

func TestUsuario_Deactivate(t *testing.T) {
	t.Run("should deactivate user", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Deactivate()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.Status != UserStatusInactive {
			t.Errorf("expected status inactive, got %s", user.Status)
		}
	})

	t.Run("should return error when user is already inactive", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Deactivate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Deactivate()

		if !errors.Is(err, ErrUserAlreadyInactive) {
			t.Fatalf("expected ErrUserAlreadyInactive, got %v", err)
		}
	})
}

func TestUsuario_Activate(t *testing.T) {
	t.Run("should activate inactive user", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Deactivate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Activate()

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if user.Status != UserStatusActive {
			t.Errorf("expected status active, got %s", user.Status)
		}
	})

	t.Run("should return error when user is already active", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Activate()

		if !errors.Is(err, ErrUserAlreadyActive) {
			t.Fatalf("expected ErrUserAlreadyActive, got %v", err)
		}
	})
}

func TestUsuario_CanAuthenticate(t *testing.T) {
	t.Run("should return true when user is active", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		canAuthenticate, err := user.CanAuthenticate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !canAuthenticate {
			t.Error("expected user can authenticate")
		}
	})

	t.Run("should return false when user is inactive", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Deactivate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		canAuthenticate, err := user.CanAuthenticate()
		if !errors.Is(err, ErrUserInactive) {
			t.Fatalf("expected ErrUserInactive, got %v", err)
		}

		if canAuthenticate {
			t.Error("expected user cannot authenticate")
		}
	})
}

func TestUsuario_CanCreateFinancialData(t *testing.T) {
	t.Run("should return true when user is active", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		canCreate, err := user.CanCreateFinancialData()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !canCreate {
			t.Error("expected user can create financial data")
		}
	})

	t.Run("should return false when user is inactive", func(t *testing.T) {
		user, err := NewUser("John Doe", "john@example.com", "hashed-password")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = user.Deactivate()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		canCreate, err := user.CanCreateFinancialData()
		if !errors.Is(err, ErrUserInactive) {
			t.Fatalf("expected ErrUserInactive, got %v", err)
		}

		if canCreate {
			t.Error("expected user cannot create financial data")
		}
	})
}
