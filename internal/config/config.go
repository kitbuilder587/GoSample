package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Env string `mapstructure:"env" validate:"required"`
	DB  struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     int    `mapstructure:"port" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
	} `mapstructure:"db" validate:"required"`
}

func MustReturnConfig() *Config {
	configureViper()
	var cfg *Config = &Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		panic("can't read config: " + err.Error())
	}

	mustValidateConfig(cfg)
	return cfg
}

func configureViper() {
	viper.SetConfigName("config") // name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config") // look in current directory
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic("can't read config: " + err.Error())
	}
}

func mustValidateConfig(cfg *Config) {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		panic("Config is not validated: " + err.Error())
	}
}
