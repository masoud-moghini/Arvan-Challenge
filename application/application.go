package application

import (
	"arvan-challenge/application/config"
	"arvan-challenge/application/pg"
	"arvan-challenge/application/rds"
	"arvan-challenge/application/router"
	"arvan-challenge/application/syncjobs"
)

// represents application with its config and methods
type Application struct {
	Config                config.AppConfig
	Router                router.Router
	RedisClients          rds.RedisClients
	DatabaseQueries       pg.DatabaseQueries
	ApplicationJobRunners syncjobs.Runners
}
