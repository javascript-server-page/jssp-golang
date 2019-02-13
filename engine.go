package main

import (
	"container/list"
	"github.com/robertkrimen/otto"
)

const cache_max = 500

var ancestor *otto.Otto = otto.New()

var cache *list.List = list.New()

var isGenerate = make(chan bool)

func init() {
	go generate()
}

func generate() {
	for {
		<-isGenerate
		for cache.Len() < cache_max {
			cache.PushBack(NewJsEngine())
		}
	}
}

func NewJsEngine() *otto.Otto {
	return ancestor.Copy()
}

func GetJsEngine() *otto.Otto {
	if cache.Len() == 0 {
		isGenerate <- true
		return NewJsEngine()
	}
	return cache.Remove(cache.Front()).(*otto.Otto)
}
