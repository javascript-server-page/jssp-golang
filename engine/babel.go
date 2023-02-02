package engine

import (
	"errors"
	"jssp/tool"
	"sync"

	"jssp/config"

	"github.com/dop251/goja"
)

const CompileES6Function = `;
function compile(code) {
	try {
		var obj = Babel.transform(code, { presets: ['es6'] });
		obj.isError = false;
		return obj;
	} catch(e) {
		e.isError = true;
		return e;
	}
}`

const CompileTSFunction = `;
function compile(code) {
	try {
		var obj = Babel.transform(code, { presets: ['typescript'], filename: '.ts' });
		obj.isError = false;
		return obj;
	} catch(e) {
		e.isError = true;
		return e;
	}
}`

var babelAst *goja.Program
var babelPool *sync.Pool

func init() {
	if !config.Babel.Enable {
		return
	}
	js, err := tool.GetResource(config.Babel.Path)
	if err != nil {
		panic("load Babel error:" + err.Error())
	}
	var CompileFunction string
	if config.Babel.Ts {
		CompileFunction = CompileTSFunction
	} else {
		CompileFunction = CompileES6Function
	}
	babelAst, err = goja.Compile("", string(js)+CompileFunction, true)
	if err != nil {
		panic(err)
	}
	babelPool = &sync.Pool{
		New: func() interface{} {
			babelVm := goja.New()
			//babelVm.Set("console", babelVm.CreateObject(nil))
			_, err := babelVm.RunProgram(babelAst)
			if err != nil {
				panic(err)
			}
			var babelCompile func(string) *goja.Object
			babelVm.ExportTo(babelVm.Get("compile"), &babelCompile)
			return babelCompile
		},
	}
	babelPool.Put(babelPool.Get())
}

func BabelCompile(jscode string) (string, error) {
	if babelPool == nil {
		return jscode, nil
	}
	babelCompile := babelPool.Get().(func(string) *goja.Object)
	defer babelPool.Put(babelCompile)
	res := babelCompile(jscode)
	if res.Get("isError").Export().(bool) {
		return "", errors.New(res.Get("message").String())
	}
	return res.Get("code").String(), nil
}
