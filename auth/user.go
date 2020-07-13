package user

import (
	"context"

	"ccloud_hdd_server/data"
	"ccloud_hdd_server/db"
)


type UserMap struct {
	m map[string] struct{
		basePath string
		ctx context.Context
	}
}

var um = func() UserMap {
	return {
		make(map[string]struct{string,context.Context})
	}
}()


func Login(hash []byte) error {

}

func Logout(hash []byte) error {

}

func GetRootFs() 