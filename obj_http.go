package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjHttp(jse *JsEngine) *otto.Object {
	obj := jse.CreateObject()
	return obj
}
