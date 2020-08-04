package main

import (
	"net/http"
	"strconv"

	"ccloud_hdd_server/pkg/worker"
	"ccloud_hdd_server/pkg/auth"
)


func workHandler(w http.ResponseWriter,r *http.Request) {
	cki,err := r.Cookie("session")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	var s_key int
	s_key,err = strconv.Atoi(cki.Value)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	session_key := uint32(s_key)
	
	key := auth.GetUserCryptoKey(session_key)
	target := ""
	work := worker.GetWoker(target)
	work.Do(w,r,key)
}