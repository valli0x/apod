package config

import (
	"flag"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var ConfigFile = flag.String("config", "config/config.yml", "config file")

type Config struct {
	Apodkey  string `json:"apodkey"`
	Datapath string `json:"datapath"`

	Server struct {
		Address string `json:"address"`
	} `json:"server"`

	Logger struct {
		ServiceName string `json:"serviceName"`
	} `json:"logger"`

	Storage struct {
		DBname   string `json:"dbname"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"storage"`
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(*ConfigFile)

	v.AutomaticEnv()
	v.SetEnvPrefix("SERVICE")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		log.Error().Stack().Err(err).Msg("ReadInConfig")
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
