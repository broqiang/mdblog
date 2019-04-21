package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/broqiang/mdblog/app/helper"
)

// 用来测试是否正确加载配置文件，因为加载配置文件涉及到了路径，所以需要传入项目根目录的参数，
// 如： go test -root=../../
func TestNewConfig(t *testing.T) {
	// 手动初始化一个配置，需要从 config/config.toml 中获取，要保证一直
	cfg := manuallyGenerateConfiguration()

	tomlCfg := New()
	log.Println("toml: ", tomlCfg)

	if !reflect.DeepEqual(cfg, tomlCfg) {
		t.Errorf("want %v, have %v", cfg, tomlCfg)
	}
}

func TestLogNewLog(t *testing.T) {
	// 初始化配置文件
	cfg = manuallyGenerateConfiguration2()

	// 准备写入日志的类型
	str := "Hello World " + time.Now().String()

	// 初始化普通日志，写入
	Info := NewLogInfo()
	Info.Println(str)

	// 初始化错误日志，写入
	logErr := NewLogError()
	logErr.Println(str)
	logErr.Close()

	// 拿出日志文件的完整路径及名称
	filePath := Info.GetFileFullName()

	// 读取文件
	file, err := os.Open(filePath)
	helper.PanicErr(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

}

func manuallyGenerateConfiguration() Config {
	return Config{
		// 基本信息
		Name:  "BroQiang 博客",
		Host:  "",
		Port:  8080,
		Debug: true,

		// 日志
		Log: Log{
			Dir:    "logs",
			Mode:   "close",
			Access: true,
		},
	}
}

func manuallyGenerateConfiguration2() Config {
	return Config{
		// 基本信息
		Name:  "BroQiang 博客",
		Host:  "",
		Port:  8080,
		Debug: true,

		// 日志
		Log: Log{
			Dir:    "logs",
			Mode:   "close",
			Access: true,
		},
	}
}
