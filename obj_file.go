package main

import (
	"github.com/robertkrimen/otto"
)

func GenerateObjFile(jse *JsEngine) *otto.Object {
	obj, _ := jse.Object("{}")
	return obj
}
