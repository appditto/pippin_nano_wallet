package database

import (
	"database/sql"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func NewEntClient(connInfo SqlDBConn) (*ent.Client, error) {
	db, err := sql.Open(connInfo.Dialect(), connInfo.DSN())
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(connInfo.Dialect(), db)
	return ent.NewClient(ent.Driver(drv)), nil
}
