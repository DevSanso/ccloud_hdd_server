package worker

import (
	"context"
	"encoding/json"
	"net/http"

	"path/filepath"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/spf13/afero"

	"ccloud_hdd_server/auth"
	"ccloud_hdd_server/data"
	"ccloud_hdd_server/db_sql"
	"ccloud_hdd_server/get_db"
	err_msg "ccloud_hdd_server/worker/internal"
)

func mergePath(dir, file string)string {
	return filepath.FromSlash(dir + "/" + file)
}

type FileMeta struct{
	Name string
	Offset int64
	Size int64
	FullSize int64
}
type FileDataServ struct {}


func (fvs *FileDataServ)getUserBaseId(key int) (int,error) {
	c,err := get_db.GetDbConn(context.Background())
	if err != nil {return 0,err}
	defer c.Close()

	return db_sql.GetBasePathId(c,key)

}
func (fvs *FileDataServ)getObjectHeader(path string)(db_sql.Header,error) {
	c,err := get_db.GetDbConn(context.Background())
	if err != nil {return nil,err}
	defer c.Close()
	return db_sql.LoadObjectHeader(c,path)
}

func (fvs *FileDataServ)Do(w http.ResponseWriter,r *http.Request,key []byte) {
	ck,ck_err := r.Cookie("session")
	if ck_err != nil {err_msg.NotLoginResponse(w);return}
	var using_key int
	using_key,ck_err = strconv.Atoi(ck.Value)

	if ck_err != nil {err_msg.BadCookieValueResponse(w);return}
	if using_key,ck_err =auth.GetUesrId(uint32(using_key));ck_err != nil {
		err_msg.BadCookieValueResponse(w)
		return;
	}

	


	file_name := r.URL.Query().Get("file")
	dir_path := r.URL.Query().Get("dir")

	header,db_err := fvs.getObjectHeader(dir_path+"/"+file_name)
	if db_err != nil {err_msg.CantSearchDataResponse(w);return}

	fs := afero.NewBasePathFs(afero.NewOsFs(),header.BaseDir())
	path := filepath.Join(header.SubDir(),header.Name())
	f,f_err :=fs.Open(path)
	if f_err != nil {err_msg.BadCookieValueResponse(w);return}

	upgrader := websocket.Upgrader{}
	ws_conn ,ws_err := upgrader.Upgrade(w,r,nil)
	if ws_err != nil {err_msg.CantStartWSLoopResponse(w);return}
	fvs.loop(ws_conn,f,key)
	
}
func(fvs *FileDataServ)loop(conn *websocket.Conn,f afero.File,key []byte) {
	for {

	}
}


