package main

import (
	"arvan-challenge/application"
	"arvan-challenge/application/config"
	"arvan-challenge/application/pg"
	"arvan-challenge/application/rds"
	"arvan-challenge/application/router"
	"arvan-challenge/application/syncjobs"
	"database/sql"
	"log"

	"github.com/go-co-op/gocron/v2"

	"github.com/go-chi/chi/v5"
)

func main() {
	app := InitiateConfigs()
	rtHndlr := app.Router
	rtHndlr.Listen()
}
func InitiateConfigs() application.Application {

	//initiate application configurations
	var cfg config.AppConfig = config.InitConfig()

	var redisClients rds.RedisClients = rds.RedisClients{
		//for caching users minute quota
		//this kept using expiration time
		//if users minute quota has been expired, it will be restored from Relational db
		RedisClientForMinuteQuota: rds.RedisClientForMinuteQuota(cfg),

		//for caching users monthly quota and last activation value
		//db size is limited to CONSTANT number of last active users
		//if users minute quota does not exist, it will be restored from Relational db
		//a specific job exists to keep size of db size, runs asyncronously with go
		RedisClientForMonthQuota: rds.RedisClientForMonthQuota(cfg),

		//for caching processed data
		//db size is kept by using expiration keys
		//if requested data exists in cache, it will returned back to user, irrespective of his quota
		RedisClientForDataProcessed: rds.RedisClientForDataProcessing(cfg),
	}

	//define request handler and pass it to application
	var routerHandler router.RouterHandler = router.RouterHandler{
		Routes: chi.NewMux(),
		RequestHandlers: router.RequestHandlers{

			InMemoryServices: rds.InMemoryServices{
				RedisClients: redisClients,
			},
		},
	}

	//users remaining quota is kept in relational database
	//it is intervally synchronized with redis cache
	//WE HAVE NOT IMPLEMENTED RELATIONAL DATABASE TO PERSIST USERS QUOTA!
	databaseQueries := InitiateRelationalDataBase(cfg)
	return application.Application{
		Config:                cfg,
		Router:                routerHandler,
		RedisClients:          redisClients,
		DatabaseQueries:       databaseQueries,
		ApplicationJobRunners: InitiateRunnersAndJobs(databaseQueries),
	}
}

func InitiateRelationalDataBase(cfg config.AppConfig) pg.DatabaseQueries {
	db, err := sql.Open("postgres", cfg.PGConfig.PG_CONN)
	if err != nil {
		log.Fatal(err)
	}
	return pg.DatabaseQueries{DBObject: db}
}
func InitiateRunnersAndJobs(queries pg.DatabaseQueries) syncjobs.Runners {
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		panic(err)
	}
	runner := syncjobs.Runners{
		MaximumCacheSlotInMonthlyQuota: 500,
	}
	runner.Synchronize(queries, s)
	return runner
}
