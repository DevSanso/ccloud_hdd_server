package cryp

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	NotMatchArraySizeError = errors.New("NotMatchArraySizeError")
)
func getBlock(key []byte) (cipher.Block,error) {
	if len(key) != 32 {return nil,NotMatchArraySizeError}
	return aes.NewCipher(key)
}

type Encoder struct {
	block cipher.Block
	iv []byte
}

func NewEncoder(key []byte,iv []byte) (*Encoder,error) {
	block,err := getBlock(key)
	if err != nil && len(iv) != aes.BlockSize {return nil,err}

	return &Encoder{
		block,
		iv,
	},nil
}

func (e *Encoder)Encrypt(src,dst []byte) error {
	if len(src) != len(dst) {return NotMatchArraySizeError}
	cipher.NewCFBEncrypter(e.block,e.iv).XORKeyStream(dst,src)
	return nil
}



type Decoder struct {
	block cipher.Block
	iv []byte
}
func NewDecoder(key []byte,iv []byte) (*Decoder,error) {
	block,err := getBlock(key)
	if err != nil && len(iv) != aes.BlockSize {return nil,err}

	return &Decoder{
		block,
		iv,
	},nil
}
func (e *Decoder)Decrypt(src,dst []byte) error {
	if len(src) != len(dst) {return NotMatchArraySizeError}
	cipher.NewCFBDecrypter(e.block,e.iv).XORKeyStream(dst,src)
	return nil
}

