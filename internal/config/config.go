package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"
)

type Config struct {
	Env             string
	DBConn          DBConnConf
	Timeout         time.Duration
	IdleTimeout     time.Duration
	AliasLifeLength time.Duration
}

type DBConnConf struct {
	User     string
	Password string
	Host     string
	Port     int
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
