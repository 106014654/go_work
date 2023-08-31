package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"go_work/user_webook/internal/domain"
)

var ERRREDISNIL = redis.Nil

type RedisUserCache struct {
	// 传单机 Redis 可以
	// 传 cluster 的 Redis 也可以
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) *RedisUserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (cache *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.Key(id)

	value, err := cache.client.Get(ctx, key).Bytes()

	if err != nil || err == ERRREDISNIL {
		return domain.User{}, err
	}

	var user domain.User

	err = json.Unmarshal(value, &user)

	return user, nil
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.Key(u.Id)

	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}
