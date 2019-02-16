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
	flag.StringVar(&(paras.Dir), "d", ".", "jssp folder")
	flag.StringVar(&(paras.Port), "p", "2019", "listening port")
	flag.Parse()
	if help {
		printUsage()
		flag.PrintDefaults()
		os.Exit(1)
	} else if version {
		printVersion()
		os.Exit(1)
	} else {
		paras.organize()
	}
}

// Organize command line parameters
func (paras *Parameter) organize() {
	_, err := strconv.Atoi(paras.Port)
	if err != nil {
		println("Port " + paras.Port + ":illegal")
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
