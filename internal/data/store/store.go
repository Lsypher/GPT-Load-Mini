package store

import "errors"

var ErrNotFound = errors.New("store: key not found")

type Store interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	HSet(key string, values map[string]any) error
	HGetAll(key string) (map[string]string, error)
	HIncrBy(key, field string, incr int64) (int64, error)
	LPush(key string, values ...any) error
	LRem(key string, count int64, value any) error
	Rotate(key string) (string, error)
	LLen(key string) (int64, error)
	Close() error
}
