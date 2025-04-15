package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/codepnw/gopher-social/internal/auth"
	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/codepnw/gopher-social/internal/handler"
	"github.com/codepnw/gopher-social/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type middleware struct {
	auth  auth.JWTAuthenticator
	store store.Storage
}

func InitMiddleware() *middleware {
	return &middleware{}
}

func (m *middleware) AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization header is missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization header is invalid"})
			return
		}

		token, err := m.auth.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, _ := token.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		user, err := store.GetUserRepo().GetByID(c, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func (m *middleware) CheckPostOwnership(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := handler.GetUserFromContext(c)
		post := handler.GetPostFromContext(c)

		// check user post
		if post.UserID == user.ID {
			c.Next()
			return
		}

		// role precedence check
		allowed, err := m.checkRolePrecedence(c, user, role)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

func (m *middleware) checkRolePrecedence(ctx context.Context, user *entity.User, roleName string) (bool, error) {
	role, err := store.GetRoleRepo().GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
