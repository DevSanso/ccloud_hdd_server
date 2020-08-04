package use_hash

import (
	"hash"
	"crypto/md5"
	"crypto/sha256"
)

type _HashDecryption struct {}
func (*_HashDecryption)first() hash.Hash {
	return md5.New()
}
func (*_HashDecryption)second() hash.Hash {
	return sha256.New()
}

var _hd = _HashDecryption{}

func Sum(b []byte) []byte {
	f,s := _hd.first(),_hd.second()
	f.Write(b)
	s.Write(f.Sum(nil))
	return s.Sum(nil)
}