package config

import (
	"log"
	"reflect"
	"testing"
)

// 用来测试是否正确加载配置文件，因为加载配置文件涉及到了路径，所以需要传入项目根目录的参数，
// 如： go test -root=../../
func TestNewConfig(t *testing.T) {
	// 手动初始化一个配置，需要从 config/config.toml 中获取，要保证一直
	cfg := manuallyGenerateConfiguration()

	log.Println("toml: ", Cfg)

	if !reflect.DeepEqual(cfg, Cfg) {
		t.Errorf("want %v, have %v", cfg, Cfg)
	}
}

func manuallyGenerateConfiguration() Config {
	return Config{
		// 基本信息
		Name:  "BroQiang 博客",
		URL:   "http://localhost:8080",
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
