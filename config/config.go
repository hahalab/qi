package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"

	"github.com/go-playground/validator"
	"github.com/sakeven/go-env"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var (
	cfg    = GwConf{}
	config = Config{}
)

type CommonConf struct {
	AliyunConfig      `json:"ali" env:"ALI" validate:"required,dive"`
	Debug      bool   `json:"debug" env:"DEBUG"`
	RouterPath string `json:"router_path" env:"ROUTER_PATH,router.conf" validate:"required"`
}

type GwConf struct {
	CommonConf `validate:"required,dive"`

	Port string `json:"port" env:"PORT,8080"`
}

type Config struct {
	CommonConf `validate:"required,dive"`

	CodePath string `validate:"required"`
	CodeConfig      `validate:"required,dive"`
}

type CodeConfig struct {
	Name       string   `yaml:"name" json:"name"`
	Command    string   `yaml:"command" json:"command"`
	Build      string   `yaml:"build" json:"build"`
	Files      []string `yaml:"files" json:"files"`
	MemorySize int      `yaml:"memorysize" json:"memorysize"`
	Timeout    int      `yaml:"timeout" json:"timeout"`
}

func (c CodeConfig) String() string {
	return fmt.Sprintf("Project name: %s  |  App init command: %s  |  Build command: %s", c.Name, c.Command, c.Build)
}

func LoadConfig(config string) (*CodeConfig, error) {
	var c CodeConfig
	data, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, &c); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", c)
	return &c, nil
}

func MustParseGWConfig(ctx *cli.Context) error {
	err := env.Decode(&(cfg.CommonConf))
	if err != nil {
		return err
	}

	if f := ctx.String(FlagPortPath); f != "" {
		cfg.Port = f
	}

	usr, err := user.Current()
	if err != nil {
		logrus.Fatal(err)
	}

	file, err := ioutil.ReadFile(usr.HomeDir + "/.ha.conf")
	if err == nil && len(file) != 0 {
		err = json.Unmarshal(file, &cfg)
		if err != nil {
			return err
		}
	}

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return validator.New().Struct(cfg)
}

func LoadQiConfig(path string, qiConfig *Config) (err error) {
	file, err := ioutil.ReadFile(path)
	if err == nil {
		err = json.Unmarshal(file, qiConfig)
	}
	return
}

func LoadCodeConfig(path string, codeConfig *CodeConfig) (err error) {
	file, err := ioutil.ReadFile(path)
	if err == nil {
		err = yaml.Unmarshal(file, codeConfig)
	}
	return
}

func MustParseConfig(ctx *cli.Context) error {
	config.CodePath = ctx.String(FlagCodePath)
	err := env.Decode(&(config.CommonConf))
	if err != nil {
		return err
	}

	usr, err := user.Current()
	if err != nil {
		logrus.Fatal(err)
	}

	err = LoadQiConfig(usr.HomeDir+"/.ha.conf", &config)
	if err != nil {
		return err
	}

	err = LoadCodeConfig(config.CodePath+"/qi.yml", &(config.CodeConfig))
	if err != nil {
		return err
	}

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return validator.New().Struct(config)
}

func GetGWConf() GwConf {
	return cfg
}

func GetConfig() Config {
	return config
}
