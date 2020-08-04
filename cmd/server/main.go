package main

import (
	"sync"
	"strconv"
	"net/http"

	"ccloud_hdd_server/pkg/get_db"
)


func HttpServer(addr string) *http.Server {
	http_server := &http.Server{Addr : addr ,Handler : nil}
	return http_server
}

func WsServer(addr string) *http.Server {
	ws_server := &http.Server{Addr : addr ,Handler : nil}
	return ws_server
}
func sAddr(host string,port int) string {
	
	return host + ":" + strconv.Itoa(port)
}

func main() {
	cfg := GetConfig()
	dberr := get_db.OpenDb("mysql","")
	if dberr != nil {panic(dberr)}
	
	certf,keyf := cfg.CertFile,cfg.KeyFile

	var ws_err,err error
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err = HttpServer(sAddr(cfg.Http.Host,cfg.Http.Port)).ListenAndServeTLS(certf,keyf)
		wg.Done()
	}()
	go func() {
		ws_err = WsServer(sAddr(cfg.Ws.Host,cfg.Ws.Port)).ListenAndServeTLS(certf,keyf)
		wg.Done()
	}()
	wg.Wait()
	if ws_err != nil && err != nil {
		panic("ws server err :" + ws_err.Error()+ 
		" " + "http server err : "+err.Error())
	}
	
}