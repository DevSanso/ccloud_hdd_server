package data

import (
	"crypto/aes"
	"crypto/des"
	"os"
	"io"
	"crypto/cipher"
	"errors"
)


const FileMode = os.FileMode(0)
const (
	AES256 = iota + 1
	DES
)

var (
	NoReleaseCryptoError = errors.New("NoReleaseCryptoError")
	OverflowSizeError = errors.New("OverflowSizeError")
	NotMatchArraySizeError = errors.New("NotMatchArraySizeError")
)
type Object struct {
	f *os.File
	key []byte
	iv []byte
	cryt byte
	tokenSize int64
	dataSize int64
}

func newObject(path string,key []byte,iv []byte,cryt byte,tokenSize,dataSize int64,flag int) (Object,error) {
	f,err := os.OpenFile(path,flag,FileMode)
	if err != nil {return Object{nil,nil,nil,0,0,0},err}
	return Object{f,key,iv,cryt,tokenSize,dataSize},nil
}

func (o *Object)GetKey() []byte {
	return o.key
}

func(o *Object)ReadAt(b []byte,offset int64) (int,error) {
	f_size,err := o.FileSize()

	if err != nil {return 0,err}

	if err = o.checkOffset(offset,f_size); err != nil {return 0,err} 

	buf := make([]byte,o.tokenSize)
	r_at := io.ReaderAt(o.f)
	if _,err = r_at.ReadAt(buf,offset);err != nil {return 0,err}

	block,cipher_err := getBlock(o.key,o.cryt)
	if cipher_err != nil {return 0,cipher_err}

	decrypt(block,buf,b,o.iv)
	return 0,nil
}


func(o *Object)WriteAt(b []byte,offset int64)(int,error) {
	f_size,err := o.FileSize()

	if err != nil {return 0,err}

	if err = o.checkOffset(offset,f_size); err != nil {return 0,err} 

	buf := make([]byte,o.tokenSize)
	r_at := io.WriterAt(o.f)
	if _,err = r_at.WriteAt(buf,offset);err != nil {return 0,err}

	block,cipher_err := getBlock(o.key,o.cryt)
	if cipher_err != nil {return 0,cipher_err}

	encrypt(block,buf,b,o.iv)
	return 0,nil
}

func decrypt(block cipher.Block,src,dst []byte,iv []byte) {
	cipher.NewCFBDecrypter(block,iv).XORKeyStream(dst,src)
}

func encrypt(block cipher.Block,src,dst []byte,iv []byte) {
	cipher.NewCFBEncrypter(block,iv).XORKeyStream(dst,src)
}

func(o *Object)Seek(offset int64, whence int) (int64,error) {
	return o.f.Seek(offset,whence)
}

func(o *Object)Close() error {
	o.key = nil
	return o.f.Close()
}

func (o *Object)checkOffset(offset int64,fSize int64) error {
	if offset + o.tokenSize > fSize {
		return OverflowSizeError
	}
	if offset % o.tokenSize-1 != 0 {
		return OverflowSizeError
	}
	
	return nil
	
}


func getBlock(key []byte,cryt byte) (cipher.Block,error) {
	switch cryt {
	case AES256:
		if len(key) != 32 {return nil,NotMatchArraySizeError}
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
	return i.Size(),nil 
}

func (o *Object)DataSize() int64 {return o.dataSize}
func (o *Object)TokenSize() int64 {return o.tokenSize}

