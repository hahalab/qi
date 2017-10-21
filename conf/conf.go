package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/go-playground/validator"
	"github.com/sakeven/go-env"
	"github.com/sirupsen/logrus"
	"github.com/todaychiji/ha/aliyun"
	"github.com/urfave/cli"
)

var (
	cfg   = GwConf{}
	upCfg = UpConf{}
)

type CommonConf struct {
	aliyun.Config `json:"ali" env:"ALI" validate:"required,dive"`
	Debug         bool   `json:"debug" env:"DEBUG"`
	RouterPath    string `json:"router_path" env:"ROUTER_PATH,router.conf" validate:"required"`
}

type GwConf struct {
	CommonConf `validate:"required,dive"`

	Port string `json:"port" env:"PORT,8080"`
}

type UpConf struct {
	CommonConf `validate:"required,dive"`

	CodePath   string `validate:"required"`
	CodeConfig `validate:"required,dive"`
}

type CodeConfig struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
	Build   string `yaml:"build"`
}

func (c CodeConfig) String() string {
	return fmt.Sprintf(`
	Project name: %s
	App init command: %s
	Build command: %s`, c.Name, c.Command, c.Build)
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
	logrus.Infof("%+v\n", c)
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

	file, err := ioutil.ReadFile("~/.ha.conf")
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

func MustParseUpConfig(ctx *cli.Context) error {
	upCfg.CodePath = ctx.String(FlagCodePath)
	err := env.Decode(&(upCfg.CommonConf))
	if err != nil {
		return err
	}

	file, err := ioutil.ReadFile("~/.ha.conf")
	if err == nil && len(file) != 0 {
		err = json.Unmarshal(file, &cfg)
		if err != nil {
			return err
		}
	}

	file, err = ioutil.ReadFile(upCfg.CodePath + "/ha.yml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &(upCfg.CodeConfig))
	if err != nil {
		return err
	}

	if upCfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return validator.New().Struct(upCfg)
}

func GetGWConf() GwConf {
	return cfg
}

func GetUPConf() UpConf {
	return upCfg
}
