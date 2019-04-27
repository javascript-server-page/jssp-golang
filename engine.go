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

func (e *JsEngine) CreateObject() *otto.Object {
	val, _ := e.Otto.Run("({})")
	return val.Object()
}

func (e *JsEngine) CreateObjectValue() *otto.Value {
	val, _ := e.Otto.Run("({})")
	return &val
}

func (e *JsEngine) CreateString(s string) *otto.Value {
	v, _ := e.ToValue(s)
	return &v
}

func (e *JsEngine) CreateArray() *otto.Value {
	val, _ := e.Otto.Run("[]")
	return &val
}

func (e *JsEngine) CreateError(err error) *otto.Value {
	if err != nil {
		val, err := e.Call("Error", e, err.Error())
		if err == nil {
			return &val
		}
	}
	return &otto.Value{}
}

func (e *JsEngine) CreateAny(any interface{}) *otto.Value {
	if any == nil {
		null := otto.NullValue()
		return &null
	}
	v, err := e.ToValue(any)
	if err != nil {
		return e.CreateError(err)
	}
	return &v
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
	jse.Set("file", GenerateObjFile(jse, "."+r.RequestURI))
	jse.Set("req", GenerateObjReq(jse, r))
	jse.Set("res", GenerateObjRes(jse, w))
	return jse
}
