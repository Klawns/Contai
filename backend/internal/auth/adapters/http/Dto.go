package http

import userports "contai/internal/users/app/ports"

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type authenticatedUserResponse struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func toUserResponse(user userports.UserDTO) userResponse {
	return userResponse{
		ID:        string(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Status:    string(user.Status),
		CreatedAt: user.CreatedAt.Format(timeFormatRFC3339),
	}
}
