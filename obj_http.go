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
	client := &http.Client{Transport: nil}
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
	obj.Set("delete", func(call otto.FunctionCall) otto.Value {
		res, err := def_request(client, "DELETE", &call)
		return *build_response(js, res, err)
	})
	return obj
}

// http request method
func def_request(client *http.Client, method string, call *otto.FunctionCall) (*http.Response, error) {
	url := call.Argument(0).String()
	body := call.Argument(1).Object()
	header := call.Argument(2).Object()
	req, err := http.NewRequest(convert_url_body(method, url, body))
	if err != nil {
		return nil, err
	}
	if header != nil && header.Value().IsObject() {
		for _, key := range header.Keys() {
			value, _ := header.Get(key)
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
func convert_url_body(method string, url string, params *otto.Object) (string, string, io.Reader) {
	u := url
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
func params_string(params *otto.Object) *bytes.Buffer {
	buf := &bytes.Buffer{}
	if (params != nil) {
		for _, k := range params.Keys() {
			v, e := params.Get(k)
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
	obj.Set("get", func(key *string) *string {
		if key == nil {
			return nil
		}
		val := h.Get(*key)
		return &val
	})
	obj.Set("set", func(key *string, val *string) {
		if key == nil {
			return
		}
		if val == nil {
			h.Del(*key)
		} else {
			h.Add(*key, *val)
		}
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
