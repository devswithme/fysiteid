package helper

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHelper interface {
	Client() *redis.Client
	Set(key string, value any, ttl time.Duration) error
	Script(key string) (int, error)
	Exist(key string) (bool, error)
	Get(key string) (string, error)
	Incr(key string) error
	Delete(key string) error
}

type redisHelper struct {
	client  *redis.Client
	timeout time.Duration
	prefix  string
	script  string
}

func NewRedisHelper(addr string, password string, db int, prefix string, timeout time.Duration) RedisHelper {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &redisHelper{
		client:  client,
		timeout: timeout,
		prefix:  prefix,
		script: `
            if tonumber(redis.call("GET", KEYS[1])) > 0 then
                redis.call("DECR", KEYS[1])
                return 1
            else
                return 0
            end
        `,
	}
}

func (r *redisHelper) Client() *redis.Client {
	return r.client
}

func (r *redisHelper) Set(key string, value any, ttl time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Set(ctx, r.prefix+key, value, ttl).Err()
}

func (r *redisHelper) Exist(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	count, err := r.client.Exists(ctx, key).Result()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *redisHelper) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Get(ctx, r.prefix+key).Result()
}

func (r *redisHelper) Script(key string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return redis.NewScript(r.script).Run(ctx, r.client, []string{r.prefix + key}, nil).Int()
}

func (r *redisHelper) Incr(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Incr(ctx, key).Err()
}

func (r *redisHelper) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Del(ctx, r.prefix+key).Err()
}
