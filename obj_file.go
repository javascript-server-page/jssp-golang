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
		p := call.Argument(0)
		f := def_openfile(&p, dir, os.O_RDWR)
		return *build_file(jse, f)
	})
	obj.Set("opena", func(call otto.FunctionCall) otto.Value {
		p := call.Argument(0)
		f := def_openfile(&p, dir, os.O_RDWR|os.O_APPEND)
		return *build_file(jse, f)
	})
	obj.Set("create", func(call otto.FunctionCall) otto.Value {
		p := call.Argument(0)
		f := def_openfile(&p, dir, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
		return *build_file(jse, f)
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

func def_openfile(v *otto.Value, dir string, flag int) *os.File {
	if v.IsUndefined() {
		return nil
	} else {
		f, _ := os.OpenFile(v.String(), flag, 0666)
		return f
	}
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
