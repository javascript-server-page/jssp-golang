package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjHttp(jse *JsEngine) *otto.Object {
	obj, _ := jse.Object("({})")
	return obj
}
