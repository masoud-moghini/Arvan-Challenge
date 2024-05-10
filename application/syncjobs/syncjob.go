package syncjobs

import (
	"arvan-challenge/application/pg"
	"context"
	"sort"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/redis/go-redis/v9"
)

type SyncWithDatabase interface {
	Synchronize(queries pg.Queries, s gocron.Scheduler)
}
type Runners struct {
	MaximumCacheSlotInMonthlyQuota uint32
	Job                            gocron.Job
	RedisMonthlyQuotaCache         *redis.Client
	DatabaseQueries                pg.Queries
}

// the method receives an interface of RelationalQueries
// implementation was not passed intentionally in order to preserve abstraction
// so sql.DB is not available in this function directly !!!!
func (r Runners) Synchronize(queries pg.Queries, s gocron.Scheduler) {
	sortedKeys := map[string]string{}
	r.Job, _ = s.NewJob(
		gocron.DurationJob(
			5*time.Minute,
		),
		gocron.NewTask(

			//remove more than 500 cache slots
			func(maximumAllowedSlicesInQuota int, redisDatabase *redis.Client) {
				// gather all keys inside go map
				allKeys, _ := redisDatabase.Keys(context.Background(), "*").Result()
				for i := 0; i < len(allKeys); i++ {
					content, _ := redisDatabase.HGetAll(context.Background(), allKeys[i]).Result()
					sortedKeys[allKeys[i]] = content["last_hit"]
				}

				//inplace sorting based on last_hit
				sort.SliceStable(allKeys, func(i, j int) bool {
					return sortedKeys[allKeys[i]] < sortedKeys[allKeys[j]]
				})

				//remove old enties
				key_size := len(allKeys)
				if key_size > int(r.MaximumCacheSlotInMonthlyQuota) {
					for i := maximumAllowedSlicesInQuota; i < key_size; i++ {
						queries.PreserveUsersRemainingQuotaInsideDatabase(allKeys[i], sortedKeys[allKeys[i]])
						redisDatabase.Del(context.Background(), sortedKeys[allKeys[i]])
					}
				}
			},
			r.MaximumCacheSlotInMonthlyQuota,
			r.RedisMonthlyQuotaCache,
		),
	)

	//todo add jobs for preserving all values periodically
}
