package config

import "time"

type Config struct {
	Server *Server `yaml:"server"`
	Db     *Db     `yaml:"db"`
}

type Server struct {
	Host              string        `yaml:"host"`
	Port              int           `yaml:"port"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
}

type Db struct {
	ConnectionString string `yaml:"connectionString"`
}
