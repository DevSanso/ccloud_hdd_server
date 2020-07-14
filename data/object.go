package data

import (
	"crypto/aes"
	"crypto/des"
	"os"
	"crypto/cipher"
	"errors"
)


const FileMode = os.FileMode(0)
const (
	AES256 = iota + 1
	DES
)

var NoReleaseCryptoError = errors.New("NoReleaseCryptoError")

type Object struct {
	f *os.File
	key []byte
	cryt byte
	tokenSize int64
	dataSize int64
}

func NewObject(path string,key []byte,cryt byte,tokenSize,dataSize int64,flag int) (Object,error) {
	f,err := os.OpenFile(path,flag,FileMode)
	if err != nil {return Object{nil,nil,0,0,0},err}
	return Object{f,key,cryt,tokenSize,dataSize},nil
}


func (o *Object)Read(b []byte)(int,error) {

}

func (o *Object)GetKey() []byte {
	return o.key
}

func(o *Object)ReadAt(b []byte,offset int64) (int,error) {
}
func(o *Object)Write(b []byte)(int,error) {
	
}

func(o *Object)WriteAt(b []byte,offset int64)(int,error) {

}

func(o *Object)Seek(offset int64, whence int) error {
	return o.f.Seek(offset,whence)
}

func(o *Object)Close() error {

}

func getCipher(key []byte,cryt byte) (cipher.Block,error) {
	switch cryt {
	case AES256:
		return aes.NewCipher(key)
	case DES:
		return des.NewCipher(key)
	default:
		return nil,NoReleaseCryptoError
	}
}


func(o *Object)FileSize() (int64,error) {
	i,err := o.f.Stat()
	if err != nil {return 0,err}
	return i.Size() 
}

func (o *Object)DataSize() int64 {return o.dataSize}
func (o *Object)TokenSize() int64 {return o.tokenSize}

