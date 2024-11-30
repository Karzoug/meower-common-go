package prom

import (
	"fmt"
	"time"
)

type ServerConfig struct {
	Host         string        `env:"HOST"`
	Port         uint          `env:"PORT,notEmpty" envDefault:"3002"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT"  envDefault:"5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"5s"`
}

func (cfg ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}
