package main

import (
	"net/http"
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