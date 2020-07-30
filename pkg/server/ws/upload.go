package ws

import (
	"sync"
)

func uploadMainRoutine(wp *sync.Pool, cp *sync.Pool) {
	t, ok := wp.Get().(*WsServiceCtx)
	if !ok {
		return
	}

	slm_ptr, p_ok := t.Args.(*ServiceCtxValue)
	if !p_ok {
		panic("wsThread structure assertion error")
	}
	err := t.Conn.ReadJSON(slm_ptr.format)
	if err != nil {
		errAndClose(t, err)
		cp.Put(t)
		return
	}
	obj := t.Obj
	var write_len int
	write_len, err = obj.WriteAt(slm_ptr.format.D, slm_ptr.Offset)
	if err != nil {
		errAndClose(t, err)
		cp.Put(t)
		return
	}
	if !slm_ptr.format.IsExistNext {
		cp.Put(t)
		return
	}
	go func() {
		slm_ptr.Offset += int64(write_len)
		slm_ptr.format.IsExistNext = true
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
