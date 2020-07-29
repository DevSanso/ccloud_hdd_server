package server


import (
	"net/http"
	"errors"
)

var (
	NoMatchUrlLenErr = errors.New("no match url length")
	NotExistUrlInWsErr = errors.New("no exist access url")
	InternalServerErr = errors.New("internal error")
	NotMatchIpErr = errors.New("not match ip error")
)



func writeErrToRes(w http.ResponseWriter,err error) {
	w.Header().Set("content-type","text/plain")
	w.Write([]byte(err.Error()))
}