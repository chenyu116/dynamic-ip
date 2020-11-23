package config

import (
	"github.com/spf13/viper"
	"io"
	"log"
	"time"
)

var c Config

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("/app")
	viper.AddConfigPath(".")
	ReadConfig()
}

func GetConfig() Config {
	return c
}

func GetCommonConfig() CommonConfig {
	return c.Common
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func IsSet(key string) bool {
	return viper.IsSet(key)
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

type CommonConfig struct {
	NodeName string `mapstructure:"nodeName"`
}

type ProviderConfig struct {
	URL           string        `mapstructure:"url"`
	CheckInterval time.Duration `mapstructure:"checkInterval"`
}

type ProvidersConfig struct {
	IPIP ProviderConfig `mapstructure:"ipip"`
}

type Config struct {
	Providers ProvidersConfig `mapstructure:"providers"`
	Common    CommonConfig    `mapstructure:"common"`
}

func ReadFromReader(reader io.Reader) error {
	err := viper.ReadConfig(reader)
	if err != nil {
		return err
	}
	c = Config{}

	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}
	return nil
}

func SetConfigPath(path string) {
	viper.SetConfigFile(path)
}

func ReadConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
	}
}
