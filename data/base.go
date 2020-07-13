package data


import (
	"github.com/spf13/afero"
	"os"
)

type RootFs struct {
	rd afero.Fs
}


func NewRootFs(path string) RootFs {

}

func (rd *RootFs)ExistDir(path string)bool {
	
}




func (rd *RootFs)GetDirFiles(path string) ([]os.FileInfo,error) {
	
}

func (rd *RootFs)ExistFile(dirpath,filename string) bool {
	
}





