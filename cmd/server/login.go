package main

import (
	"net/http"
	"context"
	"time"

	"ccloud_hdd_server/pkg/auth"
	"ccloud_hdd_server/pkg/get_db"
	"ccloud_hdd_server/pkg/use_hash"
)



func loginHandler(w http.ResponseWriter,r *http.Request) {
	if r.Method != "POST" {
		w.Write([]byte("not post method"))
		w.WriteHeader(400)
		return
	}
	conn,err := get_db.GetDbConn(context.Background())
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}
	defer conn.Close()
	r.ParseForm()
	passwd := use_hash.Sum([]byte(r.Form.Get("passwd")))

	key,auth_err := auth.Login(conn,passwd)
	if auth_err != nil {
		w.Write([]byte(auth_err.Error()))
		w.WriteHeader(500)
		return
	}
	expiration := time.Now().Add(24 * time.Hour)
	cki := &http.Cookie{Name : "session",Value : string(key),Expires : expiration}
	http.SetCookie(w,cki)
	w.WriteHeader(200)
}