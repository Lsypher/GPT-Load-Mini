package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}
	return &RedisStore{client: client}, nil
}

func (s *RedisStore) Set(key string, value []byte) error {
	return s.client.Set(context.Background(), key, value, 0).Err()
}

func (s *RedisStore) Get(key string) ([]byte, error) {
	val, err := s.client.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, ErrNotFound
	}
	return val, err
}

func (s *RedisStore) Delete(key string) error {
	return s.client.Del(context.Background(), key).Err()
}

func (s *RedisStore) HSet(key string, values map[string]any) error {
	return s.client.HSet(context.Background(), key, values).Err()
}

func (s *RedisStore) HGetAll(key string) (map[string]string, error) {
	return s.client.HGetAll(context.Background(), key).Result()
}

func (s *RedisStore) HIncrBy(key, field string, incr int64) (int64, error) {
	return s.client.HIncrBy(context.Background(), key, field, incr).Result()
}

func (s *RedisStore) LPush(key string, values ...any) error {
	return s.client.LPush(context.Background(), key, values...).Err()
}

func (s *RedisStore) LRem(key string, count int64, value any) error {
	return s.client.LRem(context.Background(), key, count, value).Err()
}

func (s *RedisStore) Rotate(key string) (string, error) {
	return s.client.LMove(context.Background(), key, key, "RIGHT", "LEFT").Result()
}

func (s *RedisStore) LLen(key string) (int64, error) {
	return s.client.LLen(context.Background(), key).Result()
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}

func FormatGroupKeysKey(groupID uint) string {
	return fmt.Sprintf("group:%d:keys", groupID)
}
