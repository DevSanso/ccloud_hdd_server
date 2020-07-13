package db

import (
	"database/sql"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"hash"


	_ "github.com/go-sql-driver/mysql"
)


const (
_InsertUserPwSql = "INSERT INTO user(p_hash) VALUES(?);"
	_InsertUserBaseSql = "INSERT INTO user(base_id) VALUES(?) WHERE id = ?;"
	_SelectUserSql = "SELECT index FROM user WHERE pw = ?;"
	_CreateUserTableSql = "CREATE TABLE user(" +
							"id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY," + 
							"p_hash CHAR(32) NOT NULL," +
							"base_id INTEGER);"
)


type _HashDecryption struct {}
func (*_HashDecryption)first() hash.Hash {
	return md5.New()
}
func (*_HashDecryption)second() hash.Hash {
	return sha256.New()
}

var _hd = _HashDecryption{}



func InsertPassword(conn *sql.Conn,passwd string) error  {
	p_pass := []byte(passwd)
	f := _hd.first()
	s := _hd.second()
	p_hash := s.Sum(f.Sum(p_pass))
	_,err := conn.QueryContext(context.Background(),_InsertPasswdSql,p_hash)
	return err
}