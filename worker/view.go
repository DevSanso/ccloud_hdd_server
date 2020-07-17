package worker

import (
	"ccloud_hdd_server/data"
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



func (v *ViewDir) Do(w http.ResponseWriter,r *http.Request,key []byte) {
	if r.Method != "GET" {
		w.WriteHeader(404)
		return
	}
	dirPath := r.URL.Query().Get("dir")
	
	if !data.ExistDir(dirPath) {
		w.WriteHeader(404)
		return
	}
	

	iList,err := data.GetDirFiles(dirPath)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var buf = new(bytes.Buffer)
	var encoder = json.NewEncoder(buf)
	
	var res struct {
		Arr []_FileInfo
	}
	for _,i := range(iList)  {
		res.Arr = append(res.Arr,convertFI(i,dirPath))
	}
	err = encoder.Encode(res.Arr)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("content-type","application/json")
	w.Write(buf.Bytes())
	w.WriteHeader(200)
}


