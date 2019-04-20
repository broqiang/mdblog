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

	if LastChar(finalPath) != '/' && LastChar(relativePath) == '/' {
		finalPath += "/"
	}

	return finalPath

}

// LastChar 获取字符串的最后一个字符
func LastChar(str string) uint8 {
	if str == "" {
		return 0
	}

	return str[len(str)-1]
}

// NameOfFunction 获取函数的名称
func NameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// PrepareRedirectCode 获取重定向使用的状态码
func PrepareRedirectCode(method string) int {

	if method == "GET" {
		return http.StatusTemporaryRedirect
	}

	return http.StatusMovedPermanently
}
