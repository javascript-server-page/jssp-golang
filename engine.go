package main

import (
	"container/list"
	"fmt"
	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/parser"
	"net/http"
	"sync"
)

type JavaScript struct {
	*otto.Otto
}

func (js *JavaScript) Parse(src []byte) (interface{}, error) {
	return parser.ParseFile(nil, "", src, 0)
}

func (js *JavaScript) Run(src interface{}) (fmt.Stringer, error) {
	return js.Otto.Run(src)
}

func (js JavaScript) CreateObjectValue() *otto.Value {
	val, _ := js.Otto.Run("({})")
	return &val
}

func (js JavaScript) CreateArray() *otto.Value {
	val, _ := js.Otto.Run("[]")
	return &val
}

func (js JavaScript) CreateError(err error) *otto.Value {
	if err != nil {
		ce := js.MakeCustomError("Jssp", err.Error())
		return &ce
	}
	return &otto.Value{}
}

func (js JavaScript) CreateAny(any interface{}) *otto.Value {
	if any == nil {
		null := otto.NullValue()
		return &null
	}
	v, err := js.ToValue(any)
	if err != nil {
		re := js.MakeRangeError(err.Error())
		return &re
	}
	return &v
}

func (js JavaScript) isError(val *otto.Value) bool {
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

func NewJsEngine() *JavaScript {
	js := &JavaScript{otto.New()}
	js.Set("http", GenerateObjHttp(js))
	js.Set("jsdo", GenerateObjJsdo(js))
	js.Set("jssp", GenerateObjJssp(js))
	return js
}

func GetJsEngine() *JavaScript {
	mutex.Lock()
	defer mutex.Unlock()
	if cache.Len() == 0 {
		isGenerate <- true
		return NewJsEngine()
	}
	return cache.Remove(cache.Front()).(*JavaScript)
}

func GenerateJsspEnv(s *JsspServer, w http.ResponseWriter, r *http.Request) *JavaScript {
	js := GetJsEngine()
	js.Set("file", GenerateObjFile(js, s.paras.Dir+r.RequestURI))
	js.Set("req", GenerateObjReq(js, r))
	js.Set("res", GenerateObjRes(js, w))
	return js
}
