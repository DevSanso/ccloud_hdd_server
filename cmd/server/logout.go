package main

import (
	"net/http"
	"strconv"
	"time"

	"ccloud_hdd_server/pkg/auth"
)


func logoutHandler(w http.ResponseWriter,r *http.Request) {
	cki,err := r.Cookie("session")
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	key, key_err := strconv.Atoi(cki.Value)
	if key_err != nil {
		w.Write([]byte(key_err.Error()))
		w.WriteHeader(400)
		return
	}
	err = auth.Logout(uint32(key))
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(400)
		return
	}
	expiration := time.Now()
	re_cki := &http.Cookie{Name : "session",Value : "",Expires : expiration}
	http.SetCookie(w,re_cki)
	w.WriteHeader(200)
}