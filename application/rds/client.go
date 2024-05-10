package rds

import (
	"arvan-challenge/application/config"
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClients struct {
	RedisClientForMinuteQuota   *redis.Client
	RedisClientForMonthQuota    *redis.Client
	RedisClientForTimestamp     *redis.Client
	RedisClientForDataProcessed *redis.Client
}

func RedisClientForMinuteQuota(cfg config.AppConfig) *redis.Client {
	_ = context.Background()
	return redis.NewClient(cfg.RedisConfigForMinuteQuotaDB)
}

func RedisClientForDataProcessing(cfg config.AppConfig) *redis.Client {
	_ = context.Background()
	return redis.NewClient(cfg.RedisConfigForDataProcessingDB)
}
func RedisClientForMonthQuota(cfg config.AppConfig) *redis.Client {
	_ = context.Background()
	return redis.NewClient(cfg.RedisConfigForMonthQuotaDB)
}
