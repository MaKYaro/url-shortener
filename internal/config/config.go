package config

import (
	"flag"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env    string       `yaml:"env"`
	Server HTTPServer   `yaml:"http_server"`
	DBConn DBConnConfig `yaml:"db_conn" env-required:"true"`
	Alias  AliasConfig  `yaml:"alias"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true"`
}

type AliasConfig struct {
	Length     int           `yaml:"length" env-default:"8"`
	LifeLength time.Duration `yaml:"lifeLength" env-default:"2592000000000000"`
}

type DBConnConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
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

	err = yaml.Unmarshal(configData, &cfg)
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

	flag.StringVar(&path, "config-path", "./path/to/conf.yml", "path to config file")
	flag.Parse()
	return path
}
