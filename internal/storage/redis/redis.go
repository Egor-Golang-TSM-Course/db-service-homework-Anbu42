package redis

import (
	"blog/internal/config"
	"blog/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client          *redis.Client
	cacheExpiration time.Duration
}

func NewRedisCache(cfg *config.Config) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Components.Cache.Host, cfg.Components.Cache.Port),
	})

	return &RedisCache{
		client:          client,
		cacheExpiration: cfg.Components.Cache.CacheExpiration,
	}
}

func (r *RedisCache) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key).Result()
}

func (r *RedisCache) Set(key string, value string) error {
	return r.client.Set(context.Background(), key, value, r.cacheExpiration).Err()
}

func (r *RedisCache) SetPosts(pageSize, page int, tags []string, date *time.Time, posts []*models.Post) error {
	// Собираем уникальный ключ для кеширования
	cacheKey := fmt.Sprintf("posts:%d:%d:%v:%v", pageSize, page, tags, date)

	postsJSON, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), cacheKey, postsJSON, r.cacheExpiration).Err()
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *RedisCache) UpdatePost(key string, updatedPost *models.Post) error {
	currentPostJSON, err := r.Get(key)
	if err != nil {
		return err
	}

	var currentPost models.Post
	err = json.Unmarshal([]byte(currentPostJSON), &currentPost)
	if err != nil {
		return err
	}

	currentPost.Title = updatedPost.Title
	currentPost.Content = updatedPost.Content

	updatedPostJSON, err := json.Marshal(currentPost)
	if err != nil {
		return err
	}

	return r.Set(key, string(updatedPostJSON))
}
