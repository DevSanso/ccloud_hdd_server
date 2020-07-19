package file

import (
	"os"
	"github.com/spf13/afero"
)

func GetDirList(fs afero.Fs,path string) ([]os.FileInfo,error) {
	f ,err := fs.Open(path)
	if err != nil {return nil,err}
	defer f.Close()
	var finfo []os.FileInfo
	finfo ,err = f.Readdir(-1)
	if err != nil {return nil,err}

	var info_arr []os.FileInfo
	for _,info := range finfo {
		if !info.IsDir() {continue}
		info_arr = append(info_arr,info)
	}
	return info_arr,err
}