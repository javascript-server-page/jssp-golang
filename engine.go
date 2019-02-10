package main

import (
	"github.com/robertkrimen/otto"
)

var ancestor *otto.Otto = otto.New()

var cache []*otto.Otto = make([]*otto.Otto, 100)

func NewJsEngine() *otto.Otto {
	return ancestor.Copy()
}


func GetJsEngine() *otto.Otto {
	if len(cache) == 0 {

	}
	return NewJsEngine()
}
