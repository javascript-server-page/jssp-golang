package main

import (
	"bytes"
	"github.com/robertkrimen/otto"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func GenerateObjHttp(jse *JsEngine) *otto.Object {
	client := &http.Client{}
	obj := jse.CreateObject()
	obj.Set("get", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(2)
		res, err := def_request(client, "GET", &url, &body, &header)
		return *build_response(jse, res, err)
	})
	obj.Set("head", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(2)
		res, err := def_request(client, "HEAD", &url, &body, &header)
		return *build_response(jse, res, err)
	})
	obj.Set("post", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(2)
		res, err := def_request(client, "POST", &url, &body, &header)
		return *build_response(jse, res, err)
	})
	obj.Set("put", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(2)
		res, err := def_request(client, "PUT", &url, &body, &header)
		return *build_response(jse, res, err)
	})
	return obj
}

func def_request(client *http.Client, method string, url, body, header *otto.Value) (*http.Response, error) {
	req, err := http.NewRequest(convert_url_body(method, url, body))
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

func build_response(jse *JsEngine, response *http.Response, err error) *otto.Value {
	val := jse.CreateObjectValue()
	obj := val.Object()
	if err != nil {
		obj.Set("status", -1)
		obj.Set("error", err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		obj.Set("status", response.StatusCode)
		obj.Set("body", string(data))
		obj.Set("header", build_header(jse, response.Header))
	}
	return val
}

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

func build_header(jse *JsEngine, h http.Header) *otto.Value {
	val := jse.CreateObjectValue()
	obj := val.Object()
	for k := range h {
		v := h.Get(k)
		obj.Set(k, v)
	}
	return val
}
