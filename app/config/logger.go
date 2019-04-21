package config

import (
	"log"
	"os"
	"path"

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

	// ************************************************
	// 默认配置，当不能获取到相应的配置时，使用这里定义的默认配置
	// ************************************************

	// 日志的默认格式，当配置文件没有配置的时候使用这个
	defaultFormat = "2006-01/02"
)

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
		log.Println("========== log mode: ", cfg.Log.Mode)
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
	return path.Join(Root, cfg.Log.Dir, lw.fileName)
}

func (lw *LogWriter) newFile() *os.File {
	// logPath := path.Join(Root, lw.subDir, lw.fileName)
	dir := path.Join(Root, cfg.Log.Dir)
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
