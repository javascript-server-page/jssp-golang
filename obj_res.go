package main

import (
	"github.com/robertkrimen/otto"
	"mime"
	"net/http"
	"strings"
)

func GenerateObjRes(js *JavaScript, w http.ResponseWriter) *otto.Object {
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
		json, err := js.Get("JSON")
		if err != nil {
			return *js.CreateError(err)
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
	obj.Set("type", func(call otto.FunctionCall) otto.Value {
		fval := call.Argument(0)
		if fval.IsUndefined() {
			return fval
		}
		ct := mime.TypeByExtension("." + fval.String())
		w.Header().Set("Content-Type", ct)
		return otto.Value{}
	})
	obj.Set("include", func(call otto.FunctionCall) otto.Value {
		file, err := js.Get("file")
		if err != nil {
			return *js.CreateError(err)
		}
		fval := call.Argument(0)
		if fval.IsUndefined() {
			return fval
		}
		fname := fval.String()
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
			return *def_runsrc(js, src.String())
		} else if strings.HasSuffix(fname, ".jssp") {
			src := jssp_jsjs([]byte(src.String()))
			return *def_runsrc(js, src)
		} else {
			w.Write([]byte(src.String()))
			return otto.UndefinedValue()
		}
	})
	return obj
}

// run js src code
func def_runsrc(js *JavaScript, src interface{}) *otto.Value {
	str, err := js.Eval(src)
	if err != nil {
		return js.CreateError(err)
	} else {
		return &str
	}
}
