package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

func GenerateObjReq(jse *JsEngine, r *http.Request) *otto.Object {
	obj := jse.CreateObject()
	obj.Set("method", r.Method)
	return obj
}
