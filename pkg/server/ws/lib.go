package ws

import (
	"sync"

	"github.com/gorilla/websocket"

	"ccloud_hdd_server/pkg/data"
)

type WsServiceCtx struct {
	Conn *websocket.Conn
	Obj  *data.Object

	Wg   *sync.WaitGroup
	Args interface{}
}

func (wlc *WsServiceCtx) setWg(wg *sync.WaitGroup) { wlc.wg = wg }

type Thread interface {
	Push(loop *WsServiceCtx)
}
