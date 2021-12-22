package engine

import (
	"math/rand"
	"time"
)

func RandInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func InitConfig(configfile string) interface{} {
	iniConf := IniConfig{FilePath: configfile}
	config, err := iniConf.Load()
	if err != nil {
		panic("can not load config")
	}
	return config
}
