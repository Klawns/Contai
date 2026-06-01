package http

import (
	"errors"
	"net/http"
	"strings"

	"contai/internal/auth/app/contracts"
	authports "contai/internal/auth/app/ports"
	authdomain "contai/internal/auth/domain"
	userports "contai/internal/users/app/ports"
	userdomain "contai/internal/users/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService authports.AuthService
	userService userports.UserService
	cookies     CookieService
}

func NewHandler(authService authports.AuthService, userService userports.UserService, cookies CookieService) Handler {
	return Handler{
		authService: authService,
		userService: userService,
		cookies:     cookies,
	}
}

func (handler Handler) CreateUser(ctx *gin.Context) {
	var request createUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		if isBodyTooLarge(err) {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := handler.userService.CreateUser(ctx.Request.Context(), userports.CreateUserInput{
		Name:          request.Name,
		Email:         request.Email,
		PlainPassword: request.Password,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	_, tokens, err := handler.authService.Login(ctx.Request.Context(), contracts.LoginInput{
		Email:         request.Email,
		PlainPassword: request.Password,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	handler.cookies.SetAccessCookie(ctx.Writer, tokens.AccessToken, tokens.AccessClaims.ExpiresAt)
	ctx.JSON(http.StatusCreated, toUserResponse(user))
}

func (handler Handler) Login(ctx *gin.Context) {
	var request loginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		if isBodyTooLarge(err) {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	authenticatedUser, tokens, err := handler.authService.Login(ctx.Request.Context(), contracts.LoginInput{
		Email:         request.Email,
		PlainPassword: request.Password,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	handler.cookies.SetAccessCookie(ctx.Writer, tokens.AccessToken, tokens.AccessClaims.ExpiresAt)
	ctx.JSON(http.StatusOK, authenticatedUserResponse{
		ID:     string(authenticatedUser.UserID),
		Email:  authenticatedUser.Email,
		Status: string(authenticatedUser.Status),
	})
}

func (handler Handler) Logout(ctx *gin.Context) {
	if err := handler.authService.Logout(ctx.Request.Context()); err != nil {
		writeError(ctx, err)
		return
	}

	handler.cookies.ClearAccessCookie(ctx.Writer)
	ctx.Status(http.StatusNoContent)
}

func (handler Handler) Me(ctx *gin.Context) {
	authenticatedUser, ok := AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, authenticatedUserResponse{
		ID:     string(authenticatedUser.UserID),
		Email:  authenticatedUser.Email,
		Status: string(authenticatedUser.Status),
	})
}

func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, authdomain.ErrInvalidCredentials):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	case errors.Is(err, authdomain.ErrInvalidToken), errors.Is(err, authdomain.ErrExpiredToken):
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	case errors.Is(err, userdomain.ErrUserInactive):
		ctx.JSON(http.StatusForbidden, gin.H{"error": "user inactive"})
	case errors.Is(err, userdomain.ErrUserEmailAlreadyExists):
		ctx.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
	case errors.Is(err, userdomain.ErrUserPasswordTooWeak):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "weak password"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isBodyTooLarge(err error) bool {
	return strings.Contains(err.Error(), "request body too large")
}
