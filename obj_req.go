package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

func GenerateObjReq(jse *JsEngine, r *http.Request) *otto.Object {
	obj := jse.CreateObject()
	return obj
}
