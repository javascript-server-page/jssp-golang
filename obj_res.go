package main

import (
	"github.com/robertkrimen/otto"
	"mime"
	"net/http"
	"strings"
)

func GenerateObjRes(js *JavaScript, w http.ResponseWriter) *otto.Object {
	json, _ := js.Get("JSON")
	file, _ := js.Get("file")
	js.Set("echo", func(call otto.FunctionCall) otto.Value {
		for _, e := range call.ArgumentList {
			str := e.String()
			if len(str) > 0 {
				w.Write([]byte(str))
			}
		}
		return otto.Value{}
	})
	js.Set("print", func(call otto.FunctionCall) otto.Value {
		n := len(call.ArgumentList)
		if n == 0 {
			return otto.UndefinedValue()
		}
		var value interface{}
		if n == 1 {
			value = call.Argument(0).Object()
		} else {
			arr := js.CreateArray().Object()
			for _, e := range call.ArgumentList {
				arr.Call("push", e)
			}
			value = arr
		}
		str, err := json.Object().Call("stringify", value)
		if err != nil {
			return *js.CreateError(err)
		}
		w.Write([]byte(str.String()))
		return otto.Value{}
	})
	obj := js.CreateObjectValue().Object()
	obj.Set("header", build_editableheader(js, w.Header()))
	obj.Set("type", func(name string) {
		ct := mime.TypeByExtension("." + name)
		w.Header().Set("Content-Type", ct)
	})
	obj.Set("include", func(fname string) otto.Value {
		f, err := file.Object().Call("open", fname)
		if err != nil {
			return *js.CreateError(err)
		}
		if js.isError(&f) {
			return f
		}
		defer f.Object().Call("close")
		src, err := f.Object().Call("read")
		if err != nil {
			return *js.CreateError(err)
		}
		if js.isError(&src) {
			return src
		}
		if strings.HasSuffix(fname, "js") {
			val, err := js.Eval(src.String())
			if err != nil {
				return *js.CreateError(err)
			}
			return val
		} else if strings.HasSuffix(fname, ".jssp") {
			val, err := js.Eval(src.String())
			if err != nil {
				return *js.CreateError(err)
			}
			return val
		} else {
			w.Write([]byte(src.String()))
			return otto.UndefinedValue()
		}
	})
	return obj
}
