package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

func GenerateObjRes(jse *JsEngine, w http.ResponseWriter) *otto.Object {
	obj, _ := jse.Object("{}")
	return obj
}
