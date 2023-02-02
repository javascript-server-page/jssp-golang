package http

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"jssp/engine"
	"net/http"
	"reflect"
	"strings"

	"github.com/dop251/goja"
)

func init() {
	engine.DefaultModules["http"] = GenerateObjHttp
}

var METHODS = []string{
	"OPTIONS",
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"TRACE",
	"CONNECT",
}

func GenerateObjHttp(js *engine.VM) *goja.Object {
	client := &http.Client{Transport: nil}
	obj := js.CreateObject(nil)
	for _, m := range METHODS {
		obj.Set(strings.ToLower(m), func(vals ...goja.Value) goja.Value {
			if len(vals) == 0 {
				panic(js.NewGoError(errors.New("")))
			}
			res, err := build_request(js, client, m, vals...)
			return build_response(js, res, err)
		})
	}
	obj.Set("req", func(obj *goja.Object) {
		if obj == nil {
			panic(js.NewGoError(errors.New("")))
		}
	})
	return obj
}

// http request method
func build_request(js *engine.VM, client *http.Client, method string, vals ...goja.Value) (*http.Response, error) {
	url := vals[0].String()
	var body goja.Value
	if len(vals) >= 2 {
		body = vals[1]
	}
	req, err := http.NewRequest(convert_url_body(js, method, url, body))
	if err != nil {
		return nil, err
	}
	if len(vals) >= 3 {
		header := vals[2].ToObject(js.Runtime)
		for _, key := range header.Keys() {
			val := header.Get(key)
			req.Header.Add(key, val.String())
		}
	}
	return client.Do(req)
}

// build  jssp.Response object with *http.Response and error
func build_response(js *engine.VM, res *http.Response, err error) goja.Value {
	obj := js.CreateObject(nil)
	if err != nil {
		obj.Set("status", -1)
		obj.Set("error", js.NewGoError(err))
	} else {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			obj.Set("status", -2)
			obj.Set("error", js.NewGoError(err))
		} else {
			obj.Set("status", res.StatusCode)
			obj.Set("body", string(data))
			obj.Set("header", build_header(js, res.Header))
		}
	}
	return obj
}

// generate the http request parameters
func convert_url_body(js *engine.VM, method string, url string, params goja.Value) (string, string, io.Reader) {
	u := url
	p := params_string(js, params)
	if method == "POST" || method == "PUT" {
		return method, u, p
	}
	if strings.LastIndex(u, "?") > 0 {
		u = u + "&" + p.String()
	} else {
		u = u + "?" + p.String()
	}
	return method, u, strings.NewReader("")
}

// convert key-value pairs to http parameters
func params_string(js *engine.VM, params goja.Value) *bytes.Buffer {
	buf := &bytes.Buffer{}
	if params == nil {
		return buf
	}
	kind := params.ExportType().Kind()
	if kind == reflect.Map || kind == reflect.Array || kind == reflect.Slice {
		ps := params.ToObject(js.Runtime)
		for _, k := range ps.Keys() {
			v := ps.Get(k)
			buf.WriteString(k)
			buf.WriteByte('=')
			buf.WriteString(v.String())
			buf.WriteByte('&')
		}
		if buf.Len() > 0 {
			buf.Truncate(buf.Len() - 1)
		}
	} else {
		buf.WriteString(params.String())
	}
	return buf
}

// build an uneditable jssp.header object
func build_header(js *engine.VM, h http.Header) *goja.Object {
	obj := js.CreateObject(nil)
	for k := range h {
		v := h.Get(k)
		obj.Set(k, v)
	}
	return obj
}

/*
// build an editable jssp.header object
func build_editableheader(js *engine.VM, h http.Header) *goja.Object {
	obj := js.CreateObject(nil)
	obj.Set("get", h.Get)
	obj.Set("set", h.Set)
	obj.Set("del", h.Del)
	obj.Set("clone", func() *goja.Object { return build_editableheader(js, h.Clone()) })
	obj.Set("keys", func() goja.Value {
		arr := make([]string, 0)
		for k := range h {
			arr = append(arr, k)
		}
		return js.ToValue(arr)
	})
	return obj
}
*/
