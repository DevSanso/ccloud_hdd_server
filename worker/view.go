package worker

import (
	"net/http"
)


type View struct {}


func (v *View) Do(w http.ResponseWriter,r *http.Request,next Worker) {
	if r.Method != "GET" {
		w.WriteHeader(400)
		return
	}
	path := r.URL.Query().Get("path")
}