package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

type Config struct {
	Env         string        `json:"Env"`
	DBConn      DBConnConfig  `json:"DBConn" env-required:"true"`
	Timeout     time.Duration `json:"Timeout"`
	IdleTimeout time.Duration `json:"IdleTimeout"`
	Alias       AliasConfig   `json:"Alias"`
}

type AliasConfig struct {
	Length     int           `json:"Length" env-default:"8"`
	LifeLength time.Duration `json:"LifeLength" env-default:"2592000000000000"`
}

type DBConnConfig struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
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

	if err = json.Unmarshal(configData, &cfg); err != nil {
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

	flag.StringVar(&path, "config-path", "./path/to/conf.json", "path to config file")
	flag.Parse()
	return path
}
