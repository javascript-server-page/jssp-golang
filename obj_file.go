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
		return *def_invokefunc(jse, call, func(s string) error { return os.Mkdir(s, 0666) })
	})
	obj.Set("mkdirall", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc(jse, call, func(s string) error { return os.MkdirAll(s, 0666) })
	})
	obj.Set("remove", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc(jse, call, os.Remove)
	})
	obj.Set("removeall", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc(jse, call, os.RemoveAll)
	})
	obj.Set("rename", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
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
	f, err := os.OpenFile(name, flag, 0666)
	if err != nil {
		return jse.CreateError(err)
	}
	val := jse.CreateObjectValue()
	obj := val.Object()
	obj.Set("write", func(call otto.FunctionCall) otto.Value {
		return *def_invokefunc(jse, call, func(s string) error {
			_, err := f.WriteString(s)
			return err
		})
	})
	obj.Set("read", func(call otto.FunctionCall) otto.Value {
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return *jse.CreateError(err)
		}
		return *jse.CreateString(string(data))
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
	obj.Set("close", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateError(f.Close())
	})
	return val
}

// execute (func(string) error) by js calling the parameter
func def_invokefunc(jse *JsEngine, call otto.FunctionCall, fun func(string) error) *otto.Value {
	p := call.Argument(0)
	if p.IsUndefined() {
		return &p
	}
	return jse.CreateError(fun(p.String()))
}
