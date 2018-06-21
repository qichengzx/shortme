package config

import (
	"github.com/Unknwon/goconfig"
	"log"
	"os"
)

var (
	configFile = ".env"
	Conf       = &goconfig.ConfigFile{}
	RUNMODE    string
	err        error
)

func init() {
	RUNMODE, err = GetEnv("RUNMODE", "dev")

	Conf, err = goconfig.LoadConfigFile(configFile)
	if err != nil {
		log.Fatal("Load Config Failed :%v\n", err)
	}
}

func Get(name string) (string, error) {
	return Conf.GetValue(RUNMODE, name)
}

func GetByBlock(block, name string) (string, error) {
	return Conf.GetValue(block, name)
}

func GetEnv(name, dft string) (string, error) {
	k := os.Getenv(name)
	if len(k) == 0 {
		k = dft
	}
	return k, nil
}
