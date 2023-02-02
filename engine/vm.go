package engine

import (
	"jssp/config"
	"math/rand"
	"sync"

	"github.com/dop251/goja"
	"github.com/dop251/goja/parser"
	"github.com/valyala/fasthttp"
)

type GetDefaultJsValue func(*VM) *goja.Object
type GetContextJsValue func(*VM, *fasthttp.RequestCtx) *goja.Object

var DefaultModules map[string]GetDefaultJsValue = make(map[string]GetDefaultJsValue)
var ContextModules map[string]GetContextJsValue = make(map[string]GetContextJsValue)

func init() {
	if config.Vmpool.Enable {
		pool = new(VMPool)
		pool.data = new(sync.Map)
		go pool.Reset()
	}
}

var pool *VMPool

type VM struct {
	*goja.Runtime
}

type VMPool struct {
	data *sync.Map
}

func (p *VMPool) Reset() {
	for i := 0; i < config.Vmpool.Size; i++ {
		pool.data.Store(i, NewVM())
	}
}

func (p *VMPool) Get(i int) *VM {
	index := rand.Intn(config.Vmpool.Size)
	vm, loaded := p.data.LoadAndDelete(index)
	if loaded {
		go p.Store(index)
		return vm.(*VM)
	}
	if i > config.Vmpool.Retry {
		// println("new vm retry size is", i)
		return NewVM()
	}
	return p.Get(i + 1)
}

func (p *VMPool) Store(index int) {
	p.data.Store(index, NewVM())
}

func NewVM() *VM {
	vm := &VM{goja.New()}
	vm.SetParserOptions(parser.WithDisableSourceMaps)
	for name, getDefaultJsValue := range DefaultModules {
		vm.Set(name, getDefaultJsValue(vm))
	}
	return vm
}

func NewVMByContext(ctx *fasthttp.RequestCtx) *VM {
	var vm *VM
	if pool == nil {
		vm = NewVM()
	} else {
		vm = pool.Get(1)
	}
	for name, getContextJsValue := range ContextModules {
		vm.Set(name, getContextJsValue(vm, ctx))
	}
	return vm
}
