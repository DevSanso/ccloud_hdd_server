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
	fvs.wsHook(r_ip,obj,header)
	
	upgrader := websocket.Upgrader{}
	ws_conn, ws_err := upgrader.Upgrade(w, r, nil)
	if ws_err != nil {
		pkg_internal.CantStartWSLoopResponse(w)
		return
	}
	fvs.loop(ws_conn, obj, header)
}

type wsFileFormat struct {
	Name        string
	Size        int64
	Offset      int64
	IsExistNext bool
	D           []byte
}

func (fvs *FileDataServ)wsHook(ip net.IP,obj *data.Object, h *db_sql.Header) {
	var req servers.WsRequest
	req.Ip = ip
	req.Obj = obj
	req.WsMethod = servers.WsSER
	
	hook := servers.WsServerHooking()
}

func (fvs *FileDataServ) loop(conn *websocket.Conn,
	obj *data.Object, h *db_sql.Header) {
	defer conn.Close()
	defer obj.Close()

	offset := int64(0)
	var buf = bytes.Buffer{}
	encode := json.NewEncoder(&buf)

	var data_buf = make([]byte, obj.TokenSize())
	var res_format = wsFileFormat{
		Name: h.Name(),
		Size: obj.DataSize(),
	}

	var isOverDataRange = func() bool {
		return offset+int64(obj.TokenSize()) > res_format.Size
	}
	var cutData = func() []byte {
		return data_buf[:offset+int64(obj.TokenSize())-res_format.Size]
	}

	var is_next = true
	for _, err := obj.ReadAt(data_buf, 0); err == nil && is_next; _, err = obj.ReadAt(data_buf, offset) {
		if isOverDataRange() {
			data_buf = cutData()
			is_next = false
		}

		res_format.Offset = offset
		res_format.IsExistNext = is_next
		res_format.D = data_buf

		if err = encode.Encode(&res_format); err != nil {
			fvs.serveWsErr(conn, err)
			break
		}

		if err = conn.WriteMessage(websocket.TextMessage, buf.Bytes()); err != nil {
			fvs.serveWsErr(conn, err)
			break
		}
		buf.Reset()
	}

	if !is_next {
		conn.WriteMessage(websocket.CloseMessage, []byte("done"))
	}

}

func (fvs *FileDataServ) serveWsErr(conn *websocket.Conn, err error) {
	conn.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
}
