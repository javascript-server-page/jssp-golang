package engine

import (
	"bytes"
	"io/ioutil"
	"jssp/config"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/bluele/gcache"
	"github.com/dop251/goja"
)

var astCache gcache.Cache

func init() {
	if config.Astpool.Enable {
		astCache = gcache.New(config.Astpool.Size * 1024 * 1024).EvictType(config.Astpool.Mode).Build()
	}
}

type ast_item struct {
	time.Time
	*goja.Program
}

func getAstByCache(key string, mod time.Time) *goja.Program {
	ast, err := astCache.Get(key)
	if err != nil {
		return nil
	}
	if ast.(*ast_item).Time != mod {
		astCache.Remove(key)
		return nil
	}
	return ast.(*ast_item).Program
}

func setAstByCache(key string, ast *goja.Program, mod time.Time) {
	astCache.Set(key, &ast_item{mod, ast})
}

var fileLockerMap = new(sync.Map)

func getFileLocker(path string) sync.Locker {
	if locker, ok := fileLockerMap.Load(path); ok {
		return locker.(sync.Locker)
	}
	locker, _ := fileLockerMap.LoadOrStore(path, &sync.Mutex{})
	return locker.(sync.Locker)
}

func GetAstByFile(path string, modtime time.Time, isTemplate bool, f http.File) (*goja.Program, error) {
	if astCache == nil {
		return getAstByFile(path, isTemplate, f)
	}
	ast := getAstByCache(path, modtime)
	if ast != nil {
		return ast, nil
	}
	locker := getFileLocker(path)
	locker.Lock()
	defer locker.Unlock()
	ast = getAstByCache(path, modtime)
	if ast != nil {
		return ast, nil
	}
	ast, err := getAstByFile(path, isTemplate, f)
	if err != nil {
		return nil, err
	}
	setAstByCache(path, ast, modtime)
	return ast, nil
}

func getAstByFile(path string, isTemplate bool, f http.File) (*goja.Program, error) {
	if f == nil {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		f = file
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if isTemplate {
		data = jssp_jsjs(data)
	}
	res, err := BabelCompile(string(data))
	if err != nil {
		return nil, err
	}
	return goja.Compile("", res+";res.flush(0)", true)
}

func jssp_jsjs(data []byte) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString(`echo("`)
	isJsjs := false
	isPrint := false
	for i, n := 0, len(data); i < n; i++ {
		c := data[i]
		if isJsjs {
			if c == '%' && data[i+1] == '>' {
				if isPrint {
					isPrint = false
					buf.WriteString(`);echo("`)
				} else {
					buf.WriteString(`;echo("`)
				}
				i++
				isJsjs = false
			} else {
				buf.WriteByte(c)
			}
		} else {
			if c == '<' && data[i+1] == '%' {
				buf.WriteString(`");`)
				if data[i+2] == '=' {
					i++
					isPrint = true
					buf.WriteString(`echo(`)
				}
				i++
				isJsjs = true
			} else {
				switch c {
				case '\n':
					buf.WriteString("\\n")
					// buf.WriteString(`\n");`)
					// buf.WriteByte(c)
					// buf.WriteString(`echo("`)
				case '\r':
					continue
				case '"':
					buf.WriteString(`\"`)
				default:
					buf.WriteByte(c)
				}
			}
		}
	}
	buf.WriteString(`");`)
	return buf.Bytes()
}
