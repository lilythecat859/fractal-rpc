package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort   int    `mapstructure:"http_port"`
	JWTSecret  string `mapstructure:"jwt_secret"`
	ClickHouse struct {
		Addr     string `mapstructure:"addr"`
		Database string `mapstructure:"database"`
		Auth     struct {
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		} `mapstructure:"auth"`
	} `mapstructure:"clickhouse"`
	S3 struct {
		Bucket    string `mapstructure:"bucket"`
		Region    string `mapstructure:"region"`
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Endpoint  string `mapstructure:"endpoint"`
	} `mapstructure:"s3"`
}

func MustLoad() *Config {
	viper.SetConfigName("fractal")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
	var c Config
	_ = viper.Unmarshal(&c)
	return &c
}
