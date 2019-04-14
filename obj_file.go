package main

import (
	"github.com/robertkrimen/otto"
	"os"
	"path"
)

func GenerateObjFile(jse *JsEngine, jspath string) *otto.Object {
	dir, _ := path.Split(jspath)
	obj := jse.CreateObject()
	obj.Set("open", func(call otto.FunctionCall) otto.Value {
		p := call.Argument(0)
		f, err := def_openfile(&p, dir, os.O_RDWR)
		return *build_file(jse, f, err)
	})
	obj.Set("opena", func(call otto.FunctionCall) otto.Value {
		p := call.Argument(0)
		f, err := def_openfile(&p, dir, os.O_RDWR|os.O_APPEND)
		return *build_file(jse, f, err)
	})
	obj.Set("create", func(call otto.FunctionCall) otto.Value {
		p := call.Argument(0)
		f, err := def_openfile(&p, dir, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
		return *build_file(jse, f, err)
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

// execute (func(string) error) by js calling the parameter
func def_invokefunc(jse *JsEngine, call otto.FunctionCall, fun func(string) error) *otto.Value {
	p := call.Argument(0)
	if p.IsUndefined() {
		return &p
	}
	return jse.CreateError(fun(p.String()))
}

func def_openfile(v *otto.Value, dir string, flag int) (*os.File, error) {
	if v.IsUndefined() {
		return nil, nil
	} else {
		f, err := os.OpenFile(v.String(), flag, 0666)
		return f, err
	}
}

func build_file(jse *JsEngine, f *os.File, err error) *otto.Value {
	if f == nil {
		if err != nil {
			return &otto.Value{}
		} else {
			return &otto.Value{}
		}
	}
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
