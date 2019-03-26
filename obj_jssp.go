package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjJssp(jse *JsEngine) *otto.Object {
	obj := jse.CreateObject()
	obj.Set("exec", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return obj
}
