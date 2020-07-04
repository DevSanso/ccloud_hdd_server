package object

import (
	"context"
	"crypto/cipher"
	"errors"
	"github.com/spf13/afero"

	"github.com/mattn/go-sqlite3"
)

const TokenSize = 4096
var ClosedChanErr = errors.New("channel close")
var CantUseMethodErr = errors.New("this method not allow")
type Object struct {
	f afero.File
}

type ObjChannel struct {
	channel chan []byte
	wFlag bool //ture : write ,false : read
	cancel context.CancelFunc
	err map[int]error
}
func (ch *ObjChannel)appendErr(key int,err error) {
	ch.err[key] = err
}

func (ch *ObjChannel)FreeErr(key int) {
	delete(ch.err,key)
}
func (ch *ObjChannel)AsyncRead() ([]byte) {
	if ch.wFlag {return nil}

	var data []byte = nil
	select {
	case buf,ok := <-ch.channel:
		if !ok {return nil}
		data = buf
	}
	return data
}

func (ch *ObjChannel)AsyncWrite(b []byte) ([]byte,error) {
	if !ch.wFlag {return nil,CantUseMethodErr}
	ch.channel <- b

}

func NewObject(f afero.File) Object {
	return Object{f}
}

func (o *Object)MakeReadChan(cryp cipher.Block) (<-chan []byte,context.CancelFunc,error) {
	stat,err := o.f.Stat()
	if err != nil {return nil,nil,err}

	channel := make(chan []byte)
	ctx,cancel := context.WithCancel(context.Background())

	go func() {
		var max = stat.Size()
		var offset int64 = 0
		var buf = make([]byte,TokenSize)

		for offset < max{
			select {
			case <-ctx.Done():
				break
			default:
				o.f.ReadAt(buf,offset)
			}
		}


		close(channel)
	}()
	

	return channel,cancel,nil
}





