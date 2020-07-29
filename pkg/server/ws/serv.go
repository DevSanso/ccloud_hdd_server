package ws

import (
	"sync"

	"github.com/gorilla/websocket"
)

type wsFileRes struct {
	Name string
	//데이터 사이즈
	Size        int64
	Offset      int64
	IsExistNext bool
	D           []byte
}
type ServeLoop struct {
	waitLoopP   sync.Pool
	closedLoopP sync.Pool

	isRunning bool
}
type ServeServiceCtxValue struct {
	Name   string
	Offset int64

	format *wsFileRes
}

func NewServeLoop() *ServeLoop {
	var s = new(ServeLoop)
	go func() {
		for s.isRunning {
			s.runThread()
			s.closeLoop()
		}
	}()
	return s
}
func (sl *ServeLoop) Push(wsCtx *WsServiceCtx) { sl.waitLoopP.Put(wsCtx) }

func (sl *ServeLoop) runThread() {

	t, ok := sl.waitLoopP.Get().(*WsServiceCtx)
	if !ok {
		return
	}

	slm_ptr, p_ok := t.Args.(*ServeServiceCtxValue)
	if !p_ok {
		panic("ServeLoop structure assertion error")
	}

	obj := t.Obj

	var res_format = slm_ptr.format

	var b = make([]byte, obj.TokenSize())
	_, err := obj.ReadAt(b, slm_ptr.Offset)

	if err != nil {
		sl.errAndClose(t, err)
		return
	}

	if sl.isOverDataRange(
		slm_ptr.Offset, obj.DataSize(), obj.TokenSize()) {

		b = sl.cutData(b, slm_ptr.Offset, obj.DataSize())
		res_format.Size = obj.DataSize() - slm_ptr.Offset
		res_format.IsExistNext = false
	}

	res_format.D = b
	res_format.Offset = slm_ptr.Offset

	err = t.Conn.WriteJSON(&res_format)

	if err != nil {
		sl.errAndClose(t, err)
		return
	}

	if !res_format.IsExistNext {
		sl.closedLoopP.Put(t)
		return
	}

	go func() {
		slm_ptr.Offset += int64(len(res_format.D))
		slm_ptr.format.IsExistNext = true
		sl.waitLoopP.Put(t)
	}()

}
func (sl *ServeLoop) closeLoop() {
	t, ok := sl.closedLoopP.Get().(*WsServiceCtx)
	if !ok {
		return
	}

	t.Obj.Close()
	t.Conn.Close()
	t.Wg.Done()
}
func (sl *ServeLoop) errAndClose(t *WsServiceCtx, err error) {
	t.Conn.WriteMessage(websocket.CloseMessage,
		[]byte(err.Error()))
	sl.closedLoopP.Put(t)
}
func (sl *ServeLoop) isOverDataRange(off, dataSize int64, tokenSize int) bool {
	return off+int64(tokenSize) > dataSize
}

func (sl *ServeLoop) cutData(d []byte, off int64, dataSize int64) []byte {
	return d[:dataSize-off]
}
