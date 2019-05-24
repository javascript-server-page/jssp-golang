package main

import (
	"github.com/robertkrimen/otto"
	"os/exec"
	"runtime"
)

var storage = make(map[string]string)

func GenerateObjJssp(jse *JsEngine) *otto.Object {
	obj := jse.CreateObjectValue().Object()
	obj.Set("exec", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) == 0 {
			return otto.UndefinedValue()
		}
		return *jse.CreateAny(def_exec(call))
	})
	obj.Set("version", Server)
	obj.Set("os", runtime.GOOS)
	obj.Set("arch", runtime.GOARCH)
	obj.Set("storage", build_storage(jse))
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
func build_storage(jse *JsEngine) *otto.Object {
	obj := jse.CreateObjectValue().Object()
	obj.Set("getItem", func(call otto.FunctionCall) otto.Value {
		if res, is := storage[call.Argument(0).String()]; !is {
			return otto.UndefinedValue()
		} else {
			return *jse.CreateAny(res)
		}
	})
	obj.Set("setItem", func(call otto.FunctionCall) otto.Value {
		storage[call.Argument(0).String()] = call.Argument(1).String()
		return otto.UndefinedValue()
	})
	obj.Set("removeItem", func(call otto.FunctionCall) otto.Value {
		delete(storage, call.Argument(0).String())
		return otto.UndefinedValue()
	})
	obj.Set("size", func(call otto.FunctionCall) otto.Value {
		return *jse.CreateAny(len(storage))
	})
	obj.Set("keys", func(call otto.FunctionCall) otto.Value {
		arr := jse.CreateArray()
		o := arr.Object()
		for key := range storage {
			o.Call("push", key)
		}
		return *arr
	})
	obj.Set("clear", func(call otto.FunctionCall) otto.Value {
		storage = make(map[string]string)
		return otto.UndefinedValue()
	})
	return obj
}
