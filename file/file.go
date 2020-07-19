package file

import (
	"database/sql"
	"github.com/spf13/afero"
	"ccloud_hdd_server/db_sql"
	"ccloud_hdd_server/data"
)


func ReadFile(conn *sql.Conn,fs afero.Fs, name string) (*data.Object,error)  {

}
func WriteFile(conn *sql.Conn,fs afero.Fs, name string) (*data.Object,error)  {

}

func CreateFile(conn *sql.Conn,fs afero.Fs, name string) (*data.Object,error)  {

}