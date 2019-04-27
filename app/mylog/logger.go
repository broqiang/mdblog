package mylog

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/broqiang/mdblog/app/config"
	"github.com/broqiang/mdblog/app/helper"
)

const (
	// *******
	// 日志模式
	// *******

	// LogClose 关闭日志
	LogClose string = "close"
	// LogFile 写入到日志文件
	LogFile string = "file"
	// LogStdout 写入到标准输出
	LogStdout string = "stdout"
)

var cfg = config.Cfg
var root = config.Root

// LogInfo 普通日志输出
var LogInfo = NewLogInfo()

// LogErr 错误日志输出
var LogErr = NewLogError()

// Logger 日志结构
type Logger struct {
	*LogWriter
	*log.Logger
}

// 有个问题， Info 和 Error 各自的锁是独立的，并发大的时候会产生问题
// 后面再解决，暂时先这样

// NewLogInfo 初始化普通日志
func NewLogInfo() *Logger {
	l := Logger{
		LogWriter: getSysLog(),
	}

	l.Logger = log.New(l.LogWriter, "[info] ", log.LstdFlags)

	return &l
}

// NewLogError 初始化错误日志
func NewLogError() *Logger {
	l := Logger{
		LogWriter: getSysLog(),
	}

	l.Logger = log.New(l.LogWriter, "[error] ", log.LstdFlags|log.Llongfile)

	return &l
}

var defaultLogWriter *LogWriter

// LogWriter 定义日志输出的结构
type LogWriter struct {
	// 日志文件的名称
	fileName string

	// 系统日志打开的文件
	*os.File
}

func getSysLog() *LogWriter {
	if defaultLogWriter == nil {
		defaultLogWriter = NewSysLog()
	}

	return defaultLogWriter
}

// NewAccessLog 初始化访问日志
func NewAccessLog() *LogWriter {
	lw := LogWriter{fileName: "access.log"}
	lw.setWriter()

	return &lw
}

// NewSysLog 初始化 LogWriter
func NewSysLog() *LogWriter {
	lw := LogWriter{fileName: "sys.log"}

	lw.setWriter()

	return &lw
}

func (lw *LogWriter) setWriter() {
	if cfg.Debug {
		lw.File = os.Stdout
		return
	}

	if cfg.Log.Mode == "file" {
		lw.File = lw.newFile()
		return
	}

	if cfg.Log.Mode == "close" {
		lw.File = nil
		return
	}

}

// GetFileFullName 获取日志文件的名称
func (lw *LogWriter) GetFileFullName() string {
	return path.Join(root, cfg.Log.Dir, lw.fileName)
}

func (lw *LogWriter) newFile() *os.File {
	dir := path.Join(root, cfg.Log.Dir)
	fileStat, err := os.Stat(dir)

	if err != nil {
		if !os.IsNotExist(err) {
			helper.Panicf("%v", err)
		}

		os.MkdirAll(dir, 0755)
	} else {
		if !fileStat.IsDir() {
			helper.Panicf("%q is an existing file", dir)
		}
	}

	logFile := path.Join(dir, lw.fileName)

	fileStat, err = os.Stat(logFile)

	if err != nil {
		if !os.IsNotExist(err) {
			helper.Panicf("faild to open %q file, %v", err)
		}
	}

	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	helper.PanicErr(err)

	return file
}

// Info 打印 info 信息
func Info(vals ...interface{}) {
	LogInfo.Output(2, fmt.Sprintln(vals...))
}

// Infof 指定格式的 info 日志
func Infof(format string, vals ...interface{}) {
	LogInfo.Output(2, fmt.Sprintf(format, vals...))
}

// Error 错误日志
func Error(vals ...interface{}) {
	LogErr.Output(2, fmt.Sprintln(vals...))
}

// Errorf 指定格式的错误日志
func Errorf(format string, vals ...interface{}) {
	LogErr.Output(2, fmt.Sprintf(format, vals...))
}
