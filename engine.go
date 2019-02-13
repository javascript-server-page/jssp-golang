package main

import (
	"container/list"
	"github.com/robertkrimen/otto"
	"sync"
)

const cache_max = 500

var mutex sync.Mutex

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
	mutex.Lock()
	defer mutex.Unlock()
	if cache.Len() == 0 {
		isGenerate <- true
		return NewJsEngine()
	}
	return cache.Remove(cache.Front()).(*otto.Otto)
}
