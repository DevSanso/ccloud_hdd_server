package db

import (
	"database/sql"
)

type FHeader struct {
	name string
	baseID int
	subDirPath string

	tokenSize int
	size int64
	date string
	ext string
}

func NewFHeader() *FHeader {

}

func (fh *FHeader)getBase()
