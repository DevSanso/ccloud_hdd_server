package user

import (
	"context"
	"crypto/aes"
	"database/sql"
	"errors"
	"runtime"

	hash "ccloud_hdd_server/use_hash"
)


const (
	_InsertUserPwSql = "INSERT INTO user(p_hash,iv) VALUES(?,?);"
	_InsertUserBaseSql = "INSERT INTO user(base_id) VALUES(?) WHERE id = ?;"
	_SelectUserSql = "SELECT id FROM user WHERE p_hash = ?;"
	_SelectUserBaseIdSql = "SELECT base_id FROM user WHERE id = ?;"
	_CreateUserTableSql = "CREATE TABLE user(" +
							"id INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY," + 
							"p_hash CHAR(32) NOT NULL UNIQUE," +
							"iv VARCHAR(16) NOT NULL," +
							"base_id INTEGER);"
	_SelectUserIvSql = "SELECT iv FROM user WHERE id = ?;"
)


var (
	NotMatchBlockSizeErr = errors.New("NotMatchBlockSizeErr")
)



func InsertUser(conn *sql.Conn,passwdString string,iv []byte) error  {
	p_pass := []byte(passwdString)
	p_hash := hash.Sum(p_pass)

	if len(iv) != aes.BlockSize {
		return NotMatchBlockSizeErr
	}
	_,err := conn.QueryContext(context.Background(),_InsertUserPwSql,p_hash,string(iv))
	return err
}

func GetUserIv(conn *sql.Conn,userId int) ([]byte,error) {
	var res string
	row := conn.QueryRowContext(context.Background(),_SelectUserIvSql,userId)
	err := row.Scan(&res)
	b_res := []byte(res)
	if len(b_res) != aes.BlockSize {
		pc,_,_,_ :=runtime.Caller(5)
		fun :=runtime.FuncForPC(pc)
		panic("user iv block size, crypto blocksize not match:" + fun.Name())
	}
	return b_res,err
}

func GetUserId(conn *sql.Conn,passwdHash []byte) (int,error) {
	row := conn.QueryRowContext(context.Background(),_SelectUserSql,passwdHash)
	var res int = 0
	err := row.Scan(&res)
	return res,err
}

func GetBasePathId(conn *sql.Conn,userId int) (int,error) {
	row := conn.QueryRowContext(context.Background(),_SelectUserBaseIdSql,userId)
	var res int = 0
	err := row.Scan(&res)
	return res,err
}

func IsExistUserPasswd(conn *sql.Conn,passwdString string) (bool,error) {
	h :=hash.Sum([]byte(passwdString))
	const s = "SELECT EXISTS(SELECT * FROM user WHERE p_hash = ?);"
	r:=conn.QueryRowContext(context.Background(),s,h)
	var res int
	err := r.Scan(&res)
	if res == 1 {
		return true,err
	}else {
		return false,err
	}
}



