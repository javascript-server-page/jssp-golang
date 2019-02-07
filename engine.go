package main

import (
	"github.com/robertkrimen/otto"
)

var ancestor *otto.Otto = otto.New()

func GetJsEngine() *otto.Otto {
	return ancestor.Copy()
}
