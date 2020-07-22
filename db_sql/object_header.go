package db_sql

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const (
	_InsertHeaderSql = "INSERT INTO object_table VALUES(?,?,?,?,?,?);"

	_UpdateHeaderSubDirPathSql = "UPDATE object_table SET sub_dir_path = ? WHERE name = ?;"
	_UpdateHeaderSizeSql       = "UPDATE object_table SET size = ? WHERE name = ?;"
	_UpdateHeaderDateSql       = "UPDATE object_table SET date = ? WHERE name = ?;"

	_DeleteHeaderSql = "DELTE FROM object_table WHERE name = ?;"

	_SelectHeaderSql = "SELECT * FROM object_table WHERE sub_dir_dir = ? AND name = ?   LIMIT 1;"

	_CreateHeaderTableSql = "CREATE TABLE object_table (" +
		"name VARCHAR(256) NOT NULL," +
		"base_id INT," +
		"sub_dir_path VARCHAR(256)," +
		"token_size INT," +
		"size BIGINT," +
		"date DATETIME);"
)

var (
	NilParameterErr   = errors.New("NilParameterError")
	NotMatchColumnErr = errors.New("NotMatchColumnErr")
	AlearyExistErr    = errors.New("AlearyExistErr")
)


type Header struct {
	//수정 불가능
	name string
	//수정 불가능
	basePath    string
	subDirPath  string
	//수정 불가능
	tokenSize int
	//본래의 파일 사이즈,하드에 있는 파일사이즈랑 다름
	size int64
	date time.Time
}

type ObjectConfig struct {
	BaseId      int
	SubDirPath  string
	TokenSize   int
	Size        int64
	Date        time.Time
}

func isExistFile(conn *sql.Conn, name string, cfg *ObjectConfig) (bool, error) {
	const sql_exist = "SELECT EXISTS(" +
		"WITH paths AS (SELECT CONCAT(sub_dir_path,\"/\",name) AS path" +
		"FROM object_table WHERE base_id = ? )" +
		"SELECT path FROM paths WHERE = ?);"
	row := conn.QueryRowContext(context.Background(),
		sql_exist, cfg.BaseId, cfg.SubDirPath+"/"+name)
	var res int = 0
	err := row.Scan(&res)
	if res == 0 {
		return false, err
	} else {
		return true, err
	}
}
func CreateHeader(conn *sql.Conn, name string, cfg *ObjectConfig) (*Header, error) {
	if cfg == nil {
		return nil, NilParameterErr
	}

	row := conn.QueryRowContext(context.Background(), _SelectBasePathSql, cfg.BaseId)
	var base_path string
	if s_err := row.Scan(&base_path); s_err != nil {
		return nil, s_err
	}

	if ok, is_err := isExistFile(conn, name, cfg); is_err != nil {
		return nil, is_err
	} else if ok {
		return nil, AlearyExistErr
	}

	tx, err := conn.BeginTx(context.Background(),
		&sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return nil, err
	}
	_, err = tx.Query(_InsertHeaderSql,
		name, cfg.BaseId, cfg.SubDirPath, cfg.TokenSize, cfg.Size, cfg.Date)

	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}
	return &Header{
		name,
		base_path,
		cfg.SubDirPath,
		cfg.TokenSize,
		cfg.Size,
		cfg.Date,
	}, nil
}

func splitPath(path string) (name string, subDir string) {
	name = filepath.Base(path)
	subDir = filepath.Dir(subDir)
	return
}

func LoadHeader(conn *sql.Conn, path string) (*Header, error) {
	name, sub_dir := splitPath(path)
	row := conn.QueryRowContext(context.Background(), _SelectHeaderSql, name, sub_dir)
	var res struct {
		sub_dir    string
		base_id    int
		token_size int
		date       time.Time
		size       int64
	}
	err := row.Scan(&(res.base_id), &(res.sub_dir),
		 &(res.token_size), &(res.size), &(res.date))
	if err != nil {
		return nil, err
	}
	row = conn.QueryRowContext(context.Background(), _SelectBasePathSql, res.base_id)
	var basePath string
	err = row.Scan(&basePath)
	if err != nil {
		return nil, err
	}

	return &Header{
		name,
		basePath,
		res.sub_dir,
		res.token_size,
		res.size,
		res.date,
	}, nil
}

func DeleteHeader(conn *sql.Conn, name string) error {
	tx, err := conn.BeginTx(context.Background(),
		&sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: false})
	if err != nil {
		return err
	}
	_, err = tx.Query(_DeleteHeaderSql, name)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (oh *Header) Name() string        { return oh.name }
func (oh *Header) Size() int64         { return oh.size }
func (oh *Header) Mode() os.FileMode   { return os.FileMode(0) }
func (ob *Header) ModTime() time.Time  { return ob.date }
func (ob *Header) BaseDir() string     { return filepath.FromSlash(ob.basePath) }
func (ob *Header) SubDir() string      { return filepath.FromSlash(ob.subDirPath) }
func (ob *Header) TokenSize() int      { return ob.tokenSize }


