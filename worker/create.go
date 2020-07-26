package worker

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"database/sql"

	"github.com/gorilla/websocket"
	"github.com/spf13/afero"

	"ccloud_hdd_server/auth"
	"ccloud_hdd_server/data"
	"ccloud_hdd_server/db_sql"
	db_user "ccloud_hdd_server/db_sql/user"
	"ccloud_hdd_server/file"
	"ccloud_hdd_server/get_db"

	pkg_internal "ccloud_hdd_server/worker/internal"
)

const DataToKenSize = 4096
type CreateFileWork struct {}

func (cfw *CreateFileWork) Do(w http.ResponseWriter, r *http.Request, key []byte) {
	var upgrade  = websocket.Upgrader{}
	ws_conn,ws_err := upgrade.Upgrade(w,r,nil)
	if ws_err != nil {pkg_internal.CantStartWSLoopResponse(w);return}
	defer ws_conn.Close()

	_,db_conn,iv,base_path,setting_err :=authAndGetInfo(r)
	if setting_err != nil {pkg_internal.RawErrorResponse(w,setting_err,400)}

	dir := r.Form.Get("subDir")
	name := r.Form.Get("name")

	path := dir+"/"+name
	base_fs := afero.NewBasePathFs(afero.NewOsFs(),base_path)
	obj,obj_err := file.CreateEmptyFile(base_fs,path,key,iv,DataToKenSize)
	if obj_err != nil {pkg_internal.CantCreateDataResponse(w);return}
	
	ws_loop_err := cfw.uploadFileToWs(ws_conn,db_conn,obj)
	if ws_loop_err != nil {pkg_internal.CantStartWSLoopResponse(w);return}

	
}

func (cfw *CreateFileWork)uploadFileToWs(ws *websocket.Conn,c *sql.Conn,o *data.Object) error {
	return nil
}

type CreateDirWork struct {}
func (cdw *CreateDirWork) Do(w http.ResponseWriter, r *http.Request, key []byte) {
	if r.Method != "POST" {pkg_internal.BadMethodResponse(w);return}
	
	ck,ck_err :=r.Cookie("session")
	if ck_err != nil {pkg_internal.BadCookieValueResponse(w);return}
	

	using_key,atol_err := strconv.Atoi(ck.Value)
	if atol_err != nil {pkg_internal.BadCookieValueResponse(w);return}

	using_key,ck_err = auth.GetUesrId(uint32(using_key))
	if ck_err != nil {pkg_internal.NotLoginResponse(w);return}

	conn,db_err := get_db.GetDbConn(context.Background())
	if db_err != nil {pkg_internal.CantConnectDbResponse(w);return}
	defer conn.Close()

	var iv []byte
	iv,db_err = db_user.GetUserIv(conn,using_key)
	if db_err != nil {pkg_internal.CantSearchDataResponse(w);return}
	
	base_id,base_err := db_user.GetBasePathId(conn,using_key)
	var  base_path string
	base_path,base_err = db_sql.GetBasePath(conn,base_id)
	if base_err != nil {pkg_internal.CantSearchDataResponse(w);return}

	r.ParseForm()

	dir := r.Form.Get("subDir")
	name := r.Form.Get("name")

	path, cryp_err := file.EncodeFilePath(key,iv,dir+"/"+name)
	if cryp_err != nil {pkg_internal.CantCreateDataResponse(w);return}
	base_fs := afero.NewBasePathFs(afero.NewOsFs(),base_path)
	err := base_fs.Mkdir(path,os.ModeDir)
	if err != nil {pkg_internal.CantCreateDataResponse(w);return}

	w.WriteHeader(200)
}