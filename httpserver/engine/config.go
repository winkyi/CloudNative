package engine

import (
	"gopkg.in/ini.v1"
)

type Config interface {
	Load() (interface{}, error)
}

type IniConfig struct {
	FilePath string
}

func (c *IniConfig) Load() (interface{}, error) {
	cfg, err := ini.Load(c.FilePath)
	return cfg, err
}
