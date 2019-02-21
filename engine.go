package main

import (
	"container/list"
	"fmt"
	"github.com/robertkrimen/otto"
	"sync"
)

type JsEngine struct {
	*otto.Otto
}

func (e *JsEngine) Run(src interface{}) (fmt.Stringer, error) {
	return e.Otto.Run(src)
}

const cache_max = 500

var mutex = new(sync.Mutex)

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

func NewJsEngine() *JsEngine {
	return &JsEngine{otto.New()}
}

func GetJsEngine() *JsEngine {
	mutex.Lock()
	defer mutex.Unlock()
	if cache.Len() == 0 {
		isGenerate <- true
		return NewJsEngine()
	}
	return cache.Remove(cache.Front()).(*JsEngine)
}
