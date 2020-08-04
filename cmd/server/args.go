package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
)


type Config struct {
	Http struct {
		Host string
		Port int
	}
	Ws struct { 
		Host string
		Port int
	}

	CertFile string
	KeyFile string

	DbConfig struct {
		Host string
		Port int
		User string
		Passwd string
		DbName string
	}
}


func GetConfig() Config {
	var c Config
	s := flag.String("cfg_file","./cfg.json","server config file")
	flag.Parse()
	data,err := ioutil.ReadFile(*s)
	if err != nil {panic(err)}

	if err = json.Unmarshal(data,&c);err != nil {
		panic(err)
	}

	return c
}