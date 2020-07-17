package get_db

import (
	"database/sql"
	"context"

	_ "github.com/go-sql-driver/mysql"
)


func GetDbConn(ctx context.Context) (*sql.Conn,error) {

}