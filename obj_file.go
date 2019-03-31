package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjFile(jse *JsEngine) *otto.Object {
	obj := jse.CreateObject()
	obj.Set("open", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
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
