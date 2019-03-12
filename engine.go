package main

import (
	"container/list"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
	"net/http"
	"sync"
)

type JsEngine struct {
	*otto.Otto
}

func (e *JsEngine) Parse(src []byte) (interface{}, error) {
	return parser.ParseFile(nil, "", src, 0)
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
	js := &JsEngine{otto.New()}
	js.Set("file", GenerateObjFile(js))
	js.Set("http", GenerateObjHttp(js))
	js.Set("jsdo", GenerateObjJsdo(js))
	js.Set("jssp", GenerateObjJssp(js))
	return js
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

func GenerateJsspEnv(w http.ResponseWriter, r *http.Request) *JsEngine {
	jse := GetJsEngine()
	jse.Set("req", GenerateObjReq(jse, r))
	jse.Set("res", GenerateObjRes(jse, w))
	jse.Set("echo", func(call otto.FunctionCall) otto.Value {
		for _, e := range call.ArgumentList {
			w.Write([]byte(e.String()))
		}
		return otto.Value{}
	})
	return jse
}
