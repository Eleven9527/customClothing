package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var (
	Mode string //开发环境：dev prd
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

type EnvConfig struct {
}

type DbConfig struct {
	Dsn     string `toml:"dsn"`
	MaxIdle int    `toml:"maxIdle"`
	MaxOpen int    `toml:"mxOpen"`
	Level   string `toml:"level"`
}

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DBIndex  int    `toml:"dbindex"`
	PoolSize int    `toml:"poolSize"`
	IdleSize int    `toml:"idleConns"`
}

type TokenConfig struct {
	Timeout        int64  `toml:"timeout"`
	CacheKey       string `toml:"cacheKey"`
	TokenDelimiter string `toml:"tokenDelimiter"`
	HeaderKey      string `toml:"headerKey"`
	EncryptKey     string `toml:"encryptKey"`
	MultiLogin     bool   `toml:"multiLogin"`
}

type Config struct {
	DbCfg    *DbConfig    `toml:"db"`
	RedisCfg *RedisConfig `toml:"redis"`
	TokenCfg *TokenConfig `toml:"token"`
	EnvCfg   *EnvConfig   `toml:"env"`
}

func Cfg() *Config {
	if cfg != nil {
		return cfg
	}

	InitConfig()
	return cfg
}

func InitConfig() {
	cPath := ""

	switch Mode {
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
