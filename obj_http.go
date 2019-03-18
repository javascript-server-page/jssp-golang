package main

import (
	"github.com/robertkrimen/otto"
	"io"
	"net/http"
)

func GenerateObjHttp(jse *JsEngine) *otto.Value {
	client := &http.Client{}
	val, obj := jse.CreateObject()
	obj.Set("get", func(call otto.FunctionCall) otto.Value {
		url, body, header := call.Argument(0), call.Argument(1), call.Argument(1)
		return *def_request(client, "GET", &url, &body, &header)
	})
	return val
}

func def_request(client *http.Client, method string, url, body, header *otto.Value) *otto.Value {
	req, err := http.NewRequest(convert_url_body(method, url, body))
	if err != nil {
		return build_response(nil, err)
	}
	if header.IsObject() {
		h := header.Object()
		for _, key := range h.Keys() {
			value, _ := h.Get(key)
			req.Header.Add(key, value.String())
		}
	}
	return build_response(client.Do(req))
}

func build_response(response *http.Response, err error) *otto.Value {
	return nil
}

func convert_url_body(method string, url, params *otto.Value) (string, string, io.Reader) {
	return method, "", nil
}
