package bro

import (
	"net/http"
	"path"
	"reflect"
	"runtime"
)

// JoinPath 合并路径
func JoinPath(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)

	if lastChar(finalPath) != '/' && lastChar(relativePath) == '/' {
		finalPath += "/"
	}

	return finalPath

}

// 获取字符串的最后一个字符
func lastChar(str string) uint8 {
	if str == "" {
		return 0
	}

	return str[len(str)-1]
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// PrepareRedirectCode 获取重定向使用的状态码
func PrepareRedirectCode(method string) int {

	if method == "GET" {
		return http.StatusTemporaryRedirect
	}

	return http.StatusMovedPermanently
}
