package file

import (
	"ccloud_hdd_server/pkg/data"
	"ccloud_hdd_server/pkg/db_sql"
	"os"

	"github.com/spf13/afero"
)

func ReadFile(fs afero.Fs, name string, key []byte, iv []byte, header *db_sql.Header) (*data.Object, error) {
	enc_name, encode_err := EncodeFilePath(key, iv, name)
	if encode_err != nil {
		return nil, encode_err
	}
	f, err := fs.Open(enc_name)
	if err != nil {
		return nil, err
	}
	return data.NewObject(f, key, iv, header.TokenSize(), header.Size())
}
func WriteFile(fs afero.Fs, name string, key []byte, iv []byte, header *db_sql.Header) (*data.Object, error) {
	enc_name, encode_err := EncodeFilePath(key, iv, name)
	if encode_err != nil {
		return nil, encode_err
	}
	f, err := fs.OpenFile(enc_name, os.O_WRONLY|os.O_CREATE, os.FileMode(0))
	if err != nil {
		return nil, err
	}
	return data.NewObject(f, key, iv, header.TokenSize(), header.Size())
}

func CreateEmptyFile(fs afero.Fs, name string, key []byte, iv []byte, tokenSize int) (*data.Object, error) {
	enc_name, encode_err := EncodeFilePath(key, iv, name)
	if encode_err != nil {
		return nil, encode_err
	}
	f, err := fs.Create(enc_name)
	if err != nil {
		return nil, err
	}
	return data.NewObject(f, key, iv, tokenSize, 0)
}
