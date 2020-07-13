package data

import (
	"crypto/aes"
	"os"
)


const FileMode = os.FileMode(0)
const (
	AES256 = iota + 1
)
type Object struct {
	*os.File
	key []byte
	cryt int
	tokenSize int64
	dataSize int64
}

func NewObject(path string,key []byte,cryt int,tokenSize int64,flag int) (Object,error) {
	f,err := os.OpenFile(path,flag,FileMode)
	if err != nil {return Object{nil,nil,0,0},err}
	return Object{f,key,cryt,tokenSize},nil
}


func (o *Object)Read(b []byte)(int,error) {

}

func (o *Object)GetKey() []byte {
	return obj.key
}

func(o *Object)ReadAt()(b []byte,offset int64)(int,error) {

}
func(o *Object)Write(b []byte)(int,error) {

}

func(o *Object)WriteAt(b []byte,offset int64)(int,error) {

}

func(o *Object)Seek(offset int64, whence int) error {

}



func(o *Object)FileSize() (int64,error) {
	i,err := o.File.Stat()
	if err != nil {return 0,err}
	return i.Size()
}

func (o *Object)DataSize() int64 {return o.dataSize}
func (o *Object)TokenSize() int64 {return o.tokenSize}

