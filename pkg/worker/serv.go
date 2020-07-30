package worker

import (
	"bytes"
	
	"context"
	"database/sql"
	"encoding/json"
	"net"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/spf13/afero"

	"ccloud_hdd_server/pkg/auth"
	"ccloud_hdd_server/pkg/data"
	"ccloud_hdd_server/pkg/db_sql"
	"ccloud_hdd_server/pkg/file"
	"ccloud_hdd_server/pkg/get_db"
	servers "ccloud_hdd_server/pkg/server"
	ws_service "ccloud_hdd_server/pkg/server/ws"
	db_user "ccloud_hdd_server/pkg/db_sql/user"
	pkg_internal "ccloud_hdd_server/pkg/worker/internal"
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

func (fvs *FileDataServ) getUserBaseId(c *sql.Conn, key int) (int, error) {
	return db_user.GetBasePathId(c, key)

}
func (fvs *FileDataServ) getHeader(c *sql.Conn, path string) (*db_sql.Header, error) {
	return db_sql.LoadHeader(c, path)
}

func (fvs *FileDataServ) Do(w http.ResponseWriter, r *http.Request, key []byte) {
	ck, ck_err := r.Cookie("session")
	if ck_err != nil {
		pkg_internal.NotLoginResponse(w)
		return
	}
	var using_key int
	using_key, ck_err = strconv.Atoi(ck.Value)

	if ck_err != nil {
		pkg_internal.BadCookieValueResponse(w)
		return
	}
	if using_key, ck_err = auth.GetUesrId(uint32(using_key)); ck_err != nil {
		pkg_internal.BadCookieValueResponse(w)
		return
	}

	file_name := r.URL.Query().Get("file")
	dir_path := r.URL.Query().Get("dir")
	c, db_err := get_db.GetDbConn(context.Background())
	if db_err != nil {
		pkg_internal.CantConnectDbResponse(w)
		return
	}
	var header *db_sql.Header
	header, db_err = fvs.getHeader(c, dir_path+"/"+file_name)
	if db_err != nil {
		pkg_internal.CantSearchDataResponse(w)
		return
	}
	var iv []byte
	iv, db_err = db_user.GetUserIv(c, using_key)
	if db_err != nil {
		pkg_internal.CantSearchDataResponse(w)
		return
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), header.BaseDir())

	obj, f_err := file.ReadFile(fs,
		filepath.Join(header.SubDir(), header.Name()), key, iv, header)
	if f_err != nil {
		pkg_internal.BadCookieValueResponse(w)
		return
	}
	h, _, _ := net.SplitHostPort(r.RemoteAddr)
	r_ip := net.ParseIP(h)
	url_key := fvs.wsHook(r_ip,obj,header)
	
	w.Header().Set("content-type","text/plain")
	w.Write(url_key[:])
	w.WriteHeader(204)

}

type wsFileFormat struct {
	Name        string
	Size        int64
	Offset      int64
	IsExistNext bool
	D           []byte
}

func (fvs *FileDataServ)wsHook(ip net.IP,obj *data.Object, h *db_sql.Header) [64]byte{
	var req = new(servers.WsRequest)
	req.Ip = ip
	req.Obj = obj
	req.WsMethod = servers.WsSER

	req.Args = context.WithValue(context.Background(),
	 ws_service.CtxIndex,&ws_service.WsFileApiFormat{
		h.Name(),
		h.Size(),
		0,
		true,
		[]byte(""),
	})

	hook := servers.WsServerHooking()
	return hook.RequestWsService(req)
	
}


func (fvs *FileDataServ) serveWsErr(conn *websocket.Conn, err error) {
	conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
}
