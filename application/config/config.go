package config

import (
	"github.com/redis/go-redis/v9"
)

type (
	PGcfg struct {
		PG_CONN string `default:"host=postgres dbname=arvan_queue user=postgres password=itsasecret"`
	}
	AppConfig struct {
		RedisConfigForMinuteQuotaDB    *redis.Options
		RedisConfigForMonthQuotaDB     *redis.Options
		RedisConfigForDataProcessingDB *redis.Options
		PGConfig                       PGcfg
	}
)

func InitConfig() (cfg AppConfig) {
	RedisAddr := "redis:6379"
	RedisPasswd := ""
	RedisMinuteQuotaDB := 0
	RedisMonthQuotaDB := 1

	redisCfgForMonthQuota := &redis.Options{
		Addr:     RedisAddr,
		Password: RedisPasswd, // no password set
		DB:       RedisMonthQuotaDB,
	}

	redisCfgMinuteQuotaDB := &redis.Options{
		Addr:     RedisAddr,
		Password: RedisPasswd, // no password set
		DB:       RedisMinuteQuotaDB,
	}
	pgCfg := PGcfg{}
	return AppConfig{
		RedisConfigForMinuteQuotaDB: redisCfgMinuteQuotaDB,
		RedisConfigForMonthQuotaDB:  redisCfgForMonthQuota,
		PGConfig:                    pgCfg,
	}

}
