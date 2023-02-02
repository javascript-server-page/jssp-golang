package filetype

import (
	"io/fs"
	"jssp/config"
	"os"
	"path/filepath"
	"strings"
)

type FileType uint8

const (
	FILE FileType = iota + 1
	DIR
	JSSP
	JSJS
)

func GetFileInfo(path string) (string, fs.FileInfo, FileType, error) {
	const JSSP_TYPE = ".jssp."
	const JSJS_TYPE = ".jsjs."
	const INDEX_JSSP = "index" + JSSP_TYPE + "html"
	const INDEX_JSJS = "index" + JSJS_TYPE + "html"
	if path[0] != '/' {
		path = "/" + path
	}
	if path[len(path)-1] == '/' {
		jssp_path := filepath.Join(config.Server.Dir, path, INDEX_JSSP)
		if stat, _ := os.Stat(jssp_path); stat != nil && !stat.IsDir() {
			return jssp_path, stat, JSSP, nil
		}
		jsjs_path := filepath.Join(config.Server.Dir, path, INDEX_JSJS)
		if stat, _ := os.Stat(jsjs_path); stat != nil && !stat.IsDir() {
			return jsjs_path, stat, JSJS, nil
		}
		path = path[:len(path)-1]
	}
	path = filepath.Join(config.Server.Dir, path)
	stat, err := os.Stat(path)
	if err != nil {
		return path, nil, 0, err
	}
	if stat.IsDir() {
		return path, stat, DIR, nil
	}
	if strings.Contains(filepath.Base(path), JSSP_TYPE) {
		return path, stat, JSSP, nil
	}
	if strings.Contains(filepath.Base(path), JSJS_TYPE) {
		return path, stat, JSJS, nil
	}

	return path, stat, FILE, nil
}

func GetFile(path string, ft FileType) (*os.File, error) {
	if (ft == DIR && config.Server.Enable.Dir) ||
		(ft == FILE && config.Server.Enable.File) ||
		(ft == JSSP && config.Server.Enable.Jssp) ||
		(ft == JSJS && config.Server.Enable.Jsjs) {
		return os.Open(path)
	}
	return nil, os.ErrPermission
}
