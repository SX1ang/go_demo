package config

import (
	"fmt"

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
