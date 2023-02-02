package console

import (
	"bytes"
	"jssp/engine"
	"jssp/log"
	"time"

	"github.com/dop251/goja"
)

func init() {
	engine.DefaultModules["console"] = GenerateObjConsole
}

func GenerateObjConsole(js *engine.VM) *goja.Object {
	obj := js.CreateObject(nil)
	obj.Set("log", func(vals ...goja.Value) {
		if len(vals) == 0 {
			return
		}
		vs := make([]*goja.Object, len(vals))
		for i, v := range vals {
			vs[i] = v.ToObject(js.Runtime)
		}
		log.SInfo(&logValue{time.Now(), vs})
	})
	obj.Set("error", func(vals ...goja.Value) {
		if len(vals) == 0 {
			return
		}
		vs := make([]*goja.Object, len(vals))
		for i, v := range vals {
			vs[i] = v.ToObject(js.Runtime)
		}
		log.SError(&logValue{time.Now(), vs})
	})
	return obj
}

type logValue struct {
	time time.Time
	vals []*goja.Object
}

func (lv *logValue) Time() time.Time {
	return lv.time
}

func (lv *logValue) String() string {
	buf := &bytes.Buffer{}
	for _, v := range lv.vals {
		data, err := v.MarshalJSON()
		if err != nil {
			buf.WriteString(err.Error())
		} else {
			buf.Write(data)
		}
		buf.WriteRune(' ')
	}
	return buf.String()
}
