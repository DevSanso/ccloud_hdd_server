package object

import (
	"os"
	"github.com/spf13/afero"
	"path/filepath"
)


type Object struct {
	phy afero.Fs
}

var dataRootPath string

func getRootPath() string {return dataRootPath}
func objectPath(path string) string {return filepath.FromSlash(dataRootPath + "/"+ path)}

func setRootPath(path string) {
	dataRootPath = path
}

func NewObject(name string) *Object {
	return &Object{
		afero.NewBasePathFs(afero.NewOsFs(),objectPath(name)),
	}
}

func (obj *Object) UploadAll(filename string, data []byte) error {
	f,err := obj.phy.OpenFile(filename, os.O_TRUNC | os.O_WRONLY ,os.FileMode(0644))
	if err != nil {return err}
	defer f.Close()

	_,err = f.Write(data)
	return err
}

func (obj *Object) Delete(filename string) error  {
	return obj.phy.Remove(filename)
}

func (obj *Object) Create(filename string) error {
	f,err:= obj.phy.Create(filename)
	if err != nil {return err}
	return f.Close()
}
func (obj *Object) GetAll(filename string) ([]byte,error) {
	f,err := obj.phy.Open(filename)
	if err != nil {return nil,err}
	defer f.Close()

	var info os.FileInfo
	info,err = f.Stat()
	if err != nil {return nil,err}

	var data = make([]byte,info.Size())
	_,err = f.Read(data)
	if err != nil {return nil,err}

	return data,nil
}

func (obj *Object)



