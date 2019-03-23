package main

import (
	"bytes"
	"github.com/robertkrimen/otto"
	"io"
	"net/http"
)

func GenerateObjHttp(jse *JsEngine) *otto.Value {
	client := &http.Client{}
	val, obj := jse.CreateObject()
	obj.Set("get", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(1)
		res, err := def_request(client, "GET", &url, &body, &header)
		return *build_response(jse, res, err)
	})
	obj.Set("head", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(1)
		res, err := def_request(client, "HEAD", &url, &body, &header)
		return *build_response(jse, res, err)
	})
	obj.Set("post", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(1)
		res, err := def_request(client, "POST", &url, &body, &header)
		return *build_response(jse, res, err)
	})
	obj.Set("put", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(1)
		res, err := def_request(client, "PUT", &url, &body, &header)
		return *build_response(jse, res, err)
	})
	return val
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
	val, obj := jse.CreateObject()
	if err != nil {
		obj.Set("status", -1)
		obj.Set("error", err.Error())
	} else {
		obj.Set("status", response.StatusCode)
	}
	return val
}

func convert_url_body(method string, url, params *otto.Value) (string, string, io.Reader) {
	return method, "", nil
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
	buf.Truncate(buf.Len() - 1)
	return buf
}
