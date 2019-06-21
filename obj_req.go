package main

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

func GenerateObjReq(js *JavaScript, r *http.Request, w http.ResponseWriter, s *Session) *otto.Object {
	r.ParseForm()
	obj := js.CreateObjectValue().Object()
	obj.Set("header", build_header(js, r.Header))
	obj.Set("host", r.Host)
	obj.Set("method", r.Method)
	obj.Set("path", r.URL.Path)
	obj.Set("proto", r.Proto)
	obj.Set("remoteAddr", r.RemoteAddr)
	obj.Set("cookie", build_cookie(js, r, w))
	obj.Set("session", build_session(js, s))
	obj.Set("parm", r.Form.Get)
	obj.Set("file", "")
	return obj
}

// build req.cookie object
func build_cookie(js *JavaScript, r *http.Request, w http.ResponseWriter) *otto.Object {
	obj := js.CreateObjectValue().Object()
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
		c := &http.Cookie{Name: *key,}
		if val == nil {
			c.MaxAge = -1
		} else {
			c.Value = *val
		}
		http.SetCookie(w, c)
	})
	obj.Set("map", func(call otto.FunctionCall) otto.Value {
		val := js.CreateObjectValue()
		obj := val.Object()
		for _, k := range r.Cookies() {
			obj.Set(k.Name, k.Value)
		}
		return *val
	})
	return obj
}

// build req.session object
func build_session(js *JavaScript, s *Session) *otto.Object {
	obj := js.CreateObjectValue().Object()
	obj.Set("get", func(key *string) *otto.Value {
		if key == nil {
			return nil
		}
		val, ok := s.data[*key]
		if ok {
			return val
		}
		return nil
	})
	obj.Set("set", func(key *string, val *otto.Value) {
		if key == nil {
			return
		}
		s.data[*key] = val
	})
	obj.Set("start", s.mutex.Lock)
	obj.Set("close", s.mutex.Unlock)
	return obj
}
