package main

import (
	"github.com/robertkrimen/otto"
	"os"
	"path"
	"path/filepath"
)

func GenerateObjFile(jse *JsEngine, dir string) *otto.Object {
	obj := jse.CreateObject()
	obj.Set("open", func(call otto.FunctionCall) otto.Value {
		fn := call.Argument(0).String()
		f, _ := os.OpenFile(fn, os.O_RDWR, 0666)
		return *build_file(jse, f)
	})
	obj.Set("create", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("remove", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("mkdir", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("rmdir", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return obj
}

func build_file(jse *JsEngine, f *os.File) *otto.Value {
	val := jse.CreateObjectValue()
	obj := val.Object()
	obj.Set("write", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("read", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("isdir", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("parent", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("children", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return val
}

func def_getfullname(dir, name string) string {
	if dir == "" {
		dir = "."
	}
	return filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
}
