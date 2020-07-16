package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	hash "ccloud_hdd_server/use-hash"
)

const (

	_InsertObjectHeaderSql = "INSERT INTO object_table VALUES(?,?,?,?,?,?);"

	_UpdateObjectHeaderSubDirPathSql = "UPDATE object_table SET sub_dir_path = ? WHERE name = ?;"
	_UpdateObjectHeaderSizeSql = "UPDATE object_table SET size = ? WHERE name = ?;"
	_UpdateObjectHeaderDateSql = "UPDATE object_table SET date = ? WHERE name = ?;"

	_DeleteObjectHeaderSql = "DELTE FROM object_table WHERE name = ?;"


	_SelectObjectHeaderSql = "SELECT * FROM object_table WHERE name = ? LIMIT 1;"

	_CreateObjectHeaderTableSql =  "CREATE TABLE object_table ("+
		"name VARCHAR(256) NOT NULL PRIMARY KEY,"+
		"base_id INT,"+
		"sub_dir_path VARCHAR(256),"+
		"token_size INT,"+
		"size BIGINT,"+
		"date DATETIME);"
)
type HeaderColumn int
const (
	SubPath HeaderColumn = iota
	Size
	Date
)

var (
	NilParameterErr = errors.New("NilParameterError")
	NotMatchColumnErr = errors.New("NotMatchColumnErr")
)

type Header interface {
	os.FileInfo
	DirPath() string
	TokenSize() int
	HashName() []byte
	UpdateValue(conn *sql.Conn,column HeaderColumn,value interface{})error
}

type objectHeader struct {
	name string
	basePath string
	subDirPath string

	//수정 불가능
	tokenSize int
	//본래의 파일 사이즈,하드에 있는 파일사이즈랑 다름
	size int64
	date time.Time
}

type ObjectConfig struct {
	BaseId int
	SubDirPath string
	TokenSize int
	Size int64
	Date time.Time
}

func CreateObjectHeader(conn *sql.Conn,name string,cfg *ObjectConfig) (Header,error) {
	if cfg == nil {return nil,NilParameterErr}

	row := conn.QueryRowContext(context.Background(),_SelectBasePathSql,cfg.BaseId)
	var base_path string
	if s_err := row.Scan(&base_path);s_err != nil {return nil,s_err}


	tx ,err := conn.BeginTx(context.Background(),
	&sql.TxOptions{Isolation: sql.LevelReadCommitted,ReadOnly : false})
	if err != nil {return nil,err}
	_,err = tx.Query(_InsertObjectHeaderSql,
		name,cfg.BaseId,cfg.SubDirPath,cfg.TokenSize,cfg.Size,cfg.Date)

	if err != nil {
		tx.Rollback()
		return nil,err
	}
	if err = tx.Commit();err!=nil {tx.Rollback();return nil,err}
	return &objectHeader{
		name,
		base_path,
		cfg.SubDirPath,
		cfg.TokenSize,
		cfg.Size,
		cfg.Date,
	},nil
}
func LoadObjectHeader(conn *sql.Conn,name string) (Header,error) {
	row := conn.QueryRowContext(context.Background(),_SelectObjectHeaderSql,name)
	var res struct {
		sub_dir string
		base_id int
		token_size int
		date time.Time
		size int64
	}
	err:=row.Scan(&(res.base_id),&(res.sub_dir),&(res.token_size),&(res.size),&(res.date))
	if err != nil {return nil,err}
	row = conn.QueryRowContext(context.Background(),_SelectBasePathSql,res.base_id)
	var basePath string;err = row.Scan(&basePath)
	if err != nil {return nil,err}
	
	return &objectHeader{
		name,
		basePath,
		res.sub_dir,
		res.token_size,
		res.size,
		res.date,
	},nil	
}
func DeleteObjectHeader(conn *sql.Conn,name string) error {
	tx ,err := conn.BeginTx(context.Background(),
	&sql.TxOptions{Isolation: sql.LevelSerializable,ReadOnly : false})
	if err != nil {return err}
	_,err = tx.Query(_DeleteObjectHeaderSql,name)
	if err != nil {tx.Rollback();return err}

	return tx.Commit()
}


func (oh *objectHeader)Name() string {return oh.name}
func (oh *objectHeader)Size() int64 {return oh.size}
func (oh *objectHeader)Mode() os.FileMode {return os.FileMode(0)}
func (ob *objectHeader)ModTime() time.Time {return ob.date}
func (*objectHeader)IsDir() bool {return false}
func (*objectHeader)Sys() interface{} {return nil}
func (ob *objectHeader)DirPath() string {
	return filepath.FromSlash(ob.basePath + "/" + ob.subDirPath)
}
func (ob *objectHeader)TokenSize() int {return ob.tokenSize}
func (ob *objectHeader)HashName() []byte {
	return hash.Sum([]byte(ob.name))
}
func(ob *objectHeader)UpdateValue(conn *sql.Conn,
	column HeaderColumn,value interface{})error {

	tx ,err := conn.BeginTx(context.Background(),
	&sql.TxOptions{Isolation: sql.LevelReadCommitted,ReadOnly : false})
	if err != nil {return err}
	var query string = ""
	switch column{
	case SubPath:
		query = _UpdateObjectHeaderSubDirPathSql
	case Size:
		query = _UpdateObjectHeaderSizeSql
	case Date:
		query = _UpdateObjectHeaderDateSql
	default:
		return NotMatchColumnErr
	}

	_,err = tx.Query(query,value,ob.name)
	if err != nil {tx.Rollback();return err}
	ob.patchValue(column,value)
	return tx.Commit()
}
func(ob *objectHeader)patchValue(column HeaderColumn,value interface{}) {
	switch column{
	case SubPath:
		ob.subDirPath = value.(string)		
	case Size:
		ob.size = value.(int64)
	case Date:
		ob.date = value.(time.Time)
	default:		
	}
}