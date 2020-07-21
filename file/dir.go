package file

import (
	"os"
	"time"

	"github.com/spf13/afero"
	"ccloud_hdd_server/cryp"
)


type decodeInfo struct {
	name string
	raw os.FileInfo
}


func(dei *decodeInfo)Name() string {return dei.name}
func(dei *decodeInfo)Size() int64 {return dei.raw.Size()}
func(dei *decodeInfo)Mode() os.FileMode {return dei.raw.Mode()}
func(dei *decodeInfo)ModTime() time.Time {return dei.raw.ModTime()}
func(dei *decodeInfo)IsDir() bool {return dei.raw.IsDir()}
func(dei *decodeInfo)Sys() interface{} {return dei.raw.Sys()}


func GetFileList(fs afero.Fs,key,iv []byte,path string) ([]os.FileInfo,error) {
	enc_path, enc_err := EncodeFilePath(key,iv,path)
	if enc_err != nil {return nil,enc_err}
	f ,err := fs.Open(enc_path)
	if err != nil {return nil,err}
	defer f.Close()

	var infos,res []os.FileInfo
	infos,err = f.Readdir(-1)
	if err != nil {return nil,err}
	
	decode,dec_err := cryp.NewDecoder(key,iv)
	if dec_err != nil {return nil,dec_err}

	var buf []byte
	for _,info := range infos {
		enc_name := []byte(info.Name())
		buf = make([]byte,len(enc_name))

		err = decode.Decrypt(enc_name,buf)
		if err != nil {return nil,err}
		
		res = append(res,&decodeInfo{
			string(buf),
			info,
		})
	}

	return res,nil
}