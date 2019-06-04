package main

import (
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func GenerateObjFile(js *JavaScript, jspath string) *otto.Object {
	dir, _ := path.Split(jspath)
	obj := js.CreateObjectValue().Object()
	obj.Set("open", func(filename string) otto.Value {
		return *def_openfile(js, filepath.Join(dir, filename), os.O_RDWR)
	})
	obj.Set("opena", func(filename string) otto.Value {
		return *def_openfile(js, filepath.Join(dir, filename), os.O_RDWR|os.O_APPEND)
	})
	obj.Set("create", func(filename string) otto.Value {
		return *def_openfile(js, filepath.Join(dir, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC)
	})
	obj.Set("mkdir", func(dir string) otto.Value {
		return *js.CreateError(os.Mkdir(dir, 0666))
	})
	obj.Set("mkdirall", func(dir string) otto.Value {
		return *js.CreateError(os.MkdirAll(dir, 0666))
	})
	obj.Set("remove", func(file string) otto.Value {
		return *js.CreateError(os.Remove(file))
	})
	obj.Set("removeall", func(file string) otto.Value {
		return *js.CreateError(os.RemoveAll(file))
	})
	obj.Set("readdir", func(dirname string) otto.Value {
		fis, err := ioutil.ReadDir(dirname)
		if err != nil {
			return *js.CreateError(err)
		}
		fs := make([]string, len(fis))
		for i := range fis {
			fs[i] = fis[i].Name()
		}
		return *js.CreateAny(fs)
	})
	return obj
}

// get jssp.file object
func def_openfile(js *JavaScript, name string, flag int) *otto.Value {
	f, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		return js.CreateError(err)
	}
	val := js.CreateObjectValue()
	obj := val.Object()
	obj.Set("write", func(str string) otto.Value {
		_, err := f.WriteString(str)
		return *js.CreateError(err)
	})
	obj.Set("read", func() otto.Value {
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return *js.CreateError(err)
		}
		return *js.CreateAny(string(data))
	})
	obj.Set("info", func() otto.Value {
		return *build_fileinfo(js, f)
	})
	obj.Set("parent", func() string {
		return path.Base(name)
	})
	obj.Set("close", func() otto.Value {
		return *js.CreateError(f.Close())
	})
	obj.Set("move", func(newName string) otto.Value {
		err := os.Rename(name, newName)
		if err == nil {
			name = newName
		}
		return *js.CreateError(err)
	})
	return val
}

// build jssp.fileinfo by os.FileInfo
func build_fileinfo(js *JavaScript, f *os.File) *otto.Value {
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	val := js.CreateObjectValue()
	obj := val.Object()
	obj.Set("name", fi.Name)
	obj.Set("size", fi.Size)
	obj.Set("mode", fi.Mode().String)
	obj.Set("time", fi.ModTime)
	obj.Set("isdir", fi.IsDir)
	return val
}
