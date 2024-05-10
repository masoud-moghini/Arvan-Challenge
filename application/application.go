package application

import (
	"arvan-challenge/application/config"
	"arvan-challenge/application/rds"
	"arvan-challenge/application/router"
	"database/sql"
)

// represents application with its config and methods
type Application struct {
	Config       config.AppConfig
	Router       router.Router
	RedisClients rds.RedisClients
	PQsql        *sql.DB
}
