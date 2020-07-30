package ws

import (
	"sync"
	"database/sql"


	"ccloud_hdd_server/pkg/db_sql"
)

func uploadMainRoutine(wp *sync.Pool, cp *sync.Pool) {
	t, ok := wp.Get().(*WsServiceCtx)
	if !ok {
		return
	}

	ptr, p_ok := t.Args.Value(CtxIndex).(*WsFileApiFormat)
	if !p_ok {
		panic("wsThread structure assertion error")
	}
	err := t.Conn.ReadJSON(ptr)
	if err != nil {
		errAndClose(t, err)
		cp.Put(t)
		return
	}
	obj := t.Obj
	var write_len int
	write_len, err = obj.WriteAt(ptr.D, ptr.Offset)
	if err != nil {
		errAndClose(t, err)
		cp.Put(t)
		return
	}
	if !ptr.IsExistNext {
		db_conn := t.Args.Value("db-conn").(*sql.Conn)
		name := t.Args.Value("origin-name").(string)
		cfg := t.Args.Value("cfg").(*db_sql.ObjectConfig)
		cfg.Size = ptr.Offset + int64(write_len)
		
		go func() {
			_,err := db_sql.CreateHeader(db_conn,name,cfg)
			if err != nil {panic(err)}
		}()

		cp.Put(t)
		return
	}
	go func() {
		ptr.Offset += int64(write_len)
		ptr.IsExistNext = true
		wp.Put(t)
	}()

}

func uploadCloseRoutine(cp *sync.Pool) {
	t, ok := cp.Get().(*WsServiceCtx)
	if !ok {
		return
	}
	t.Obj.Close()
	t.Conn.Close()
	t.Wg.Done()
}
