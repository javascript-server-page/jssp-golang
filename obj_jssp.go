package main

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
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
		return *js.CreateAny(def_jssp_exec(&call))
	})
	obj.Set("version", Server)
	obj.Set("os", runtime.GOOS)
	obj.Set("arch", runtime.GOARCH)
	obj.Set("cypto", build_jssp_cypto(js))
	obj.Set("storage", build_jssp_storage(js))
	return obj
}

// execute system command line
func def_jssp_exec(call *otto.FunctionCall) string {
	var cmd *exec.Cmd
	switch len(call.ArgumentList) {
	case 0:
		return ""
	case 1:
		cmd = exec.Command(call.Argument(0).String())
	default:
		ss := make([]string, 0)
		for _, v := range call.ArgumentList {
			ss = append(ss, v.String())
		}
		cmd = exec.Command(ss[0], ss[1:]...)
	}
	res, err := cmd.CombinedOutput()
	if err != nil {
		return err.Error()
	} else {
		return string(res)
	}
}

// build jssp.storage object
func build_jssp_storage(js *JavaScript) *otto.Object {
	obj := js.CreateObjectValue().Object()
	obj.Set("getItem", func(key string) *string {
		rwMutex.RLock()
		defer rwMutex.RUnlock()
		if res, is := storage[key]; !is {
			return nil
		} else {
			return &res
		}
	})
	obj.Set("setItem", func(key, val string) {
		rwMutex.Lock()
		defer rwMutex.Unlock()
		storage[key] = val
	})
	obj.Set("removeItem", func(key string) otto.Value {
		rwMutex.Lock()
		defer rwMutex.Unlock()
		delete(storage, key)
		return otto.UndefinedValue()
	})
	obj.Set("size", func() int {
		rwMutex.RLock()
		defer rwMutex.RUnlock()
		return len(storage)
	})
	obj.Set("keys", func() []string {
		rwMutex.RLock()
		defer rwMutex.RUnlock()
		i := 0
		keys := make([]string, len(storage))
		for key := range storage {
			keys[i] = key
			i++
		}
		return keys
	})
	obj.Set("clear", func() {
		rwMutex.Lock()
		defer rwMutex.Unlock()
		storage = make(map[string]string)
	})
	return obj
}

// build jssp.cypto object
func build_jssp_cypto(js *JavaScript) *otto.Object {
	obj := js.CreateObjectValue().Object()
	obj.Set("md5", func(key string) string {
		return fmt.Sprintf("%x", md5.Sum([]byte(key)))
	})
	obj.Set("sha1", func(key string) string {
		return fmt.Sprintf("%x", sha1.Sum([]byte(key)))
	})
	return obj
}
