package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port         int    `yaml:"port"`
	JwtSecretKey string `yaml:"jwt_secret_key"`
	Components   struct {
		Database struct {
			Name               string `yaml:"name"`
			Username           string `yaml:"username"`
			Password           string `yaml:"password"`
			Host               string `yaml:"host"`
			Port               string `yaml:"port"`
			ConnectionsLimit   int    `yaml:"connectionslimit"`
			ConnectionTimeout  string `yaml:"connectiontimeout"`
			ConnectionLifetime string `yaml:"connectionlifetime"`
		} `yaml:"database"`
	} `yaml:"components"`
}

func GetConfig(path string, validator *validator.Validate) *Config {
	const f = "NewConfig "

	config := &Config{}
	err := cleanenv.ReadConfig(path, config)
	if err != nil {
		panic(f + err.Error())
	}

	err = validator.Struct(config)
	if err != nil {
		panic(f + err.Error())
	}

	return config
}
