package config

import (
	"net/http"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// Config Service config
type Config struct {
	Grpc Grpc `json:"grpc" yaml:"grpc"`
	Http Http `json:"http" yaml:"http"`
	Pg   Pg   `json:"pg" yaml:"pg"`
}

// NewConfig Initial service's config
func NewConfig(cfg string) *Config {

	if cfg == "" {
		panic("load config file failed.config file can not be empty.")
	}

	viper.SetConfigFile(cfg)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		panic("read config failed.[ERROR]=>" + err.Error())
	}
	conf := &Config{}
	// Assign the overloaded configuration to the global
	if err := viper.Unmarshal(conf); err != nil {
		panic("assign config failed.[ERROR]=>" + err.Error())
	}

	return conf

}

// Grpc Grpc server config
type Grpc struct {
	Host   string `json:"host" yaml:"host"`
	Port   string `json:"port" yaml:"port"`
	Name   string `json:"name" yaml:"name"`
	Server *grpc.Server
}

// Http Http server config
type Http struct {
	Host   string `json:"host" yaml:"host"`
	Port   string `json:"port" yaml:"port"`
	Name   string `json:"name" yaml:"name"`
	Server *http.Server
}

type Pg struct {
	Host                         string `json:"host"                    yaml:"host"`
	Username                     string `json:"username"                 yaml:"username"`
	Password                     string `json:"-"                                  yaml:"password"`
	Database                     string `json:"database"                           yaml:"database"`
	MaxIdleConnections           int    `json:"max-idle-connections"     yaml:"max-idle-connections"`
	MaxOpenConnections           int    `json:"max-open-connections"     yaml:"max-open-connections"`
	MaxConnectionLifeTimeSeconds int    `json:"max-connection-life-time" yaml:"max-connection-life-time-seconds"`
	LogLevel                     int    `json:"log-level"                          yaml:"log-level"`
}
