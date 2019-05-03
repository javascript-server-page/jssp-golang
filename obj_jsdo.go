package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjJsdo(jse *JsEngine) *otto.Object {
	obj := jse.CreateObjectValue().Object()
	obj.Set("mysql", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("sqlserver", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("postgres", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return obj
}
