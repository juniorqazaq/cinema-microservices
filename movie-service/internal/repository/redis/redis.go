package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/cinema-booking-system/movie-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	movieKeyPrefix          = "movie:"
	sessionsByDatePrefix    = "sessions:date:"
	sessionKeyPrefix        = "session:"
	movieCacheTTL           = time.Hour
	sessionsByDateTTL       = 15 * time.Minute
	sessionByIDTTL          = 15 * time.Minute
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
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

func (c *Cache) SetMovieCache(ctx context.Context, movie *domain.Movie) error {
	data, err := json.Marshal(movie)
	if err != nil {
		return fmt.Errorf("redis: marshal movie: %w", err)
	}
	key := movieKeyPrefix + movie.ID
	return c.client.Set(ctx, key, data, movieCacheTTL).Err()
}

func (c *Cache) GetMovieCache(ctx context.Context, id string) (*domain.Movie, error) {
	key := movieKeyPrefix + id
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis: get movie: %w", err)
	}
	var m domain.Movie
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("redis: unmarshal movie: %w", err)
	}
	return &m, nil
}

func (c *Cache) DeleteMovieCache(ctx context.Context, id string) error {
	return c.client.Del(ctx, movieKeyPrefix+id).Err()
}

func (c *Cache) SetSessionsByDateCache(ctx context.Context, date string, sessions []*domain.Session) error {
	if date == "" {
		return nil
	}
	data, err := json.Marshal(sessions)
	if err != nil {
		return fmt.Errorf("redis: marshal sessions: %w", err)
	}
	key := sessionsByDatePrefix + date
	return c.client.Set(ctx, key, data, sessionsByDateTTL).Err()
}

func (c *Cache) GetSessionsByDateCache(ctx context.Context, date string) ([]*domain.Session, error) {
	if date == "" {
		return nil, nil
	}
	key := sessionsByDatePrefix + date
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis: get sessions by date: %w", err)
	}
	var list []*domain.Session
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("redis: unmarshal sessions: %w", err)
	}
	return list, nil
}

func (c *Cache) DeleteSessionsByDateCache(ctx context.Context, date string) error {
	if date == "" {
		return nil
	}
	return c.client.Del(ctx, sessionsByDatePrefix+date).Err()
}

func (c *Cache) SetSessionCache(ctx context.Context, session *domain.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("redis: marshal session: %w", err)
	}
	key := sessionKeyPrefix + session.ID
	return c.client.Set(ctx, key, data, sessionByIDTTL).Err()
}

func (c *Cache) GetSessionCache(ctx context.Context, id string) (*domain.Session, error) {
	key := sessionKeyPrefix + id
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis: get session: %w", err)
	}
	var s domain.Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("redis: unmarshal session: %w", err)
	}
	// Restore time location
	s.StartTime = s.StartTime.UTC()
	return &s, nil
}

func (c *Cache) DeleteSessionCache(ctx context.Context, id string) error {
	return c.client.Del(ctx, sessionKeyPrefix+id).Err()
}

// SessionDateKey returns YYYY-MM-DD in UTC for a session start time.
func SessionDateKey(t time.Time) string {
	u := t.UTC()
	return u.Format("2006-01-02")
}
