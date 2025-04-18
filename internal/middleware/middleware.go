package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/codepnw/gopher-social/internal/auth"
	"github.com/codepnw/gopher-social/internal/entity"
	"github.com/codepnw/gopher-social/internal/handler"
	"github.com/codepnw/gopher-social/internal/store"
	"github.com/codepnw/gopher-social/internal/store/cache"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type middleware struct {
	auth auth.JWTAuthenticator
	// store store.Storage
	redis cache.Storage
}

func InitMiddleware(redis cache.Storage) *middleware {
	return &middleware{
		redis: redis,
	}
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

		// // get user from cache
		// user, err := m.redis.Users.Get(c, userID)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// 	return
		// }

		// if user == nil {
		// 	// get user from db
		// 	user, err = store.GetUserRepo().GetByID(c, userID)
		// 	if err != nil {
		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// 		return
		// 	}

		// 	// set user
		// 	if err := m.redis.Users.Set(c, user); err != nil {
		// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// 		return
		// 	}
		// }

		user, err := m.getUser(c, userID)
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

func (m *middleware) getUser(ctx context.Context, userID int64) (*entity.User, error) {
	// get user from cache
	user, err := m.redis.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	log.Println("user 1", user)

	if user == nil {
		// get user from db
		user, err = store.GetUserRepo().GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		// set user
		if err := m.redis.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}
