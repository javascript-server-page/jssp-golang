package main

import (
	"github.com/robertkrimen/otto"
	"os/exec"
	"runtime"
)

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
