package config

import (
	"log"
	"os"
	"path"
	"time"

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

// NewLogInfo 初始化普通日志
func NewLogInfo() Logger {
	l := Logger{
		LogWriter: getSysLog(),
	}

	l.Logger = log.New(l.LogWriter, "[info] ", log.LstdFlags)

	return l
}

var defaultLogWriter *LogWriter

// LogWriter 定义日志输出的结构
type LogWriter struct {
	// 日志文件的名称
	fileName string

	// 系统日志打开的文件
	*os.File

	// 日志文件的格式
	format string

	// 当前的子目录，按照 mode 保存的格式
	subDir string
}

func getSysLog() *LogWriter {
	if defaultLogWriter == nil {
		defaultLogWriter = NewSysLog()
	}

	return defaultLogWriter
}

// NewSysLog 初始化 LogWriter
func NewSysLog() *LogWriter {
	lw := LogWriter{fileName: "sys.log"}

	lw.setDefaultConfig()
	lw.setWriter()

	log.Println(lw)

	return &lw
}

// AccessLogWriter 返回访问日志的 Writer， 暂时是配合 gin 使用
func AccessLogWriter() *os.File {
	al := NewAccessLog()

	return al.File
}

// NewAccessLog 初始化访问日志
func NewAccessLog() *LogWriter {
	lw := LogWriter{fileName: "access.log"}
	lw.setDefaultConfig()
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

	log.Println(cfg)
}

// GetFileName 获取日志文件的名称
func (lw *LogWriter) GetFileName() string {
	return path.Join(Root, cfg.Log.Dir, lw.subDir, lw.fileName)
}

func (lw *LogWriter) setDefaultConfig() {
	// 设置日志格式
	lw.format = defaultFormat
	if cfg.Log.Format != "" {
		lw.format = cfg.Log.Format
	}

	lw.subDir = time.Now().Format(lw.format)
}

func (lw *LogWriter) newFile() *os.File {
	// logPath := path.Join(Root, lw.subDir, lw.fileName)
	dir := path.Join(Root, cfg.Log.Dir, lw.subDir)
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
