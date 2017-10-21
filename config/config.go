package config

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Command string `yaml:"command"`
	Build   string `yaml:"build"`
}

func LoadConfig(config string) Config {
	var c Config
	data, err := ioutil.ReadFile(config)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(data, &c); err != nil {
		panic(err)
	}
	fmt.Printf("load config %+v\n", c)
	return c
}
