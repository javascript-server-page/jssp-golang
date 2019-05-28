package main

import (
	"flag"
	"os"
	"strconv"
)

type Setting struct {
	Dir  string
	Log  string
	Port string
}

func (set *Setting) Init() {
	var (
		help, version bool
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	flag.StringVar(&(set.Dir), "d", ".", "jssp folder")
	flag.StringVar(&(set.Log), "l", "req.log", "log file")
	flag.StringVar(&(set.Port), "p", "2019", "listening port")
	flag.Parse()
	if help {
		printUsage()
		flag.PrintDefaults()
		os.Exit(1)
	} else if version {
		printVersion()
		os.Exit(1)
	} else {
		set.organize()
	}
}

// Organize command line parameters
func (set *Setting) organize() {
	_, err := strconv.Atoi(set.Port)
	if err != nil {
		println("Port " + set.Port + ":illegal")
		os.Exit(1)
	}
}

func printUsage() {
	println(`Usage: jssp [-p port] [-d dir] [-l log]
Example: jssp -p 2019 -d . -l req.log

Options:`)
}

func printVersion() {
	println("jssp version:", Version)
}
