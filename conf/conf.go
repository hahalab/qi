package conf

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-playground/validator"
	"github.com/sakeven/go-env"
	"github.com/sirupsen/logrus"
	"github.com/todaychiji/ha/aliyun"
	"github.com/urfave/cli"
)

var (
	cfg = Conf{}
)

type Conf struct {
	aliyun.Config `json:"ali" env:"ALI" validate:"required,dive"`

	Debug      bool   `json:"debug" env:"DEBUG"`
	RouterPath string `json:"router_path" env:"ROUTER_PATH,router.conf" validate:"required"`
	Port       string `json:"port" env:"PORT,8080"`
}

func MustParseConfig(ctx *cli.Context) error {
	err := env.Decode(&cfg)
	if err != nil {
		return err
	}

	if f := ctx.String(FlagConfigFile); f != "" {
		file, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

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

func GetConf() Conf {
	return cfg
}
