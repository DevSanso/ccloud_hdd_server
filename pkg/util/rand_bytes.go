package util

import (
	"bytes"
	"math/rand"
)




func MakeBytes(n int) []byte {
	var i = 0;var err error
	var buf = bytes.NewBuffer([]byte(""))
	for i < n{
		_,err = buf.WriteRune(rand.Int31())
		i++
	}
	if err != nil {panic(err)}
	var res []byte;copy(res,buf.Bytes())
	return res
}