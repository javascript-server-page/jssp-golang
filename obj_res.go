package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

func GenerateObjRes(jse *JsEngine, w http.ResponseWriter) *otto.Object {
	obj := jse.CreateObject()
	obj.Set("header", build_editableheader(jse, w.Header()))
	jse.Set("echo", func(call otto.FunctionCall) otto.Value {
		for _, e := range call.ArgumentList {
			w.Write([]byte(e.String()))
		}
		return otto.Value{}
	})
	jse.Set("include", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return obj
}
