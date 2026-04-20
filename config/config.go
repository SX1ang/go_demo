package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/spf13/viper"
)

type Config struct {
	*App   `mapstructure:"app"`
	*Log   `mapstructure:"log"`
	*Mysql `mapstructure:"mysql"`
	*Redis `mapstructure:"redis"`
}

type App struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Mode    string `mapstructure:"mode"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

type Log struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type Mysql struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func InitConfig() (*Config, error) {
	// 指定配置文件路径
	viper.SetConfigFile("./config/config.yaml")
	//viper.AddConfigPath(".")
	//viper.AddConfigPath("./config")
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")

	//读取配置信息
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	var Cfg Config
	if err := viper.Unmarshal(&Cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config file err: %s \n", err)
	}

	return &Cfg, nil
}

type EnvConfig struct {
	SecretKey string `env:"SECRET_KEY" envDefault:"9f4b2d7a6c3e8f1a2b4c5d6e7f8091a2b3c4d5e6f7081920a1b2c3d4e5f60718"`
}

func InitEnvConfig() (*EnvConfig, error) {
	var cfg EnvConfig
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func ProvideSecretKey(cfg *EnvConfig) string {
	return cfg.SecretKey
}
