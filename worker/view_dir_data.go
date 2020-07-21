package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/afero"

	"ccloud_hdd_server/auth"
	"ccloud_hdd_server/db_sql"
	db_user "ccloud_hdd_server/db_sql/user"
	"ccloud_hdd_server/file"
	"ccloud_hdd_server/get_db"
	pkg_internal "ccloud_hdd_server/worker/internal"
)

type ViewDir struct{}

func convertFI(finfo os.FileInfo) _FileMeta {
	return _FileMeta{
		Name:  finfo.Name(),
		IsDir: finfo.IsDir(),
		Date:  finfo.ModTime().Format("2006-01-02"),
		Size:  finfo.Size(),
	}
}

type _FileMeta struct {
	Name  string
	Size  int64
	Date  string
	IsDir bool
}
type _FileList struct {
	Dir       string
	FileInfos []_FileMeta
}

func (v *ViewDir) makeFs(key int) (afero.Fs, error) {
	conn, err := get_db.GetDbConn(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	base_key, sql_err := db_user.GetBasePathId(conn, key)
	if sql_err != nil {
		return nil, sql_err
	}
	var p string
	p, sql_err = db_sql.GetBasePath(conn, base_key)
	if sql_err != nil {
		return nil, sql_err
	}

	return afero.NewBasePathFs(afero.NewOsFs(), p), nil
}

func (v *ViewDir) Do(w http.ResponseWriter, r *http.Request, _ []byte) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}

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
	dir := r.URL.Query().Get("Dir")

	fs, fs_err := v.makeFs(using_key)
	if fs_err != nil {
		pkg_internal.CantSearchDataResponse(w)
		return
	}

	iList, err := file.GetFileList(fs, dir)
	if err != nil {
		pkg_internal.CantSearchDataResponse(w)
		return
	}

	var buf = new(bytes.Buffer)
	var encoder = json.NewEncoder(buf)

	var res = _FileList{}
	for _, i := range iList {
		res.FileInfos = append(res.FileInfos, convertFI(i))
	}
	err = encoder.Encode(res)
	if err != nil {
		pkg_internal.CantCreateDataResponse(w)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(buf.Bytes())
	w.WriteHeader(200)
}
