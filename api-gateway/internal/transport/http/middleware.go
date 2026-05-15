package httpgateway

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/cinema-booking-system/api-gateway/internal/config"
	"github.com/cinema-booking-system/api-gateway/internal/domain"
	"github.com/cinema-booking-system/api-gateway/internal/repository"
	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const authUserKey = "auth_user"

type RateLimiter struct {
	mu       sync.Mutex
	limit    rate.Limit
	burst    int
	limiters map[string]*rate.Limiter
}

func NewRateLimiter(eventsPerSecond float64, burst int) *RateLimiter {
	if eventsPerSecond <= 0 {
		eventsPerSecond = 1
	}
	if burst <= 0 {
		burst = 1
	}
	return &RateLimiter{
		limit:    rate.Limit(eventsPerSecond),
		burst:    burst,
		limiters: make(map[string]*rate.Limiter),
	}
}

func (l *RateLimiter) Allow(key string) bool {
	return l.limiterFor(key).Allow()
}

func (l *RateLimiter) limiterFor(key string) *rate.Limiter {
	if key == "" {
		key = "unknown"
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, ok := l.limiters[key]
	if ok {
		return limiter
	}

	limiter = rate.NewLimiter(l.limit, l.burst)
	l.limiters[key] = limiter
	return limiter
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedOrigins))
	allowAny := false
	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAny = true
			continue
		}
		allowed[origin] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			if allowAny {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			} else if allowed[origin] {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func rateLimitMiddleware(limiter *RateLimiter, keyFunc func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(keyFunc(c)) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func authMiddleware(cfg config.Config, users repository.UserClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c.GetHeader("Authorization"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.RequestTimeout)
		defer cancel()
		resp, err := users.ValidateToken(ctx, &userpb.ValidateTokenRequest{AccessToken: token})
		if err != nil {
			statusCode, message := grpcError(err)
			c.JSON(statusCode, gin.H{"error": message})
			c.Abort()
			return
		}
		if !resp.GetValid() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set(authUserKey, domain.AuthUser{
			ID:    resp.GetUserId(),
			Email: resp.GetEmail(),
			Role:  roleFromProto(resp.GetRole()),
		})
		c.Next()
	}
}

func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := currentUser(c)
		if !ok || !user.IsAdmin() {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func currentUser(c *gin.Context) (domain.AuthUser, bool) {
	raw, ok := c.Get(authUserKey)
	if !ok {
		return domain.AuthUser{}, false
	}
	user, ok := raw.(domain.AuthUser)
	return user, ok
}

func roleFromProto(role userpb.Role) domain.Role {
	if role == userpb.Role_ROLE_ADMIN {
		return domain.RoleAdmin
	}
	return domain.RoleUser
}

func bearerToken(header string) string {
	prefix := "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
