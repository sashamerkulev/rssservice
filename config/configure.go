package config

import (
	"fmt"
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Connection struct {
		Mysql string
	}
	Server struct {
		Port string
	}
}

func GetConfig() (cfg Config) {
	err := gcfg.ReadFileInto(&cfg, "config.gcfg")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	return cfg
}
