package pg

import (
	"arvan-challenge/application/config"
	"database/sql" // add this
	"log"

	_ "github.com/lib/pq" // add this
)

func GetDB(cfg config.AppConfig) *sql.DB {
	db, err := sql.Open("postgres", cfg.PGConfig.PG_CONN)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
