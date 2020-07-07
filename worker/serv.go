package worker

import (
	"ccloud_hdd_server/db"
	"ccloud_hdd_server/data"
	"encoding/json"
	"bytes"
	"path/filepath"
	"net/http"
	"github.com/gorilla/websocket"
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

func (fvs *FileDataServ)Do(w http.ResponseWriter,r *http.Request,key []byte) {
	fileName := r.URL.Query().Get("file")
	dirPath := r.URL.Query().Get("dir")
	if !db.ExistFile(dirPath,fileName) {
		w.WriteHeader(404)
		return
	}

	obj,objErr := data.GetObject(key,mergePath(dirPath,fileName))
	if objErr != nil {
		w.WriteHeader(500)
		return
	}

	upgrader := websocket.Upgrader{}
	conn ,err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	fvs.loop(conn,obj)
	
}
func(fvs *FileDataServ)loop(conn *websocket.Conn,obj *data.Object) {
	channel,cancel := obj.ReadChan()
	conn.SetCloseHandler(func(code int,text string) error {

	})

	for i := int64(0); i < obj.Size(); {
		data,ok:=<-channel
		if !ok {
			
			break
		}
		if err := conn.ReadJSON(data);err != nil {cancel();break}
	}
}


