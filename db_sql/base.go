package db_sql

import (
	"database/sql"
	"context"
	_ "github.com/go-sql-driver/mysql"
)

const (
	_SelectBasePathSql = "SELECT path FROM hdd_base_path WHERE id = ?;"
	_InsertBasePathSql = "INSERT INTO hdd_base_path(id,path) VALUES(?,?);"
	_CreateBaseTableSql = "CREATE TABLE hdd_base_path(" +
					"id INTEGER NULL AUTO_INCREMENT PRIMARY KEY ," +
					"path VARCHAR(256));" 
)

func GetBasePath(conn *sql.Conn,id int) (string,error) {
	row := conn.QueryRowContext(context.Background(),_SelectBasePathSql,id)
	var path string
	if err := row.Scan(&path);err != nil {return "",err}
	return path,nil
}
func InsertBasePath(conn *sql.Conn,id int,path string) error {
	_,err := conn.QueryContext(context.Background(),_InsertBasePathSql,id,path)
	return err
}

func CreateBasePathTable(db *sql.DB) error {
	_,err := db.Exec(_CreateBaseTableSql)
	return err
}