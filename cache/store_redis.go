package cache

import (
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"time"
)

// RedisStore 实现 cache.Store interface
type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

func NewRedisStore(address string, username string, password string, db int) *RedisStore {
	rs := &RedisStore{}
	rs.RedisClient = redis.NewClient(address, username, password, db)
	rs.KeyPrefix = config.GetString("app.name") + ":cache:"
	return rs
}

func (s *RedisStore) Set(key string, value string, expireTime time.Duration) {
	s.RedisClient.Set(s.KeyPrefix+key, value, expireTime)
}

func (s *RedisStore) Get(key string) string {
	return s.RedisClient.Get(s.KeyPrefix + key)
}

func (s *RedisStore) Has(key string) bool {
	return s.RedisClient.Has(s.KeyPrefix + key)
}

func (s *RedisStore) Forget(key string) {
	s.RedisClient.Del(s.KeyPrefix + key)
}

func (s *RedisStore) Forever(key string, value string) {
	s.RedisClient.Set(s.KeyPrefix+key, value, 0)
}

func (s *RedisStore) PreDelAll(pre string) error {
	iter := s.RedisClient.Client.Scan(s.RedisClient.Context, 0, s.KeyPrefix+pre+"*", 0).Iterator()
	for iter.Next(s.RedisClient.Context) {
		s.RedisClient.Del(iter.Val())
	}
	return iter.Err()
}
func (s *RedisStore) Flush() {
	s.RedisClient.FlushDB()
}

func (s *RedisStore) Increment(parameters ...interface{}) {
	parameters[0] = s.KeyPrefix + parameters[0].(string)
	s.RedisClient.Increment(parameters...)
}

func (s *RedisStore) Decrement(parameters ...interface{}) {
	parameters[0] = s.KeyPrefix + parameters[0].(string)
	s.RedisClient.Decrement(parameters...)
}

func (s *RedisStore) IsAlive() error {
	return s.RedisClient.Ping()
}

func (s *RedisStore) Delete(key string) bool {
	return s.RedisClient.Del(s.KeyPrefix + key)
}

func (s *RedisStore) Exist(key string) (bool, error) {
	return s.RedisClient.Exist(s.KeyPrefix + key)
}
