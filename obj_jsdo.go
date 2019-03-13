package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjJsdo(jse *JsEngine) *otto.Object {
	obj, _ := jse.Object("({})")
	return obj
}
