package config

import (
	"flag"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Env    string       `toml:"env"`
	Server HTTPServer   `toml:"http_server"`
	DBConn DBConnConfig `toml:"db_conn" env-required:"true"`
	Alias  AliasConfig  `toml:"alias"`
}

type HTTPServer struct {
	Address     string        `toml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `toml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `toml:"idle_timeout" env-default:"60s"`
	User        string        `toml:"user" env-required:"true"`
	Password    string        `toml:"password" env-required:"true"`
}

type AliasConfig struct {
	Length     int           `toml:"length" env-default:"8"`
	LifeLength time.Duration `toml:"life_length" env-default:"2592000000000000"`
}

type DBConnConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DBName   string `toml:"name"`
}

// MustLoad parse config file in Config struct
// if there is no config file it panics
func MustLoad() *Config {
	path := fetchConfPath()
	if path == "" {
		panic("config path isn't set")
	}

	if _, err := os.Stat(path); err != nil {
		panic("incorrect config file path: " + path)
	}

	configData, err := os.ReadFile(path)
	if err != nil {
		panic("can't read config file: " + err.Error())
	}

	var cfg Config

	_, err = toml.Decode(string(configData), &cfg)
	if err != nil {
		panic("can't parse config file: " + err.Error())
	}
	return &cfg
}

// fetchConfPath finds config path in
// flag < env
func fetchConfPath() string {
	path := os.Getenv("CONFIG_PATH")
	if path != "" {
		return path
	}

	flag.StringVar(&path, "config-path", "./path/to/conf.toml", "path to config file")
	flag.Parse()
	return path
}
