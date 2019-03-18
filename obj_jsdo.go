package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjJsdo(jse *JsEngine) *otto.Value {
	val, obj := jse.CreateObject()
	obj.Set("mysql", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return val
}
