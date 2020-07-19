package internal


import (
	"net/http"
)


func NotLoginResponse(w http.ResponseWriter) {
	w.Header().Set("content-type","text/plain")
	w.Write([]byte("not login"))
	w.WriteHeader(400)
}

func BadCookieValueResponse(w http.ResponseWriter) {
	w.Header().Set("content-type","text/plain")
	w.Write([]byte("bad cookie"))
	w.WriteHeader(400)
}
func CantSearchDataResponse(w http.ResponseWriter) {
	w.Header().Set("content-type","text/plain")
	w.Write([]byte("can't search data"))
	w.WriteHeader(404)
}
func CantStartWSLoopResponse(w http.ResponseWriter) {
	w.Header().Set("content-type","text/plain")
	w.Write([]byte("can't websocket run"))
	w.WriteHeader(500)
}

func CantCreateDataResponse(w http.ResponseWriter) {
	w.Header().Set("content-type","text/plain")
	w.Write([]byte("can't create data"))
	w.WriteHeader(500)
}