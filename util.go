package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

func getFile(fs http.FileSystem, name string) *http.File {
	f, err := fs.Open(name)
	if err != nil {
		return nil
	}
	stat, err := f.Stat()
	if err != nil {
		return nil
	}
	if stat.IsDir() {
		f.Close()
		return nil
	} else {
		return &f
	}
}

func readFile(f *http.File) ([]byte, error) {
	defer (*f).Close()
	return ioutil.ReadAll(*f)
}

var repl = []byte(`);\n $1 \n  echo(`)

var reg = regexp.MustCompile("<%([\\s\\S]+?)%>")

func jssp_jsjs(data []byte) []byte {
	return reg.ReplaceAll(data, repl)
}
