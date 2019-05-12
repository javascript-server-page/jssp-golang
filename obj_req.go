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
	obj.Set("cookie", build_cookie(jse, r))
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

// build req.cookie object
func build_cookie(jse *JsEngine, r *http.Request) *otto.Object {
	obj := jse.CreateObjectValue().Object()
	obj.Set("get", func(call otto.FunctionCall) otto.Value {
		val := call.Argument(0)
		if val.IsUndefined() {
			return val
		}
		c, err := r.Cookie(val.String())
		if err != nil {
			return otto.UndefinedValue()
		}
		return *jse.CreateAny(c.Value)
	})
	obj.Set("set", func(call otto.FunctionCall) otto.Value {
		key := call.Argument(0)
		if key.IsUndefined() || key.IsNull() {
			return key
		}
		keystr := key.String()
		val := call.Argument(1)
		c := &http.Cookie{Name: keystr, Value: val.String(),}
		if val.IsUndefined() || val.IsNull() {
			c.MaxAge = -1
		}
		r.AddCookie(c)
		return otto.UndefinedValue()
	})
	obj.Set("map", func(call otto.FunctionCall) otto.Value {
		val := jse.CreateObjectValue()
		obj := val.Object()
		for _, k := range r.Cookies() {
			obj.Set(k.Name, k.Value)
		}
		return *val
	})
	return obj
}
