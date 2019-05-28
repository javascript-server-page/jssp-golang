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

type Engine struct {
	max        int
	mutex      *sync.Mutex
	cache      *list.List
	isGenerate chan bool
}

func NewEngine() *Engine {
	e := &Engine{500, new(sync.Mutex), list.New(), make(chan bool)}
	go e.generate()
	return e
}

func (e *Engine) NewJavaScript() *JavaScript {
	js := &JavaScript{otto.New()}
	js.Set("http", GenerateObjHttp(js))
	js.Set("jsdo", GenerateObjJsdo(js))
	js.Set("jssp", GenerateObjJssp(js))
	return js
}

func (e *Engine) GetJavaScript() *JavaScript {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.cache.Len() == 0 {
		e.isGenerate <- true
		return e.NewJavaScript()
	}
	return e.cache.Remove(e.cache.Front()).(*JavaScript)
}

func (e *Engine) GenJsspEnv(s *JsspServer, w http.ResponseWriter, r *http.Request) *JavaScript {
	js := e.GetJavaScript()
	js.Set("file", GenerateObjFile(js, s.set.Dir+r.RequestURI))
	js.Set("req", GenerateObjReq(js, r))
	js.Set("res", GenerateObjRes(js, w))
	return js
}

func (e *Engine) generate() {
	for {
		<-e.isGenerate
		for e.cache.Len() < e.max {
			e.cache.PushBack(e.NewJavaScript())
		}
	}
}
