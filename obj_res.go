package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

func GenerateObjRes(jse *JsEngine, w http.ResponseWriter) *otto.Value {
	val, obj := jse.CreateObject()
	jse.Set("echo", func(call otto.FunctionCall) otto.Value {
		for _, e := range call.ArgumentList {
			w.Write([]byte(e.String()))
		}
		return otto.Value{}
	})
	obj.Set("print", func(call otto.FunctionCall) otto.Value {
		if val := call.Argument(0); !val.IsUndefined() {
			w.Write([]byte(val.String()))
		}
		w.Write([]byte("\n"))
		return otto.Value{}
	})
	return val
}
