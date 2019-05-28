package main

import (
	"github.com/robertkrimen/otto"
	"os/exec"
	"runtime"
	"sync"
)

var rwMutex = new(sync.RWMutex)
var storage = make(map[string]string)

func GenerateObjJssp(js *JavaScript) *otto.Object {
	obj := js.CreateObjectValue().Object()
	obj.Set("exec", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 0 {
			return otto.UndefinedValue()
		}
		return *js.CreateAny(def_exec(call))
	})
	obj.Set("version", Server)
	obj.Set("os", runtime.GOOS)
	obj.Set("arch", runtime.GOARCH)
	obj.Set("storage", build_storage(js))
	return obj
}

// execute system command line
func def_exec(call otto.FunctionCall) string {
	var res []byte
	var err error
	if len(call.ArgumentList) == 1 {
		cmd := exec.Command(call.Argument(0).String())
		res, err = cmd.CombinedOutput()
	} else {
		ss := make([]string, 0)
		for _, v := range call.ArgumentList {
			ss = append(ss, v.String())
		}
		cmd := exec.Command(ss[0], ss[1:]...)
		res, err = cmd.CombinedOutput()
	}
	if err != nil {
		return err.Error()
	} else {
		return string(res)
	}
}

// build jssp.storage object
func build_storage(js *JavaScript) *otto.Object {
	obj := js.CreateObjectValue().Object()
	obj.Set("getItem", func(call otto.FunctionCall) otto.Value {
		rwMutex.RLock()
		defer rwMutex.RUnlock()
		if res, is := storage[call.Argument(0).String()]; !is {
			return otto.UndefinedValue()
		} else {
			return *js.CreateAny(res)
		}
	})
	obj.Set("setItem", func(call otto.FunctionCall) otto.Value {
		rwMutex.Lock()
		defer rwMutex.Unlock()
		storage[call.Argument(0).String()] = call.Argument(1).String()
		return otto.UndefinedValue()
	})
	obj.Set("removeItem", func(call otto.FunctionCall) otto.Value {
		rwMutex.Lock()
		defer rwMutex.Unlock()
		delete(storage, call.Argument(0).String())
		return otto.UndefinedValue()
	})
	obj.Set("size", func(call otto.FunctionCall) otto.Value {
		rwMutex.RLock()
		defer rwMutex.RUnlock()
		return *js.CreateAny(len(storage))
	})
	obj.Set("keys", func(call otto.FunctionCall) otto.Value {
		rwMutex.RLock()
		defer rwMutex.RUnlock()
		arr := js.CreateArray()
		o := arr.Object()
		for key := range storage {
			o.Call("push", key)
		}
		return *arr
	})
	obj.Set("clear", func(call otto.FunctionCall) otto.Value {
		rwMutex.Lock()
		defer rwMutex.Unlock()
		storage = make(map[string]string)
		return otto.UndefinedValue()
	})
	return obj
}
