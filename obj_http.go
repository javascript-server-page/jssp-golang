package main

import (
	"github.com/robertkrimen/otto"
	"io"
	"net/http"
)

func GenerateObjHttp(jse *JsEngine) *otto.Object {
	obj := jse.CreateObject()
	return obj
}

func def_request(client *http.Client, method string, url, body, header *otto.Object) *otto.Object {
	req, err := http.NewRequest(convert_url_body(method, url, body))
	if err != nil {
		return build_response(nil, err)
	}
	for _, key := range header.Keys() {
		value, _ := header.Get(key)
		req.Header.Add(key, value.String())
	}
	header.Keys()
	return build_response(client.Do(req))
}

func build_response(response *http.Response, err error) *otto.Object {
	return nil
}

func convert_url_body(method string, url, params *otto.Object) (string, string, io.Reader) {
	return method, "", nil
}
