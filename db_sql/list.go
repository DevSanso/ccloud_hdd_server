package db_sql

import (
	"database/sql"
	"context"
	"time"

)
type FileMeta struct {
	Name string
	Size int64
	Date time.Time
}
func GetFileHeaderListFromDir(conn *sql.Conn,baseId int,dir string)([]FileMeta,error) {
	const sql_query = "WITH base AS (SELECT sub_dir_path,name,size,date AS path"+
	"FROM object_table WHERE base_id= ?)"+
	"SELECT name,size,date FROM base WHERE sub_dir_path = ?;"
	rows,err:=conn.QueryContext(context.Background(),sql_query,baseId,dir);
	if err != nil {return nil,err}
	var res []FileMeta;var read = FileMeta{}
	for ok := true; ok ;ok = rows.Next(){
		err = rows.Scan(&(read.Name),&(read.Size),&(read.Date))
		if err != nil {break}
		res = append(res,read)
	}
	return res,err	
}