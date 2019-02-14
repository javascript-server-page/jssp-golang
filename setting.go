package main

import (
	"flag"
	"os"
	"strconv"
)

type Parameter struct {
	Dir  string
	Port string
}

func (paras *Parameter) Init() {
	var (
		port          int
		dir           string
		help, version bool
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	flag.StringVar(&dir, "d", ".", "jssp folder")
	flag.IntVar(&port, "p", 2019, "listening port")
	flag.Parse()
	paras.Dir = dir
	paras.Port = strconv.Itoa(port)
	if help {
		printUsage()
		flag.PrintDefaults()
		os.Exit(1)
	} else if version {
		printVersion()
		os.Exit(1)
	}
}

func printUsage() {
	println(`Usage: jssp [-p port] [-d dir]
Example: jssp -p 2019 -d .

Options:`)
}

func printVersion() {
	println("jssp version:", Version)
}
