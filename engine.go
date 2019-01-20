package main

import (
	"github.com/robertkrimen/otto"
)

func InitOtto() *otto.Otto {
	vm := otto.New()
	return vm
}
