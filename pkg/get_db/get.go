package get_db

import (
	"database/sql"
	"context"

	_ "github.com/go-sql-driver/mysql"
)
var db *sql.DB

func OpenDb(driverName , source string) (err error) {
	db,err = sql.Open(driverName,source)
	return
}

func GetDbConn(ctx context.Context) (*sql.Conn,error) {
	return db.Conn(ctx)
}

