package pg

import (
	"database/sql" // add this
	"fmt"

	_ "github.com/lib/pq" // add this
)

type (
	Queries interface {
		PreserveUsersRemainingQuotaInsideDatabase(user_id string, remaining_quota int)
	}
	DatabaseQueries struct {
		DBObject *sql.DB
	}
)

func (dQueries DatabaseQueries) PreserveUsersRemainingQuotaInsideDatabase(user_id string, remaining_quota string) {
	fmt.Println("runinng query in database")
}
