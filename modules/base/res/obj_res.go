package http

import (
	"bytes"
	"io"
	"jssp/config"
	"jssp/engine"
	"jssp/server/filetype"
	"os"
	"path"

	"github.com/dop251/goja"
	"github.com/valyala/fasthttp"
)

func init() {
	engine.ContextModules["res"] = GenerateObjRes
}

func GenerateObjRes(js *engine.VM, ctx *fasthttp.RequestCtx) *goja.Object {
	buf := &bytes.Buffer{}
	write := ctx.Response.BodyWriter()
	flush := func(size int) {
		if buf.Len() > size {
			io.Copy(write, buf)
			// ctx.Write(buf.Bytes())
			buf.Reset()
		}
	}
	js.Set("echo", func(vals ...goja.Value) {
		for _, e := range vals {
			str := e.String()
			if len(str) > 0 {
				flush(0)
				buf.WriteString(str)
			}
		}
	})
	js.Set("print", func(vals ...interface{}) {
		n := len(vals)
		if n == 0 {
			return
		}
		var obj *goja.Object
		if n == 1 {
			obj = js.ToValue(vals[0]).ToObject(js.Runtime)
		} else {
			obj = js.NewArray(vals...)
		}
		bty, err := obj.MarshalJSON()
		if err != nil {
			panic(js.NewGoError(err))
		}
		flush(0)
		buf.Write(bty)
	})
	dir, _ := path.Split(path.Join(config.Server.Dir, string(ctx.RequestURI())))
	js.Set("include", func(fname string) goja.Value {
		path, fi, ft, _ := filetype.GetFileInfo(path.Join(dir, fname))
		if ft == filetype.DIR {
			return nil
		}
		flush(0)
		f, err := os.Open(path)
		if err != nil {
			panic(js.NewGoError(err))
		}
		if ft == filetype.FILE {
			io.Copy(ctx.Response.BodyWriter(), f)
			return nil
		} else {
			ast, err := engine.GetAstByFile(path, fi.ModTime(), ft == filetype.JSSP, nil)
			if err != nil {
				panic(js.NewGoError(err))
			}
			result, err := js.RunProgram(ast)
			if err != nil {
				panic(js.NewGoError(err))
			}
			return result
		}
	})

	obj := js.CreateObject(nil)
	// obj.Set("header", build_editableheader(js, ctx.Response.Header))
	obj.Set("flush", flush)
	obj.Set("write", func(str string) {
		flush(0)
		ctx.WriteString(str)
	})
	return obj
}
