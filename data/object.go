package data

import (
	"io"
	"errors"

	"github.com/spf13/afero"

	"ccloud_hdd_server/cryp"
)





var (
	OverflowSizeError = errors.New("OverflowSizeError")

)

type Object struct {
	f afero.File
	key []byte
	iv []byte

	tokenSize int
	dataSize int64
}

func NewObject(f afero.File,key []byte,iv []byte,tokenSize int,dataSize int64) (*Object,error) {
	
	return &Object{f,key,iv,tokenSize,dataSize},nil
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
	var read_len int
	read_len,err = r_at.ReadAt(buf,offset)
	if err != nil {return 0,err}

	
	decoder,cipher_err := cryp.NewDecoder(o.key,o.iv)
	if cipher_err != nil {return 0,cipher_err}

	cipher_err = decoder.Decrypt(buf,b)
	if cipher_err != nil {return 0,cipher_err}
	return read_len,nil
}


func(o *Object)WriteAt(b []byte,offset int64)(int,error) {
	f_size,err := o.FileSize()

	if err != nil {return 0,err}

	if err = o.checkOffset(offset,f_size); err != nil {return 0,err} 

	buf := make([]byte,o.tokenSize)
	w_at := io.WriterAt(o.f)
	

	encoder,cipher_err := cryp.NewEncoder(o.key,o.iv)
	if cipher_err != nil {return 0,cipher_err}

	cipher_err = encoder.Encrypt(b,buf)
	if cipher_err != nil {return 0,cipher_err}
	var write_len int
	write_len,err = w_at.WriteAt(buf,offset);
	if err != nil {return 0,err}

	if same,f_size := o.isSameDataSize();!same {o.updateDataSize(f_size)}
	return write_len,nil
}
func (o *Object)isSameDataSize() (bool,int64) {
	info,err := o.f.Stat()
	if err != nil {panic(err)}
	return info.Size() == o.dataSize,info.Size()
}
func (o *Object)updateDataSize(size int64) {o.dataSize = size}


func(o *Object)Seek(offset int64, whence int) (int64,error) {
	return o.f.Seek(offset,whence)
}

func(o *Object)Close() error {
	o.key = nil
	return o.f.Close()
}

func (o *Object)checkOffset(offset int64,fSize int64) error {
	if offset + int64(o.tokenSize) > fSize {
		return OverflowSizeError
	}
	if offset % int64(o.tokenSize-1) != 0 {
		return OverflowSizeError
	}
	
	return nil
	
}



func(o *Object)FileSize() (int64,error) {
	i,err := o.f.Stat()
	if err != nil {return 0,err}
	return i.Size(),nil 
}

func (o *Object)DataSize() int64 {return o.dataSize}
func (o *Object)TokenSize() int {return o.tokenSize}

