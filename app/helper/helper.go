// Package helper 是助手函数
package helper

import (
	"fmt"
	"runtime"
)

// Panicf 用来打印指定格式的 panic 信息，
// 需要注意，只能直接调用，间接调用文件的位置会定位错误（层级不准确）
func Panicf(format string, values ...interface{}) {
	if _, file, line, ok := runtime.Caller(1); ok {
		format += "\n\t in %d line of the %q file"
		values = append(values, line, file)
	}

	panic(fmt.Sprintf(format, values...))
}

// PanicErr 用来断言错误，如果存在错误就 panic
func PanicErr(err error) {
	if err != nil {
		var format string
		var values []interface{}
		if _, file, line, ok := runtime.Caller(1); ok {
			format = "\n\t in %d line of the %q file"
			values = append(values, line, file, err)
		}
		panic(fmt.Sprintf(format, values...))
	}
}
