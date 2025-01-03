package repository

import (
	"context"
	"encoding/json"
	"time"

	"auth-service/internal/domain/entity"
	"auth-service/internal/domain/repository"

	"github.com/redis/go-redis/v9"
)

type RedisSessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) repository.SessionRepository {
	return &RedisSessionRepository{client: client}
}

func (r *RedisSessionRepository) Create(ctx context.Context, sessionID string, session *entity.Session, ttl time.Duration) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, sessionID, data, ttl).Err()
}

func (r *RedisSessionRepository) Get(ctx context.Context, sessionID string) (*entity.Session, error) {
	data, err := r.client.Get(ctx, sessionID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var session entity.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *RedisSessionRepository) Delete(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, sessionID).Err()
}

func (r *RedisSessionRepository) DeleteAllUserSessions(ctx context.Context, userID string) error {
	pattern := "session:" + userID + ":*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}
