package file

import (
	"io/ioutil"
	"jssp/config"
	"jssp/engine"
	"os"
	"path"
	"path/filepath"

	"github.com/dop251/goja"
	"github.com/valyala/fasthttp"
)

func init() {
	engine.ContextModules["file"] = GenerateObjFile
}

func GenerateObjFile(js *engine.VM, ctx *fasthttp.RequestCtx) *goja.Object {
	dir, _ := path.Split(path.Join(config.Server.Dir, string(ctx.RequestURI())))
	obj := js.CreateObject(nil)
	obj.Set("open", func(filename string) goja.Value {
		return build_file(js, filepath.Join(dir, filename), os.O_RDWR)
	})
	obj.Set("opena", func(filename string) goja.Value {
		return build_file(js, filepath.Join(dir, filename), os.O_RDWR|os.O_APPEND)
	})
	obj.Set("create", func(filename string) goja.Value {
		return build_file(js, filepath.Join(dir, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC)
	})
	obj.Set("mkdir", func(dir string) goja.Value {
		return js.NewGoError(os.Mkdir(dir, 0666))
	})
	obj.Set("mkdirall", func(dir string) goja.Value {
		return js.NewGoError(os.MkdirAll(dir, 0666))
	})
	obj.Set("remove", func(file string) goja.Value {
		return js.NewGoError(os.Remove(file))
	})
	obj.Set("removeall", func(file string) goja.Value {
		return js.NewGoError(os.RemoveAll(file))
	})
	obj.Set("readdir", func(dirname string) goja.Value {
		fis, err := ioutil.ReadDir(path.Join(dir, dirname))
		if err != nil {
			panic(js.NewGoError(err))
		}
		fs := make([]string, len(fis))
		for i := range fis {
			fs[i] = fis[i].Name()
		}
		return js.ToValue(fs)
	})
	return obj
}

// get jssp.file object
func build_file(js *engine.VM, name string, flag int) *goja.Object {
	f, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		panic(js.NewGoError(err))
	}
	obj := js.CreateObject(nil)
	obj.Set("write", func(str string) int {
		n, err := f.WriteString(str)
		if err != nil {
			panic(js.NewGoError(err))
		}
		return n
	})
	obj.Set("read", func() goja.Value {
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(js.NewGoError(err))
		}
		return js.ToValue(string(data))
	})
	obj.Set("info", func() goja.Value {
		return build_fileinfo(js, f)
	})
	obj.Set("parent", func() string {
		return path.Base(name)
	})
	obj.Set("close", func() goja.Value {
		return js.NewGoError(f.Close())
	})
	obj.Set("move", func(newName string) {
		err := os.Rename(name, newName)
		if err != nil {
			panic(js.NewGoError(err))
		}
		name = newName
	})
	return obj
}

// build jssp.fileinfo by os.FileInfo
func build_fileinfo(js *engine.VM, f *os.File) *goja.Object {
	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	obj := js.CreateObject(nil)
	obj.Set("name", fi.Name)
	obj.Set("size", fi.Size)
	obj.Set("mode", fi.Mode().String)
	obj.Set("time", fi.ModTime)
	obj.Set("isdir", fi.IsDir)
	return obj
}
