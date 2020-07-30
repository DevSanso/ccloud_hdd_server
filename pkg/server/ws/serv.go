package ws

import (
	"sync"
)

func serveMainRoutine(wp *sync.Pool, cp *sync.Pool) {
	t, ok := wp.Get().(*WsServiceCtx)
	if !ok {
		return
	}

	slm_ptr, p_ok := t.Args.(*ServiceCtxValue)
	if !p_ok {
		panic("wsThread structure assertion error")
	}

	obj := t.Obj

	var res_format = slm_ptr.format

	var b = make([]byte, obj.TokenSize())
	_, err := obj.ReadAt(b, slm_ptr.Offset)

	if err != nil {
		errAndClose(t, err)
		cp.Put(t)
		return
	}

	if isOverDataRange(
		slm_ptr.Offset, obj.DataSize(), obj.TokenSize()) {

		b = cutData(b, slm_ptr.Offset, obj.DataSize())
		res_format.Size = obj.DataSize() - slm_ptr.Offset
		res_format.IsExistNext = false
	}

	res_format.D = b
	res_format.Offset = slm_ptr.Offset

	err = t.Conn.WriteJSON(&res_format)

	if err != nil {
		errAndClose(t, err)
		cp.Put(t)
		return
	}

	if !res_format.IsExistNext {
		cp.Put(t)
		return
	}

	go func() {
		slm_ptr.Offset += int64(len(res_format.D))
		slm_ptr.format.IsExistNext = true
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
