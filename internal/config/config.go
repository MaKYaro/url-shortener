package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env    string       `env:"Env"`
	Server HTTPServer   `env:"HTTPServer"`
	DBConn DBConnConfig `env:"DBConn" env-required:"true"`
	Alias  AliasConfig  `env:"Alias"`
}

type HTTPServer struct {
	Address     string        `env:"Address" env-default:"localhost:8080"`
	Timeout     time.Duration `env:"Timeout" env-default:"4s"`
	IdleTimeout time.Duration `env:"IdleTimeout" env-default:"60s"`
	User        string        `env:"User" env-required:"true"`
	Password    string        `env:"Password" env-required:"true"`
}

type AliasConfig struct {
	Length     int           `env:"Length" env-default:"8"`
	LifeLength time.Duration `env:"LifeLength" env-default:"2592000000000000"`
}

type DBConnConfig struct {
	User     string `env:"User"`
	Password string `env:"Password"`
	Host     string `env:"Host"`
	Port     int    `env:"Port"`
	DBName   string `env:"DBName"`
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

	err := godotenv.Load(path)
	if err != nil {
		panic("can't load env variables: " + err.Error())
	}

	cfg, err := fillConfig()
	if err != nil {
		panic("can't parse env variable: " + err.Error())
	}

	return cfg
}

func fillConfig() (*Config, error) {
	const op = "internal.config.fillConfig"

	var cfg *Config = &Config{}
	var err error

	cfg.Env = os.Getenv("ENV")

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")
	cfg.Server.User = os.Getenv("SERVER_USER")
	cfg.Server.User = os.Getenv("SERVER_PASSWORD")
	cfg.Server.Timeout, err = time.ParseDuration(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cfg.Server.IdleTimeout, err = time.ParseDuration(os.Getenv("SERVER_IDLE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	cfg.DBConn.User = os.Getenv("DB_USER")
	cfg.DBConn.Password = os.Getenv("DB_PASSWORD")
	cfg.DBConn.Host = os.Getenv("DB_HOST")
	cfg.DBConn.DBName = os.Getenv("DB_NAME")
	cfg.DBConn.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	cfg.Alias.Length, err = strconv.Atoi(os.Getenv("ALIAS_LENGTH"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cfg.Alias.LifeLength, err = time.ParseDuration(os.Getenv("ALIAS_LIFE_LENGTH"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return cfg, nil
}

// fetchConfPath finds config path in
// flag < env
func fetchConfPath() string {
	path := os.Getenv("CONFIG_PATH")
	if path != "" {
		return path
	}

	flag.StringVar(&path, "config-path", "./path/to/conf.env", "path to config file")
	flag.Parse()
	return path
}
