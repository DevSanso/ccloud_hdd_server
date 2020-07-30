package server

import (
	"context"
	"net"
	"net/http"
	"sync"

	"ccloud_hdd_server/pkg/data"
	ws_service "ccloud_hdd_server/pkg/server/ws"
	"ccloud_hdd_server/pkg/util"

	"github.com/gorilla/websocket"
	"github.com/scylladb/go-set/b64set"
)

const (
	WsUPLOAD = iota
	WsSER

	__WsEndPos__
)

const (
	UrlKeySize = 64
)

type WsServerHook interface {
	RequestWsService(r *WsRequest) []byte
}













type WsRequest struct {
	Ip       net.IP
	WsMethod int

	Obj  *data.Object
	Args interface{}
}

type wsServeMux struct {
	services map[int]ws_service.Thread
	urlSet   b64set.Set
	urlCtx   map[[64]byte]context.Context
}

func newWsServeMux() *wsServeMux {
	var wss = &wsServeMux{}
	wss.services[WsSER] = ws_service.NewServThread()
	wss.services[WsUPLOAD] = ws_service.NewUploadThread()
	return wss
}

func (wss *wsServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, code, err := wss.checkHttpVal(r)
	if err != nil {
		writeErrToRes(w, err)
		w.WriteHeader(code)
		return
	}
	remote_ip := ctx.Value("remote-ip").(net.IP)
	h, _, _ := net.SplitHostPort(r.RemoteAddr)

	if !remote_ip.Equal(net.ParseIP(h)) {
		writeErrToRes(w, NotMatchIpErr)
		w.WriteHeader(400)
		return
	}

	ws_method, m_ok := ctx.Value("ws-method").(int)
	if !m_ok && ws_method >= __WsEndPos__ {
		panic("server/ws_server ws_method no exist error")
	}

	var upgrade = websocket.Upgrader{}
	conn, ws_err := upgrade.Upgrade(w, r, nil)
	if ws_err != nil {
		writeErrToRes(w, ws_err)
		w.WriteHeader(400)
		return
	}
	lwg := new(sync.WaitGroup)
	lwg.Add(1)
	wss.awaitServiceCtx(ws_method, &ws_service.WsServiceCtx{
		Conn: conn,
		Obj:  ctx.Value("object").(*data.Object),
		Wg:   lwg,
		Args: ctx.Value("args"),
	})

}

func (wss *wsServeMux) checkHttpVal(r *http.Request) (context.Context, int, error) {
	url := []byte(r.URL.Query().Get("url"))
	var k64 [64]byte
	l := copy(k64[:], url)
	if l != 64 {
		return nil, 500, InternalServerErr
	}
	if !wss.urlSet.Has(k64) {
		return nil, 400, NotExistUrlInWsErr
	}
	ctx := wss.urlCtx[k64]
	if ctx == nil {
		return nil, 500, InternalServerErr
	}
	return ctx, 0, nil
}
func (wss *wsServeMux) awaitServiceCtx(method int, wlc *ws_service.WsServiceCtx) {
	wss.services[method].Push(wlc)
	wlc.Wg.Wait()
}

func (wss *wsServeMux) RequestWsService(r *WsRequest) (key [64]byte) {
	key = wss.makeKey()
	ctx := wpm.MakeReqContext(r)
	wss.urlCtx[key] = ctx
	return
}

func (wss *wsServeMux) makeKey() [64]byte {
	key := util.MakeBytes(64)
	var key64 [64]byte
	copy(key64[:], key)

	for wss.urlSet.Has(key64) {
		key = util.MakeBytes(64)
		copy(key64[:], key)
	}
	wss.urlSet.Add(key64)

	return key64
}

func (wss *wsServeMux) freeKey(k [64]byte) {
	if !wss.urlSet.Has(k) {
		return
	}
	wss.urlSet.Remove(k)
}

type wsPrivateMethod struct{}

func (wsPrivateMethod) MakeReqContext(r *WsRequest) context.Context {
	ip := context.WithValue(context.Background(), "remote-ip", r.Ip)
	obj := context.WithValue(ip, "object", r.Obj)
	ws_method := context.WithValue(obj, "ws-method", r.WsMethod)
	args := context.WithValue(ws_method, "args", r.Args)
	return args
}

var wpm wsPrivateMethod
