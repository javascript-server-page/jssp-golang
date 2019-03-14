package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjJssp(jse *JsEngine) *otto.Object {
	obj := jse.CreateObject()
	return obj
}
