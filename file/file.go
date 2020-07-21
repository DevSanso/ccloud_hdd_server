package file

import (
	"os"
	"github.com/spf13/afero"

	"ccloud_hdd_server/data"
	"ccloud_hdd_server/db_sql"
	"ccloud_hdd_server/use_hash"


)


func ReadFile(fs afero.Fs, name string,key []byte,iv[]byte, header *db_sql.Header) (*data.Object,error)  {
	hash_name := string(use_hash.Sum([]byte(name)))
	f,err:=fs.Open(hash_name)
	if err != nil {return nil,err}
	return data.NewObject(f,key,iv,data.AES256,header.TokenSize(),header.Size())
}
func WriteFile(fs afero.Fs, name string,key []byte,iv[]byte,header *db_sql.Header) (*data.Object,error)  {
	hash_name := string(use_hash.Sum([]byte(name)))
	f,err := fs.OpenFile(hash_name,os.O_WRONLY | os.O_CREATE,os.FileMode(0))
	if err != nil {return nil,err}
	return data.NewObject(f,key,iv,data.AES256,header.TokenSize(),header.Size())
}

