package worker

import (
	"context"
	"database/sql"
	"net/http"

	"path/filepath"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/spf13/afero"

	"ccloud_hdd_server/auth"
	"ccloud_hdd_server/data"
	"ccloud_hdd_server/db_sql"
	db_user "ccloud_hdd_server/db_sql/user"
	"ccloud_hdd_server/file"
	"ccloud_hdd_server/get_db"

	err_msg "ccloud_hdd_server/worker/internal"
)

func mergePath(dir, file string) string {
	return filepath.FromSlash(dir + "/" + file)
}

type FileMeta struct {
	Name     string
	Offset   int64
	Size     int64
	FullSize int64
}
type FileDataServ struct{}

func (fvs *FileDataServ) getUserBaseId(c *sql.Conn,key int) (int, error) {
	return db_user.GetBasePathId(c, key)

}
func (fvs *FileDataServ) getHeader(c *sql.Conn,path string) (*db_sql.Header, error) {
	return db_sql.LoadHeader(c, path)
}

func (fvs *FileDataServ) Do(w http.ResponseWriter, r *http.Request, key []byte) {
	ck, ck_err := r.Cookie("session")
	if ck_err != nil {
		err_msg.NotLoginResponse(w)
		return
	}
	var using_key int
	using_key, ck_err = strconv.Atoi(ck.Value)

	if ck_err != nil {
		err_msg.BadCookieValueResponse(w)
		return
	}
	if using_key, ck_err = auth.GetUesrId(uint32(using_key)); ck_err != nil {
		err_msg.BadCookieValueResponse(w)
		return
	}

	file_name := r.URL.Query().Get("file")
	dir_path := r.URL.Query().Get("dir")
	c, db_err := get_db.GetDbConn(context.Background())
	if db_err != nil {
		err_msg.CantConnectDbResponse(w);return
	}
	var header *db_sql.Header
	header, db_err = fvs.getHeader(c,dir_path + "/" + file_name)
	if db_err != nil {
		err_msg.CantSearchDataResponse(w)
		return
	}
	var iv []byte
	iv,db_err = db_user.GetUserIv(c,using_key)
	if db_err != nil {
		err_msg.CantSearchDataResponse(w)
		return
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), header.BaseDir())
	path := filepath.Join(header.SubDir(), header.Name())

	obj , f_err := file.ReadFile(fs,path,key,iv,header)
	if f_err != nil {
		err_msg.BadCookieValueResponse(w)
		return
	}

	upgrader := websocket.Upgrader{}
	ws_conn, ws_err := upgrader.Upgrade(w, r, nil)
	if ws_err != nil {
		err_msg.CantStartWSLoopResponse(w)
		return
	}
	fvs.loop(ws_conn, obj, key)

}
func (fvs *FileDataServ) loop(conn *websocket.Conn, obj *data.Object, key []byte) {

}
