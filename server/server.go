package server

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"jssp/config"
	"jssp/engine"
	"jssp/log"
	_ "jssp/modules"
	"jssp/server/filetype"

	"github.com/valyala/fasthttp"
)

type JsspServer struct {
}

// init JsspServer
func (s *JsspServer) Init() {
}

// handler func
func (s *JsspServer) ServeAll(ctx *fasthttp.RequestCtx) {
	startTime := time.Now()
	path, fi, ft, err := filetype.GetFileInfo(string(ctx.Path()))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusNotFound)
		return
	}
	f, err := filetype.GetFile(path, ft)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusForbidden)
		return
	}
	s.header(ctx, fi)
	if ft == filetype.DIR {
		if dirs, err := f.Readdir(0); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusNotFound)
		} else {
			ctx.WriteString("<pre>\n")
			path := ctx.Path()
			b := path[len(path)-1] == '/'
			for _, dir := range dirs {
				if b {
					ctx.WriteString("<a href='" + dir.Name())
				} else {
					ctx.WriteString("<a href='" + fi.Name() + "/" + dir.Name())
				}
				if dir.IsDir() {
					ctx.WriteString("/'>" + dir.Name() + "</a>\n")
				} else {
					ctx.WriteString("'>" + dir.Name() + "</a>\n")
				}
			}
			ctx.WriteString("</pre>\n")
		}
	} else if ft == filetype.FILE {
		io.Copy(ctx.Response.BodyWriter(), f)
	} else {
		ast, err := engine.GetAstByFile(string(ctx.Path()), fi.ModTime(), ft == filetype.JSSP, f)
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		js := engine.NewVMByContext(ctx)
		timer := time.AfterFunc(time.Duration(config.Server.Timeout)*time.Second, func() {
			js.Interrupt("render timeout!")
		})
		t3 := time.Now()
		_, err = js.RunProgram(ast)
		fmt.Println("RunProgram", time.Since(t3))
		if err != nil {
			ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
			return
		}
		timer.Stop()
	}
	f.Close()
	// fmt.Println(string(ctx.Request.Header.Header()))
	log.Access(startTime, string(ctx.Method()), string(ctx.Path()), ctx.RemoteIP().String())
}

func (s *JsspServer) header(ctx *fasthttp.RequestCtx, fi fs.FileInfo) {
	Path := string(ctx.Path())
	if fi.IsDir() || strings.HasSuffix(Path, "/") {
		ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	} else {
		txt := mime.TypeByExtension(path.Ext(Path))
		ctx.Response.Header.Set("Content-Type", txt)
	}
	ctx.Response.Header.Set("Last-Modified", fi.ModTime().Format(http.TimeFormat))
	ctx.Response.Header.Set("Server", config.ServerName)
}

// run Jssp server
func (s *JsspServer) Run() {
	err := fasthttp.ListenAndServe(":"+strconv.Itoa(config.Server.Port), s.ServeAll)
	if err != nil {
		println(err.Error())
	}
}
