package data


import (
	"github.com/spf13/afero"
	"path/filepath"
)

func GetBasePath() string {
	const baseP = "./datas"
	p,_ :=filepath.Abs(filepath.FromSlash(baseP))
	return p
}

var base = afero.NewBasePathFs(afero.NewOsFs(),GetBasePath())


type rootData struct {fs afero.Fs}
var rd = rootData{base}
func (rd *rootData)getObject(key []byte,path string) (*Object,error) {
	
}


func GetObject(key []byte,path string) (*Object,error) {
	return rd.getObject(key,path)
}