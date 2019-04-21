package config

import (
	"log"
	"reflect"
	"testing"
	"time"
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
	str := "Hello World" + time.Now().String()

	// 初始化普通日志，写入
	Info := NewLogInfo()
	Info.Println(str)
	Info.Close()

	// 初始化日志，然后写入一行字符串
	// l := NewSysLog()
	// fileName := l.GetFileName()
	// l.WriteString(str)
	// l.Close()

	// // 读取日志文件的最后一行内容，然后匹配
	// file, err := os.Open(fileName)
	// helper.PanicErr(err)

	// scanner := bufio.NewScanner(file)

	// var haveStr string
	// for scanner.Scan() {
	// 	haveStr = scanner.Text()
	// }

	// // 将末尾的换行去掉，因为读取文件的时候是按照换行分割的文件
	// str = strings.TrimSpace(str)
	// if haveStr != str {
	// 	t.Errorf("want %q, have %q", str, haveStr)
	// }

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
		Debug: false,

		// 日志
		Log: Log{
			Dir:    "logs",
			Mode:   "file",
			Access: true,
		},
	}
}
