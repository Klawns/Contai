package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const authenticatedUserContextKey = "authenticated_user"

func (handler Handler) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie(AccessCookieName)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		authenticatedUser, err := handler.authService.ValidateAccessToken(ctx.Request.Context(), cookie)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		ctx.Set(authenticatedUserContextKey, authenticatedUser)
		ctx.Next()
	}
}
