package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjJsdo(js *JavaScript) *otto.Object {
	obj := js.CreateObjectValue().Object()
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
