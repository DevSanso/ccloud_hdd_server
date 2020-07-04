package worker

import (
	"ccloud_hdd_server/db"
	"encoding/json"
	"bytes"
	"path/filepath"
	"net/http"
	"os"
)


type ViewDir struct {}

type _FileInfo struct {
	Name string
	Path string
	IsDir bool
	Date string
	Size int64
}

func convertFI(finfo os.FileInfo,dir_path string) _FileInfo {
	return _FileInfo{
		Name : filepath.Base(finfo.Name()),
		Path : dir_path,
		IsDir : finfo.IsDir(),
		Date : finfo.ModTime().Format("2006-01-02"),
		Size : finfo.Size(),
	}
}


func (v *ViewDir) Do(w http.ResponseWriter,r *http.Request,next Worker) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	dirPath := r.URL.Query().Get("dir")

	if !db.ExistDir(dirPath) {
		w.WriteHeader(404)
		return
	}
	

	iList,err := db.GetDirFiles(dirPath)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var buf = new(bytes.Buffer)
	var encoder = json.NewEncoder(buf)

	for _,info := range(iList)  {
		if err = encoder.Encode(convertFI(info,dirPath));err != nil {break}
	}
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("content-type","application/json")
	w.Write(buf.Bytes())
	w.WriteHeader(200)

}

type FileIconView struct {}
