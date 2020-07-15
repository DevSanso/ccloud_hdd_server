package db

import (
	"database/sql"
	"context"

	_ "github.com/go-sql-driver/mysql"


	hash "ccloud_hdd_server/use-hash"
)


const (
	_InsertUserPwSql = "INSERT INTO user(p_hash) VALUES(?);"
	_InsertUserBaseSql = "INSERT INTO user(base_id) VALUES(?) WHERE id = ?;"
	_SelectUserSql = "SELECT id FROM user WHERE pw = ?;"
	_SelectUserBaseIdSql = "SELECT base_id FROM user WHERE id = ?;"
	_CreateUserTableSql = "CREATE TABLE user(" +
							"id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY," + 
							"p_hash CHAR(32) NOT NULL," +
							"base_id INTEGER);"
)






func InsertPassword(conn *sql.Conn,passwdString string) error  {
	p_pass := []byte(passwdString)
	p_hash := hash.Sum(p_pass)

	_,err := conn.QueryContext(context.Background(),_InsertUserPwSql,p_hash)
	return err
}

func GetUserId(conn *sql.Conn,passwdHash []byte) (int,error) {
	row := conn.QueryRowContext(context.Background(),_SelectUserSql,passwdHash)
	var res int = 0
	err := row.Scan(&res)
	return res,err
}

func GetBasePathId(conn *sql.Conn,userID int) (int,error) {
	row := conn.QueryRowContext(context.Background(),_SelectUserBaseIdSql,userID)
	var res int = 0
	err := row.Scan(&res)
	return res,err
}

