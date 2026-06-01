package persistence

import "contai/internal/users/domain"

func toUserEntity(user domain.User) UserEntity {
	return UserEntity{
		ID:           string(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Status:       string(user.Status),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}

func toDomainUser(entity UserEntity) (domain.User, error) {
	return domain.RehydrateUser(
		domain.UserID(entity.ID),
		entity.Name,
		entity.Email,
		entity.PasswordHash,
		domain.UserStatus(entity.Status),
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}
