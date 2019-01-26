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
		help, version bool
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	paras.Dir = *flag.String("d", ".", "jssp folder")
	paras.Port = strconv.Itoa(*flag.Int("p", 2019, "listening port"))
	flag.Parse()
	flag.Usage = usage
	if help {
		printVersion()
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	} else if version {
		printVersion()
		os.Exit(1)
	}
}

func usage() {
	println(`Usage: jssp [-p port] [-d dir]
Example: jssp -p 2019 -d .

Options:`)
}

func printVersion() {
	println("jssp version:", Version)
}
