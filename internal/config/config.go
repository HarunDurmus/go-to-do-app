package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Server struct {
	Port string
}

type Config struct {
	Appname   string
	Server    Server
	LogLevel  string
	Couchbase Couchbase
}
type Couchbase struct {
	URL      string
	Username string
	Password string
	Buckets  []BucketConfig
}
type BucketConfig struct {
	Name               string
	CreatePrimaryIndex bool
	Scopes             []ScopeConfig
}
type ScopeConfig struct {
	Name        string
	Collections []CollectionConfig
}

type CollectionConfig struct {
	Name               string
	CreatePrimaryIndex bool
	FieldIndexes       []string
}

func New(configPath, configName string) (Config, error) {
	var config Config
	viperConfig, err := readConfig(configPath, configName)
	if err != nil {
		return config, err
	}

	if err := viperConfig.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}

func readConfig(configPath, configName string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	err := v.ReadInConfig()

	return v, err
}

func (c *Config) Print() {
	fmt.Println(c)
}
