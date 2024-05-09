package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	FleetSize     int     `envconfig:"FLEET_SIZE" default:"12"`
	OriginX       float64 `envconfig:"ORIGIN_X" default:"0"`
	OriginY       float64 `envconfig:"ORIGIN_Y" default:"0"`
	HoursPerShift int     `envconfig:"HOURS_PER_SHIFT" default:"12"`
	LogLevel      string  `envconfig:"LOG_LEVEL" default:"Debug"`
}

func NewConfig() *Config {
	result := &Config{}
	if err := envconfig.Process("", result); err != nil {
		panic(err)
	}
	return result
}
