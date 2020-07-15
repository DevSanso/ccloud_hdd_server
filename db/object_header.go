package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"time"

	hash "ccloud_hdd_server/use-hash"
)

const (
	_CreateFHeaderSql = ""
	_InsertFHeaderSql = ""
	_UpdateFHeaderSql = ""
	_DeleteFHeaderSql = ""

	_SelectFHeaderSql = ""

	_CreateFHeaderTableSql = ""
)

type Header interface {
	os.FileInfo
	DirPath() string
	TokenSize() int
	HashName() []byte
}

type objectHeader struct {
	name string
	basePath string
	subDirPath string

	tokenSize int
	//본래의 파일 사이즈,하드에 있는 파일사이즈랑 다름
	size int64
	date string
}

func LoadObjectHeader(conn *sql.Conn,name string) (Header,error) {
	row := conn.QueryRowContext(context.Background(),"")
	var res struct {
		sub_dir string
		base_id int
		token_size int
		date string
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

func (oh *objectHeader)Name() string {return oh.name}
func (oh *objectHeader)Size() int64 {return oh.size}
func (oh *objectHeader)Mode() os.FileMode {return os.FileMode(0)}
func (ob *objectHeader)ModTime() time.Time {
	t,_ := time.Parse(time.RFC3339,ob.date)
	return t
}
func (*objectHeader)IsDir() bool {return false}
func (*objectHeader)Sys() interface{} {return nil}
func (ob *objectHeader)DirPath() string {
	return filepath.FromSlash(ob.basePath + "/" + ob.subDirPath)
}
func (ob *objectHeader)TokenSize() int {return ob.tokenSize}
func (ob *objectHeader)HashName() []byte {
	return hash.Sum([]byte(ob.name))
}