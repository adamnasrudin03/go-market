package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adamnasrudin03/go-market/configs"
	"github.com/adamnasrudin03/go-market/pkg/driver"
	"github.com/sirupsen/logrus"
)

type CacheRepository interface {
	CreateCache(ctx context.Context, key string, data interface{}, ttl time.Duration)
	DeleteCache(ctx context.Context, key string)
	GetCache(ctx context.Context, key string, res interface{}) bool
}

type CacheRepo struct {
	Cache  driver.RedisClient
	Cfg    *configs.Configs
	Logger *logrus.Logger
}

func NewCacheRepository(
	redis driver.RedisClient,
	cfg *configs.Configs,
	logger *logrus.Logger,
) CacheRepository {
	return &CacheRepo{
		Cache:  redis,
		Cfg:    cfg,
		Logger: logger,
	}
}

func (r *CacheRepo) CreateCache(ctx context.Context, key string, data interface{}, ttl time.Duration) {
	var (
		opName = "CacheRepository-CreateCache"
		err    error
	)
	if ttl == 0 {
		ttl = r.Cfg.Redis.DefaultCacheTimeOut
	}

	err = r.Cache.Set(key, data, ttl)
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return
	}
}

func (r *CacheRepo) DeleteCache(ctx context.Context, key string) {
	var (
		opName = "CacheRepository-DeleteCache"
		err    error
	)
	err = r.Cache.Del(key)
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return
	}
}

func (r *CacheRepo) GetCache(ctx context.Context, key string, res interface{}) bool {
	var (
		opName = "CacheRepository-GetCache"
		err    error
	)

	data, err := r.Cache.Get(key)
	if err != nil {
		r.Logger.Errorf("%v error: %v ", opName, err)
		return false
	}

	err = json.Unmarshal([]byte(data), &res)
	if err != nil {
		r.Logger.Errorf("%v Unmarshal error: %v ", opName, err)
		return false
	}

	return true
}
