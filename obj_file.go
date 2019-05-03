package main

import (
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func GenerateObjFile(jse *JsEngine, jspath string) *otto.Object {
	dir, _ := path.Split(jspath)
	obj := jse.CreateObject()
	obj.Set("open", func(call otto.FunctionCall) otto.Value {
		return *def_openfile(jse, &call, dir, os.O_RDWR)
	})
	obj.Set("opena", func(call otto.FunctionCall) otto.Value {
		return *def_openfile(jse, &call, dir, os.O_RDWR|os.O_APPEND)
	})
	obj.Set("create", func(call otto.FunctionCall) otto.Value {
		return *def_openfile(jse, &call, dir, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
	})
	obj.Set("mkdir", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc_error(jse, call, func(s string) error { return os.Mkdir(s, 0666) })
	})
	obj.Set("mkdirall", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc_error(jse, call, func(s string) error { return os.MkdirAll(s, 0666) })
	})
	obj.Set("remove", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc_error(jse, call, os.Remove)
	})
	obj.Set("removeall", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc_error(jse, call, os.RemoveAll)
	})
	obj.Set("readdir", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc_value(call, func(dirname string) *otto.Value {
			fis, err := ioutil.ReadDir(dirname)
			if err != nil {
				return jse.CreateError(err)
			}
			val := jse.CreateArray()
			for _, fi := range fis {
				val.Object().Call("push", fi.Name())
			}
			return val
		})
	})
	return obj
}

// get jssp.file object
func def_openfile(jse *JsEngine, call *otto.FunctionCall, dir string, flag int) *otto.Value {
	p := call.Argument(0)
	if p.IsUndefined() {
		return &p
	}
	name := filepath.Join(dir, p.String())
	return def_openfilebyname(jse, name, flag)
}

// get jssp.file object
func def_openfilebyname(jse *JsEngine, name string, flag int) *otto.Value {
	f, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		return jse.CreateError(err)
	}
	val := jse.CreateObjectValue()
	obj := val.Object()
	obj.Set("write", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc_error(jse, call, func(s string) error {
			_, err := f.WriteString(s)
			return err
		})
	})
	obj.Set("read", func(call otto.FunctionCall) otto.Value {
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return *jse.CreateError(err)
		}
		return *jse.CreateAny(string(data))
	})
	obj.Set("info", func(call otto.FunctionCall) otto.Value {
		return *build_fileinfo(jse, f)
	})
	obj.Set("parent", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateAny(path.Base(name))
	})
	obj.Set("close", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateError(f.Close())
	})
	obj.Set("move", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc_error(jse, call, func(newName string) error {
			err := os.Rename(name, newName)
			if err == nil {
				name = newName
			}
			return err
		})
	})
	return val
}

// execute (func(string) error) by js calling the parameter
func def_invokefunc_error(jse *JsEngine, call otto.FunctionCall, fun func(string) error) *otto.Value {
	p := call.Argument(0)
	if p.IsUndefined() {
		return &p
	}
	return jse.CreateError(fun(p.String()))
}

// execute (func(string) *otto.Value) by js calling the parameter
func def_invokefunc_value(call otto.FunctionCall, fun func(string) *otto.Value) *otto.Value {
	p := call.Argument(0)
	if p.IsUndefined() {
		return &p
	}
	return fun(p.String())
}

// build jssp.fileinfo by os.FileInfo
func build_fileinfo(jse *JsEngine, f *os.File) *otto.Value {
	fi, err := f.Stat()
	if err != nil {
		return jse.CreateError(err)
	}
	val := jse.CreateObjectValue()
	obj := val.Object()
	obj.Set("name", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateAny(fi.Name())
	})
	obj.Set("size", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateAny(fi.Size())
	})
	obj.Set("mode", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateAny(fi.Mode().String())
	})
	obj.Set("time", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateAny(fi.ModTime())
	})
	obj.Set("isdir", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateAny(fi.IsDir())
	})
	return val
}
