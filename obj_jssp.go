package main

import (
	"github.com/robertkrimen/otto"
	"os/exec"
	"runtime"
)

func GenerateObjJssp(jse *JsEngine) *otto.Object {
	obj := jse.CreateObject()
	obj.Set("exec", func(call otto.FunctionCall) otto.Value {
		var res []byte
		var err error
		switch len(call.ArgumentList) {
		case 0:
			return otto.UndefinedValue()
		case 1:
			cmd := exec.Command(call.Argument(0).String())
			res, err = cmd.CombinedOutput()
		default:
			ss := make([]string, 0)
			for _, v := range call.ArgumentList {
				ss = append(ss, v.String())
			}
			cmd := exec.Command(ss[0], ss[1:]...)
			res, err = cmd.CombinedOutput()
		}
		if err != nil {
			return *jse.CreateString(err.Error())
		} else {
			return *jse.CreateString(string(res))
		}
	})
	obj.Set("version", Server)
	obj.Set("os", runtime.GOOS)
	obj.Set("arch", runtime.GOARCH)
	return obj
}
