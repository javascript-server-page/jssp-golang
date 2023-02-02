package http

import (
	"io"
	"jssp/engine"
	"net/http"
	"os"
	"path"

	"github.com/dop251/goja"
	"github.com/valyala/fasthttp"
)

func init() {
	engine.ContextModules["req"] = GenerateObjReq
}

func GenerateObjReq(js *engine.VM, ctx *fasthttp.RequestCtx) *goja.Object {
	// r.ParseForm()
	obj := js.CreateObject(nil)
	obj.Set("header", build_header(js, ctx.Request.Header))
	obj.Set("host", string(ctx.Host()))
	obj.Set("method", ctx.Method())
	obj.Set("path", ctx.URI().String())
	obj.Set("remoteAddr", ctx.RemoteAddr().String())
	obj.Set("remoteIP", ctx.RemoteIP().String())

	// obj.Set("cookie", build_cookie(js, r, w))
	// obj.Set("session", build_session(js, s))
	// obj.Set("parm", ctx.QueryArgs())
	// obj.Set("file", func(key string) goja.Value {
	// 	return build_upload_file(js, r, key).Value()
	// })
	return obj
}

func build_header(js *engine.VM, header fasthttp.RequestHeader) goja.Value {
	// header
	return nil
}

// build req.cookie object
func build_cookie(js *engine.VM, r *http.Request, w http.ResponseWriter) *goja.Object {
	obj := js.CreateObject(nil)
	obj.Set("get", func(key *string) *string {
		if key == nil {
			return key
		}
		c, err := r.Cookie(*key)
		if err != nil {
			return nil
		}
		return &c.Value
	})
	obj.Set("set", func(key *string, val *string) {
		if key == nil {
			return
		}
		c := &http.Cookie{Name: *key}
		if val == nil {
			c.MaxAge = -1
		} else {
			c.Value = *val
		}
		http.SetCookie(w, c)
	})
	obj.Set("map", func(call goja.FunctionCall) *goja.Object {
		obj := js.CreateObject(nil)
		for _, k := range r.Cookies() {
			obj.Set(k.Name, k.Value)
		}
		return obj
	})
	return obj
}

// build req.session object
// func build_session(js *engine.VM, s *Session) *goja.Object {
// 	obj := js.CreateObjectValue().Object()
// 	obj.Set("get", func(key *string) *goja.Value {
// 		if key == nil {
// 			return nil
// 		}
// 		val, ok := s.data[*key]
// 		if ok {
// 			return val
// 		}
// 		return nil
// 	})
// 	obj.Set("set", func(key *string, val *goja.Value) {
// 		if key == nil {
// 			return
// 		}
// 		s.data[*key] = val
// 	})
// 	obj.Set("start", s.mutex.Lock)
// 	obj.Set("close", s.mutex.Unlock)
// 	return obj
// }

// build req.upload_file object
func build_upload_file(js *engine.VM, r *http.Request, key string) *goja.Object {
	obj := js.CreateObject(nil)
	src, header, err := r.FormFile(key)
	if err != nil {
		return js.NewGoError(err)
	}

	obj.Set("size", header.Size)
	obj.Set("name", header.Filename)
	obj.Set("move", func(dir string, new_name *string) goja.Value {
		var new_path string
		if new_name == nil {
			new_path = path.Join(dir, header.Filename)
		} else {
			new_path = path.Join(dir, *new_name)
		}
		dst, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return js.NewGoError(err)
		}
		_, err = io.Copy(dst, src)
		if err != nil {
			return js.NewGoError(err)
		}
		err = src.Close()
		if err != nil {
			return js.NewGoError(err)
		}
		err = dst.Close()
		if err != nil {
			return js.NewGoError(err)
		}
		return goja.Null()
	})
	return obj
}
