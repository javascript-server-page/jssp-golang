package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

func GenerateObjReq(jse *JsEngine, r *http.Request) *otto.Object {
	r.ParseForm()
	obj := jse.CreateObjectValue().Object()
	obj.Set("header", build_header(jse, r.Header))
	obj.Set("host", r.Host)
	obj.Set("method", r.Method)
	obj.Set("path", r.URL.Path)
	obj.Set("proto", r.Proto)
	obj.Set("remoteAddr", r.RemoteAddr)
	obj.Set("parm", func(call otto.FunctionCall) otto.Value {
		key := call.Argument(0)
		if key.IsUndefined() {
			return key
		}
		val := r.Form.Get(key.String())
		return *jse.CreateAny(val)
	})
	return obj
}
