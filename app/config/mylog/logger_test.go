package mylog

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/broqiang/mdblog/app/config"
	"github.com/broqiang/mdblog/app/helper"
)

func TestLogNewLog(t *testing.T) {
	// 初始化配置文件
	cfg = manuallyGenerateConfiguration()

	// 准备写入日志的类型
	str := "Hello World " + time.Now().String()

	// 初始化普通日志，写入
	LogInfo.Println(str)

	// 初始化错误日志，写入
	LogErr.Println(str)
	LogErr.Close()

	// 拿出日志文件的完整路径及名称
	filePath := LogInfo.GetFileFullName()

	// 读取文件
	file, err := os.Open(filePath)
	helper.PanicErr(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

}

func manuallyGenerateConfiguration() config.Config {
	return config.Config{
		// 基本信息
		Name:  "BroQiang 博客",
		Host:  "",
		Port:  8080,
		Debug: false,

		// 日志
		Log: config.Log{
			Dir:    "logs",
			Mode:   "file",
			Access: true,
		},
	}
}
