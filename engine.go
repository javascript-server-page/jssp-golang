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

func (e *JsEngine) CreateObjectValue() *otto.Value {
	val, _ := e.Otto.Run("({})")
	return &val
}

func (e *JsEngine) CreateArray() *otto.Value {
	val, _ := e.Otto.Run("[]")
	return &val
}

func (e *JsEngine) CreateError(err error) *otto.Value {
	ce := e.MakeCustomError("Jssp", err.Error())
	return &ce
}

func (e *JsEngine) CreateAny(any interface{}) *otto.Value {
	if any == nil {
		null := otto.NullValue()
		return &null
	}
	v, err := e.ToValue(any)
	if err != nil {
		re := e.MakeRangeError(err.Error())
		return &re
	}
	return &v
}

func (e *JsEngine) isError(val *otto.Value) bool {
	if val == nil {
		return false
	}
	return val.Class() == "Error"
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

func GenerateJsspEnv(s *JsspServer, w http.ResponseWriter, r *http.Request) *JsEngine {
	jse := GetJsEngine()
	jse.Set("file", GenerateObjFile(jse, s.paras.Dir+r.RequestURI))
	jse.Set("req", GenerateObjReq(jse, r))
	jse.Set("res", GenerateObjRes(jse, w))
	return jse
}
