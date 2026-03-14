package config

import (
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

// MySQL 数据库配置
type MySQLConfig struct {
	Host         string `toml:"host"`
	Name         string `toml:"name"`
	Password     string `toml:"password"`
	Port         int    `toml:"port"`
	TablePrefix  string `toml:"tablePrefix"`
	User         string `toml:"user"`
}

// Config 应用配置
type Config struct {
	MySQL MySQLConfig `toml:"mysql"`
}

var (
	_cfg  *Config
	_once sync.Once
)

// configPaths 按顺序尝试的配置文件路径（兼容在子目录运行测试）
var configPaths = []string{"config.toml", "../config.toml", "../../config.toml"}

// GetConfig 返回全局配置，从 config.toml 加载（会从当前目录及上级目录查找）
func GetConfig() *Config {
	_once.Do(func() {
		var data []byte
		var err error
		for _, p := range configPaths {
			data, err = os.ReadFile(p)
			if err == nil {
				break
			}
		}
		if err != nil {
			panic("读取 config.toml 失败（已尝试 " + configPaths[0] + " 等）: " + err.Error())
		}
		var c Config
		if err := toml.Unmarshal(data, &c); err != nil {
			panic("解析 config.toml 失败: " + err.Error())
		}
		_cfg = &c
	})
	return _cfg
}
