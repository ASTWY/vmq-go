package config

import (
	"fmt"
	"os"

	"vmq-go/utils"

	"gopkg.in/yaml.v3"
)

type logConfig struct {
	Level string `yaml:"level"` // 日志级别
	Path  string `yaml:"path"`  // 日志文件目录
}

type jwtConfig struct {
	Secret string `yaml:"secret"` // 密钥
	Expire int    `yaml:"expire"` // 以秒为单位
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`     // 数据库地址
	Port     int    `yaml:"port"`     // 数据库端口
	User     string `yaml:"user"`     // 数据库用户名
	Password string `yaml:"password"` // 数据库密码
	DBName   string `yaml:"name"`     // 数据库名
}

type Config struct {
	Host     string         `yaml:"host"` // 服务监听地址
	Port     int            `yaml:"port"` // 服务监听端口
	Log      logConfig      `yaml:"log"`  // 日志配置
	Jwt      jwtConfig      `yaml:"jwt"`  // JWT 配置
	Database DatabaseConfig `yaml:"db"`   // 数据库配置
}

var (
	Conf *Config
)

func loadConfigFile(fielName string) (err error) {
	f, err := os.Open(fielName)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(Conf)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	Conf = &Config{
		Host: "0.0.0.0",
		Port: 8080,
		Log: logConfig{
			Level: "info",
			Path:  utils.GetCWDir() + "/logs",
		},
		Jwt: jwtConfig{
			Secret: "hfjtSdYL6j!Ts$uD",
			Expire: 3600, // 1 小时
		},
		Database: DatabaseConfig{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "root",
			Password: "123456",
			DBName:   "vpay",
		},
	}
	err := loadConfigFile(utils.GetCWDir() + "/config.yaml")
	if err != nil {
		// 配置文件加载失败，写一个默认的配置文件
		configStr := `
host: 0.0.0.0
port: 8080
log:
  level: info
  path: ./logs
jwt:
  secret: hfjtSdYL6j!Ts$uD
  expire: 3600
db:
  host: localhost
  port: 3306
  user: test
  password: test
  name: test
`
		err := os.WriteFile(utils.GetCWDir()+"/config.yaml", []byte(configStr), 0644)
		if err != nil {
			panic("Write config file failed!")
		}
		panic(fmt.Sprintf("Load config file failed: %s,writing a default config file to %s", err.Error(), utils.GetCWDir()+"/config.yaml"))
	}
}
