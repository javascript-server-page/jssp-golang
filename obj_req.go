package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

func GenerateObjReq(jse *JsEngine, r *http.Request) *otto.Object {
	obj := jse.CreateObjectValue().Object()
	obj.Set("header", build_header(jse, r.Header))
	obj.Set("host", r.Host)
	obj.Set("method", r.Method)
	obj.Set("path", r.URL.Path)
	obj.Set("proto", r.Proto)
	obj.Set("remoteAddr", r.RemoteAddr)
	return obj
}
