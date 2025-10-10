package config

import "github.com/spf13/viper"

type Config struct {
	HTTPPort   int            `mapstructure:"http_port"`
	RPCPath    string         `mapstructure:"rpc_path"`
	JWTSecret  string         `mapstructure:"jwt_secret"`
	ClickHouse ClickHouseConf `mapstructure:"clickhouse"`
	S3         S3Conf         `mapstructure:"s3"`
	Fractal    FractalConf    `mapstructure:"fractal"`
}


type ClickHouseConf struct {
	Addr     string            `mapstructure:"addr"`
	Database string            `mapstructure:"database"`
	Auth     struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"auth"`
	Codec map[string]string `mapstructure:"codec"`
}

type S3Conf struct {
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Endpoint  string `mapstructure:"endpoint"`
}

type FractalConf struct {
	ShardBits       uint8   `mapstructure:"shard_bits"`
	CacheFraction   float64 `mapstructure:"cache_fraction"`
	ParquetPageSize int     `mapstructure:"parquet_page_size"`
}

func MustLoad() *Config {
	viper.SetConfigFile("fractal.toml") // exact file, no search
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SetDefault("http_port", 8899)
	viper.SetDefault("rpc_path", "/")

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}
	return &c
}
