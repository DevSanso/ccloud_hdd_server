package file

import (
	"ccloud_hdd_server/pkg/cryp"
	"path/filepath"
	"strings"
)

func EncodeFilePath(key, iv []byte, path string) (string, error) {
	encoder, err := cryp.NewEncoder(key, iv)
	if err != nil {
		return "", err
	}

	var builder = strings.Builder{}
	var buf []byte
	var nodes = strings.Split(path, "/")
	for i, node := range nodes {
		buf = make([]byte, len(node))
		err = encoder.Encrypt([]byte(node), buf)
		if err != nil {
			return "", err
		}

		_, err = builder.Write(buf)
		if i != len(nodes)-1 {
			_, err = builder.WriteRune('/')
		}

		if err != nil {
			return "", err
		}
	}

	return filepath.FromSlash(builder.String()), nil
}

func DecodeFilePath(key, iv []byte, subDir, name string, path string) (string, error) {
	decoder, err := cryp.NewDecoder(key, iv)
	if err != nil {
		return "", err
	}

	var builder = strings.Builder{}
	var buf []byte
	var nodes = strings.Split(path, "/")
	for i, node := range nodes {
		buf = make([]byte, len(node))
		err = decoder.Decrypt([]byte(node), buf)
		if err != nil {
			return "", err
		}

		_, err = builder.Write(buf)
		if i != len(nodes)-1 {
			_, err = builder.WriteRune('/')
		}

		if err != nil {
			return "", err
		}
	}

	return filepath.FromSlash(builder.String()), nil
}
