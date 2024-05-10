package main

import (
	"arvan-challenge/application"
	"arvan-challenge/application/config"
	"arvan-challenge/application/pg"
	"arvan-challenge/application/rds"
	"arvan-challenge/application/router"
	"database/sql"
)

func main() {
	app := InitiateConfigs()
	rtHndlr := app.Router
	rtHndlr.Listen()
}
func InitiateConfigs() application.Application {

	var cfg config.AppConfig = config.InitConfig()

	var redisClients rds.RedisClients = rds.RedisClients{
		RedisClientForMinuteQuota: rds.RedisClientForMinuteQuota(cfg),
		RedisClientForMonthQuota:  rds.RedisClientForMonthQuota(cfg),
		RedisClientForTimestamp:   rds.RedisClientForTimestamp(cfg),
	}
	var routerHandler router.RouterHandler = router.RouterHandler{
		RequestHandlers: router.RequestHandlers{
			InMemoryServices: rds.InMemoryServices{
				RedisClients: redisClients,
			},
		},
	}
	var databaseClient *sql.DB = pg.GetDB(cfg)
	return application.Application{
		Config:       cfg,
		Router:       routerHandler,
		RedisClients: redisClients,
		PQsql:        databaseClient,
	}
}
