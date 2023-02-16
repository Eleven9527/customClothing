package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

const (
	//dev环境
	DevFlag    = "dev"
	DevCfgPath = "./config/config_dev.toml"

	//prd环境
	PrdFlag    = "prd"
	PrdCfgPath = "./config/config_prd.toml"
)

var (
	cfg *Config //config实例
)

type DbConfig struct {
	Dsn     string `toml:"dsn"`
	MaxIdle int    `toml:"maxIdle"`
	MaxOpen int    `toml:"mxOpen"`
	Level   string `toml:"level"`
}

type Config struct {
	DbCfg *DbConfig `toml:"db"`
}

func Cfg() *Config {
	return cfg
}

func InitConfig(env string) {
	cPath := ""

	switch env {
	case DevFlag:
		cPath = DevCfgPath
	case PrdFlag:
		cPath = PrdCfgPath
	default:
		panic("配置文件路径错误!")
	}

	if _, err := toml.DecodeFile(cPath, &cfg); err != nil {
		fmt.Println("解析配置异常:", err.Error())
		os.Exit(1)
	}
}
