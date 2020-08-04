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
	f.Write(activeBlockSize(b,f));f_b := f.Sum(nil)
	s.Write(activeBlockSize(f_b,s))
	return s.Sum(nil)
}

func activeBlockSize(b []byte,h hash.Hash) []byte {
	max := h.BlockSize()
	if len(b) > max {
		return b[:max-1]
	}else if len(b) < max {
		empty := make([]byte,max-len(b))
		return append(b,empty...)
	}
	return b
}