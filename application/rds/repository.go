package rds

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	InMemoryServices struct {
		RedisClients RedisClients
	}
)

// returns user remained quota and decrease it if its value > 0
// if value is not present in cache, restore it from database
func (i InMemoryServices) GetAndDecreaseMonthlyQuota(ctx context.Context, user_id string, clientForMonthQuota *redis.Client) (int, error) {
	val, err := clientForMonthQuota.HGet(ctx, user_id, "remaining_quota").Result()

	if err == redis.Nil {

		//recover monthly quota to redis
		val = recoverUserQuotaFromPersistanceDb(user_id)
	} else if err != nil {
		fmt.Println("failed process data %+v\n", err)
		return 0, errors.New("unexpected error in retrieving cached value")
	}

	valInNumber, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("failed process data %+v\n", err)
		return 0, errors.New("unexpected error in retrieving cached value")
	}
	if valInNumber > 0 {
		valInNumber -= 1
	}

	defer clientForMonthQuota.HSet(ctx, user_id,
		"remaining_quota", valInNumber,
		"last_hit", time.Now().String(),
	)
	return valInNumber, nil
}
func (i InMemoryServices) GetAndDecreaseMinuteQuota(ctx context.Context, user_id string, clientForMinuteQuota *redis.Client) (int, error) {
	val, err := clientForMinuteQuota.Get(ctx, user_id).Result()

	if err == redis.Nil {

		//recover monthly quota to redis
		val = recoverUserMinuteQuotaFromPersistanceDb(user_id)
	}
	if err != nil {
		fmt.Println("failed process data %+v\n", err)
		return 0, errors.New("unexpected error in retrieving cache value")
	}

	valInNumber, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("failed process data %+v\n", err)
		return 0, errors.New("unexpected error in retrieving cache value")
	}
	if valInNumber > 0 {
		valInNumber -= 1
	}

	defer clientForMinuteQuota.HSetNX(ctx, user_id, strconv.Itoa(valInNumber), 1*time.Minute)
	return valInNumber, nil
}

func recoverUserMinuteQuotaFromPersistanceDb(user_id string) string {
	return "5"
}
func recoverUserQuotaFromPersistanceDb(user_id string) string {
	return "5"
}
