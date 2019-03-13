package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjJssp(jse *JsEngine) *otto.Object {
	obj, _ := jse.Object("({})")
	return obj
}
