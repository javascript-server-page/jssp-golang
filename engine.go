package main

import (
	"container/list"
	"github.com/robertkrimen/otto"
)

var ancestor *otto.Otto = otto.New()

var cache *list.List = list.New()

func NewJsEngine() *otto.Otto {
	return ancestor.Copy()
}


func GetJsEngine() *otto.Otto {
	if cache.Len() == 0 {

	}
	return NewJsEngine()
}
