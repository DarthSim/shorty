package main

import (
	"flag"
	"fmt"

	"code.google.com/p/gcfg"
)

type Config struct {
	Url struct {
		Domain string
	}
	Server struct {
		Address string
		Port    string
	}
	Database struct {
		Host               string
		Database           string
		User               string
		Password           string
		MaxOpenConnections int
		MaxIdleConnections int
		InitSchema         bool `gcfg:"-"`
	}
	Log struct {
		Path string
	}
}

var config Config

func prepareConfig() {
	configfile := flag.String(
		"config",
		"../shorty.conf",
		"path to configuration file",
	)

	flag.BoolVar(
		&config.Database.InitSchema,
		"init-db-schema",
		false,
		"need to initialize DB schema",
	)

	flag.Parse()

	err := gcfg.ReadFileInto(&config, absPathToFile(*configfile))
	if err != nil {
		fmt.Printf("Error opening config file: %v", err)
		halt()
	}
}
