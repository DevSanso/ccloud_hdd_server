package ws_mux



import (
	"net/http"
	"errors"
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"ccloud_hdd_server/pkg/data"
)

var (
	NoMatchUrlLenErr = errors.New("no match url length")
	NotExistUrlInWsErr = errors.New("no exist access url")
	InternalServerErr = errors.New("internal error")
	NotMatchIpErr = errors.New("not match ip error")
)



func writeErrToRes(w http.ResponseWriter,err error) {
	w.Header().Set("content-type","text/plain")
	w.Write([]byte(err.Error()))
}



func WsServerHooking() WsServerHook {
	return ws_mux
}

var (
	CantConvertCtxValueErr = errors.New("CantConvertCtxValueErr")
)

const (
	CtxIndex = "ws-api-buffer"
)

type WsServiceCtx struct {
	Conn *websocket.Conn
	Obj  *data.Object

	Wg   *sync.WaitGroup
	Args context.Context
}

func (wlc *WsServiceCtx) setWg(wg *sync.WaitGroup) { wlc.Wg = wg }

type Thread interface {
	Push(loop *WsServiceCtx)
}
type MainFunc func(wait *sync.Pool, close *sync.Pool)
type CloseFunc func(close *sync.Pool)

type WsFileApiFormat struct {
	Name string
	//데이터 사이즈
	Size        int64
	Offset      int64
	IsExistNext bool
	D           []byte
}

type wsThread struct {
	waitLoopP sync.Pool
	closePool sync.Pool

	mainRoutine  MainFunc
	closeRoutine CloseFunc
}
type ServiceCtx struct {
	Name   string
	Offset int64

	format *WsFileApiFormat
}

func newServThread() Thread {
	return newWsThread(
		serveMainRoutine,
		serveCloseRoutine,
	)
}
func newUploadThread() Thread {
	return newWsThread(
		uploadMainRoutine,
		uploadCloseRoutine,
	)
}

func newWsThread(mainF MainFunc, closeF CloseFunc) *wsThread {
	var s = &wsThread{
		sync.Pool{},
		sync.Pool{},
		mainF,
		closeF,
	}
	go func() {
		for {
			s.mainRoutine(&s.waitLoopP, &s.closePool)
			s.closeRoutine(&s.closePool)
		}
	}()
	return s
}
func (sl *wsThread) Push(wsCtx *WsServiceCtx) {
	sl.waitLoopP.Put(wsCtx)
}
func (sl *wsThread) makeCtx(wsCtx *WsServiceCtx) error {
	sl.waitLoopP.Put(wsCtx)
	return nil
}

func errAndClose(t *WsServiceCtx, err error) {
	t.Conn.WriteMessage(websocket.CloseMessage,
		[]byte(err.Error()))
}
func isOverDataRange(off, dataSize int64, tokenSize int) bool {
	return off+int64(tokenSize) > dataSize
}

func cutData(d []byte, off int64, dataSize int64) []byte {
	return d[:dataSize-off]
}

