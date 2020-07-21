package file

import (
	"os"
	"github.com/spf13/afero"
)

func GetFileList(fs afero.Fs,path string) ([]os.FileInfo,error) {
	f ,err := fs.Open(path)
	if err != nil {return nil,err}
	defer f.Close()
	return  f.Readdir(-1)
}