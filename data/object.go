package data

import "context"

const TokenSize = 4096

type Object struct {
	key []byte
}

type DataMeta struct {
	Name string
	Offset int64
	Size int64
	FullSize int64
	B []byte
}

func (o *Object) ReadChan() (<-chan DataMeta, context.CancelFunc) {

}
func (o *Object) WriteChan() chan DataMeta {

}

func (o *Object)Size() int64 {

}