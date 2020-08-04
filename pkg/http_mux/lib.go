package http_mux

import (
	"net/http"
	"ccloud_hdd_server/pkg/worker"
)


type httpMux struct {}

func (*httpMux)ServeHTTP(w http.ResponseWriter,r *http.Request) {
	r
}

