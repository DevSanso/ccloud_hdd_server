package ws

import (
	"sync"
	"github.com/gorilla/websocket"

	"ccloud_hdd_server/pkg/data"
)

type WsLoopCtx struct {
	Conn   *websocket.Conn
	Obj    *data.Object

	Args interface{}
	
}

type Thread interface{
	Push(loop *WsLoopCtx,wg sync.WaitGroup)
}