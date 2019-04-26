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
		file, err := jse.Get("file")
		if err != nil {
			return *jse.CreateError(err)
		}
		f, err := file.Object().Call("open", call.Argument(0))
		if err != nil {
			return *jse.CreateError(err)
		}
		defer f.Object().Call("close")
		src, err := f.Object().Call("read")
		if err != nil {
			return *jse.CreateError(err)
		}
		str, err := jse.Run(src.String())
		if err != nil {
			return *jse.CreateError(err)
		} else {
			return *jse.CreateString(str.String())
		}
	})
	return obj
}
