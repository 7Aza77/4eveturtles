package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Postgres `yaml:"postgres"`
	Redis    `yaml:"redis"`
	Auth     `yaml:"auth"`
}

type HTTPServer struct {
	Address string `yaml:"address" env:"HTTP_ADDRESS" env-default:":8080"`
}

type Postgres struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User     string `yaml:"user" env:"DB_USER" env-default:"user"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"password"`
	DBName   string `yaml:"db_name" env:"DB_NAME" env-default:"goevent"`
}

type Redis struct {
	Host string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
}

type Auth struct {
	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
}

var (
	cfg  *Config
	once sync.Once
)

func MustLoad() *Config {
	once.Do(func() {
		cfg = &Config{}
		if err := cleanenv.ReadConfig("config/config.yaml", cfg); err != nil {
			// if config file not found, try to read from env
			if err := cleanenv.ReadEnv(cfg); err != nil {
				log.Fatalf("cannot read config: %s", err)
			}
		}
	})
	return cfg
}
