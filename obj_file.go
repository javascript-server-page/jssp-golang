package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjFile(jse *JsEngine) *otto.Object {
	obj := jse.CreateObject()
	obj.Set("open", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return obj
}
