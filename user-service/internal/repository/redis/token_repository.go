package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cinema-booking-system/user-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	userCachePrefix    = "user:"
	blacklistPrefix    = "blacklist:"
	refreshTokenPrefix = "refresh:"
	cacheTTL           = 30 * time.Minute
)

type TokenRepository struct {
	client *redis.Client
}

func New(client *redis.Client) *TokenRepository {
	return &TokenRepository{client: client}
}

func NewClient(addr, password string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: ping: %w", err)
	}
	return client, nil
}

func (r *TokenRepository) SetUserCache(ctx context.Context, user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("redis: marshal user: %w", err)
	}
	key := userCachePrefix + user.ID
	return r.client.Set(ctx, key, data, cacheTTL).Err()
}

func (r *TokenRepository) GetUserCache(ctx context.Context, id string) (*domain.User, error) {
	key := userCachePrefix + id
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("redis: get user cache: %w", err)
	}
	var user domain.User
	if err = json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("redis: unmarshal user: %w", err)
	}
	return &user, nil
}

func (r *TokenRepository) DeleteUserCache(ctx context.Context, id string) error {
	return r.client.Del(ctx, userCachePrefix+id).Err()
}

func (r *TokenRepository) SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	key := refreshTokenPrefix + token
	return r.client.Set(ctx, key, userID, ttl).Err()
}

func (r *TokenRepository) FindRefreshToken(ctx context.Context, token string) (string, error) {
	key := refreshTokenPrefix + token
	userID, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domain.ErrInvalidToken
		}
		return "", fmt.Errorf("redis: find refresh token: %w", err)
	}
	return userID, nil
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return r.client.Del(ctx, refreshTokenPrefix+token).Err()
}

func (r *TokenRepository) DeleteAllUserTokens(ctx context.Context, userID string) error {
	var cursor uint64
	pattern := refreshTokenPrefix + "*"
	for {
		keys, next, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("redis: scan refresh tokens: %w", err)
		}
		for _, key := range keys {
			val, err := r.client.Get(ctx, key).Result()
			if err != nil {
				continue
			}
			if val == userID {
				r.client.Del(ctx, key)
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (r *TokenRepository) AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := blacklistPrefix + token
	return r.client.Set(ctx, key, "1", ttl).Err()
}

func (r *TokenRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := blacklistPrefix + token
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis: check blacklist: %w", err)
	}
	return exists > 0, nil
}
