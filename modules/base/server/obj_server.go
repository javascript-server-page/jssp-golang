package server

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"jssp/engine"
	"os/exec"
	"runtime"
	"sync"

	"github.com/dop251/goja"
)

func init() {
	engine.DefaultModules["server"] = GenerateObjJssp
}

var storage = new(sync.Map)

func GenerateObjJssp(js *engine.VM) *goja.Object {
	obj := js.CreateObject(nil)

	obj.Set("exec", func(vals ...goja.Value) goja.Value {
		var cmd *exec.Cmd
		switch len(vals) {
		case 0:
			return goja.Undefined()
		case 1:
			cmd = exec.Command(vals[0].String())
		default:
			ss := make([]string, 0)
			for _, v := range vals {
				ss = append(ss, v.String())
			}
			cmd = exec.Command(ss[0], ss[1:]...)
		}
		res, err := cmd.CombinedOutput()
		if err != nil {
			panic(js.NewGoError(err))
		}
		return js.ToValue(string(res))
	})
	obj.Set("os", runtime.GOOS)
	obj.Set("arch", runtime.GOARCH)
	obj.Set("cypto", build_jssp_cypto(js))
	obj.Set("storage", build_jssp_storage(js))
	return obj
}

// build jssp.storage object
func build_jssp_storage(js *engine.VM) *goja.Object {
	obj := js.CreateObject(nil)
	obj.Set("getItem", func(key string) *string {
		if res, ok := storage.Load(key); !ok {
			return nil
		} else {
			str := res.(string)
			return &str
		}
	})
	obj.Set("setItem", func(key, val string) {
		storage.Store(key, val)
	})
	obj.Set("removeItem", func(key string) {
		storage.Delete(key)
	})
	obj.Set("clear", func() {
		storage.Range(func(key, value interface{}) bool {
			storage.Delete(key)
			return true
		})
	})
	return obj
}

// build jssp.cypto object
func build_jssp_cypto(js *engine.VM) *goja.Object {
	obj := js.CreateObject(nil)
	obj.Set("md5", func(key string) string {
		return fmt.Sprintf("%x", md5.Sum([]byte(key)))
	})
	obj.Set("sha1", func(key string) string {
		return fmt.Sprintf("%x", sha1.Sum([]byte(key)))
	})
	return obj
}
