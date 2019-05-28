package main

import (
	"bytes"
	"github.com/robertkrimen/otto"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func GenerateObjHttp(js *JavaScript) *otto.Object {
	client := &http.Client{}
	obj := js.CreateObjectValue().Object()
	obj.Set("get", func(call otto.FunctionCall) otto.Value {
		res, err := def_request(client, "GET", &call)
		return *build_response(js, res, err)
	})
	obj.Set("head", func(call otto.FunctionCall) otto.Value {
		res, err := def_request(client, "HEAD", &call)
		return *build_response(js, res, err)
	})
	obj.Set("post", func(call otto.FunctionCall) otto.Value {
		res, err := def_request(client, "POST", &call)
		return *build_response(js, res, err)
	})
	obj.Set("put", func(call otto.FunctionCall) otto.Value {
		res, err := def_request(client, "PUT", &call)
		return *build_response(js, res, err)
	})
	return obj
}

// http request method
func def_request(client *http.Client, method string, call *otto.FunctionCall) (*http.Response, error) {
	url, body, header := call.Argument(0), call.Argument(1), call.Argument(2)
	req, err := http.NewRequest(convert_url_body(method, &url, &body))
	if err != nil {
		return nil, err
	}
	if header.IsObject() {
		h := header.Object()
		for _, key := range h.Keys() {
			value, _ := h.Get(key)
			req.Header.Add(key, value.String())
		}
	}
	return client.Do(req)
}

// build  jssp.Response object with *http.Response and error
func build_response(js *JavaScript, response *http.Response, err error) *otto.Value {
	val := js.CreateObjectValue()
	obj := val.Object()
	if err != nil {
		obj.Set("status", -1)
		obj.Set("error", err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		obj.Set("status", response.StatusCode)
		obj.Set("body", string(data))
		obj.Set("header", build_header(js, response.Header))
	}
	return val
}

// generate the http request parameters
func convert_url_body(method string, url, params *otto.Value) (string, string, io.Reader) {
	u := url.String()
	p := params_string(params)
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
func params_string(params *otto.Value) *bytes.Buffer {
	buf := &bytes.Buffer{}
	if !params.IsObject() {
		return buf
	}
	for _, k := range params.Object().Keys() {
		v, e := params.Object().Get(k)
		if e != nil {
			continue
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v.String())
		buf.WriteByte('&')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	return buf
}

// build an uneditable jssp.header object
func build_header(js *JavaScript, h http.Header) *otto.Object {
	obj := js.CreateObjectValue().Object()
	for k := range h {
		v := h.Get(k)
		obj.Set(k, v)
	}
	return obj
}

// build an editable jssp.header object
func build_editableheader(js *JavaScript, h http.Header) *otto.Object {
	obj := js.CreateObjectValue().Object()
	obj.Set("get", func(call otto.FunctionCall) otto.Value {
		val := call.Argument(0)
		if val.IsUndefined() {
			return val
		}
		return *js.CreateAny(h.Get(val.String()))
	})
	obj.Set("set", func(call otto.FunctionCall) otto.Value {
		key := call.Argument(0)
		if key.IsUndefined() || key.IsNull() {
			return key
		}
		val := call.Argument(1)
		if val.IsUndefined() || val.IsNull() {
			h.Del(key.String())
			return val
		}
		keystr := key.String()
		pre := *js.CreateAny(h.Get(keystr))
		h.Set(keystr, val.String())
		return pre
	})
	obj.Set("map", func(call otto.FunctionCall) otto.Value {
		val := js.CreateObjectValue()
		obj := val.Object()
		for k := range h {
			v := h.Get(k)
			obj.Set(k, v)
		}
		return *val
	})
	return obj
}
