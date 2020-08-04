package ws_mux

import (
	"sync"
)

func serveMainRoutine(wp *sync.Pool, cp *sync.Pool) {
	t, ok := wp.Get().(*WsServiceCtx)
	if !ok {
		return
	}

	ptr, p_ok := t.Args.Value(CtxIndex).(*WsFileApiFormat)
	if !p_ok {
		panic("wsThread structure assertion error")
	}

	obj := t.Obj

	var b = make([]byte, obj.TokenSize())
	_, err := obj.ReadAt(b, ptr.Offset)

	if err != nil {
		errAndClose(t, err)
		cp.Put(t)
		return
	}

	if isOverDataRange(
		ptr.Offset, obj.DataSize(), obj.TokenSize()) {

		b = cutData(b, ptr.Offset, obj.DataSize())
		ptr.Size = obj.DataSize() - ptr.Offset
		ptr.IsExistNext = false
	}

	ptr.D = b

	err = t.Conn.WriteJSON(ptr)

	if err != nil {
		errAndClose(t, err)
		cp.Put(t)
		return
	}

	if !ptr.IsExistNext {
		cp.Put(t)
		return
	}

	go func() {
		ptr.Offset += int64(len(ptr.D))
		ptr.IsExistNext = true
		wp.Put(t)
	}()

}

func serveCloseRoutine(cp *sync.Pool) {
	t, ok := cp.Get().(*WsServiceCtx)
	if !ok {
		return
	}

	t.Obj.Close()
	t.Conn.Close()
	t.Wg.Done()
}
