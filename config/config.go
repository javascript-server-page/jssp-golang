package config

import (
	"flag"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const Version = "0.9"
const ServerName = "jssp-golang-" + Version

var (
	Server  *server
	Log     *log
	Babel   *babel
	Vmpool  *vmpool
	Astpool *astpool
	Db      *db
)

func init() {
	var (
		help, version, Version, test, Test bool
		filename                           string
	)
	flag.BoolVar(&help, "h", false, "this help")
	flag.BoolVar(&version, "v", false, "show version and exit")
	flag.BoolVar(&Version, "V", false, "show version and configure options then exit")
	flag.BoolVar(&test, "t", false, "test configuration and exit")
	flag.BoolVar(&Test, "T", false, "test configuration, dump it and exit")
	flag.StringVar(&filename, "c", "config.yml", "set configuration file")

	flag.Parse()
	if help {
		printUsage()
		flag.PrintDefaults()
		os.Exit(1)
	} else if version {
		printVersion()
		os.Exit(1)
	}

	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	var Config = new(config)
	err = yaml.Unmarshal(configFile, &Config)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	Server = Config.Server
	Log = Config.Log
	Babel = Config.Babel
	Vmpool = Config.VmPool
	Astpool = Config.Astpool
	Db = Config.Db
}

func printUsage() {
	println(`Usage: jssp [-p port] [-d dir] [-l log]
Example: jssp -p 2019 -d . -l req.log

Options:`)
}

func printVersion() {
	println("jssp version:", Version)
}
