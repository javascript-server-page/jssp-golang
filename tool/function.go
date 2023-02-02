package tool

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// GetResource
func GetResource(path string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		res, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(res.Body)
	}
	return ioutil.ReadFile(path)
}

// FileExists 判断文件是否存在,并且不能是文件夹
// func FileExists(path string) bool {
// 	if stat, err := os.Stat(path); err == nil {
// 		return !stat.IsDir()
// 	}
// 	return false
// }

// OpenFileAndInfo
// func OpenFileAndInfo(fs http.FileSystem, name string) (http.File, fs.FileInfo, error) {
// 	f, err := fs.Open(name)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	stat, err := f.Stat()
// 	if err != nil {
// 		f.Close()
// 		return nil, nil, err
// 	}
// 	return f, stat, nil
// }
