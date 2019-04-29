package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func getFile(fs http.FileSystem, name string) http.File {
	f, err := fs.Open(name)
	if err != nil {
		return nil
	}

	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return nil
	}
	if stat.IsDir() {
		f.Close()
		return nil
	} else {
		return f
	}
}

func readFile(f http.File) ([]byte, error) {
	defer f.Close()
	return ioutil.ReadAll(f)
}

func jssp_jsjs(data []byte) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString(`echo("`)
	isJsjs := false
	for i, n := 0, len(data); i < n; i++ {
		c := data[i]
		if isJsjs {
			if c == '%' && data[i+1] == '>' {
				buf.WriteString(`;echo("`)
				i++
				isJsjs = false
			} else {
				buf.WriteByte(c)
			}
		} else {
			if c == '<' && data[i+1] == '%' {
				buf.WriteString(`");`)
				i++
				isJsjs = true
			} else {
				switch c {
				case '\n':
					buf.WriteString(`\n");`)
					buf.WriteByte(c)
					buf.WriteString(`echo("`)
				case '\r':
					continue
				case '"':
					buf.WriteString(`\"`)
				default:
					buf.WriteByte(c)
				}
			}
		}
	}
	buf.WriteString(`");`)
	return buf.Bytes()
}
