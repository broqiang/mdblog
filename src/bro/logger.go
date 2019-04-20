package bro

import (
	"fmt"
	"io"
	"log"
	"os"
)

// DefaultAccessWriter 默认的日志输出
var DefaultAccessWriter io.Writer = os.Stdout

// DefaultSystemWriter 默认的系统日志输出
var DefaultSystemWriter io.Writer = os.Stdout

var sl = log.New(DefaultSystemWriter, "", log.LstdFlags|log.Llongfile)

// SystemLogError 系统错误日志
func SystemLogError(err error) {
	if err != nil {
		SystemLogf("[Error] %v", err)
	}
}

// SystemLogPanic 记录 Panic 时候的错误日志
func SystemLogPanic(err interface{}) {
	SystemLogf("[Panic] %v", err)
}

// SetSystemLogWriter 定义日志输出到哪个 Writer
func SetSystemLogWriter(writer io.Writer) {
	sl.SetOutput(writer)
}

// SystemLogRoute 注册路由的日志
func SystemLogRoute(httpMethod, absolutePath string, handlers Handlers) {
	handlerNumbers := len(handlers)
	handlerName := NameOfFunction(handlers.Last())

	sl.Output(2,
		fmt.Sprintf(" register router:\n\tMethod: %s, \n\tPath: %s \n\tHandlerName: %s \n\tTotal handlers:  %d\n\n",
			httpMethod, absolutePath, handlerName, handlerNumbers))
}

// SystemLogln 写入系统日志（非访问日志），每一条一行
func SystemLogln(values ...interface{}) {
	sl.Output(2, fmt.Sprintln(values...))
}

// SystemLogf 写入指定格式的系统日志
func SystemLogf(format string, values ...interface{}) {
	sl.Output(2, fmt.Sprintf(format, values...))
}
