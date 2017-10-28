package builder

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"encoding/base64"
	"github.com/hahalab/qi/archive"
	"github.com/hahalab/qi/config"
	"github.com/hahalab/qi/app"
)

func (b Builder) Build(path string, m chan string) (err error) {
	err = archive.Build(path, m)
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}

func (b Builder) Prepare(app *app.App, m chan string) (err error) {
	m <- "Preparing"

	err = b.EnsureLogStore(app.ProjectName(), app.StoreName())
	if err != nil {
		return
	}

	err = b.EnsureService(app.Service())

	return
}

func (b Builder) Deploy(app *app.App, m chan string) error {
	m <- "Deploying"

	file, err := ioutil.ReadFile("code.zip")
	if err != nil {
		return err
	}

	function := app.Function()

	function.Code.ZipFile = base64.StdEncoding.EncodeToString(file)

	err = b.DeployFunction(app.Service().ServiceName, function)
	if err != nil {
		return err
	}

	return nil
}

func (b Builder) Configuration(app *app.App, m chan string) (err error) {
	m <- "Configuration"
	groupAttribute, err := b.EnsureAPIGroup(app.APIGroup())
	if err != nil {
		return
	}

	for _, apiGateway := range app.APIGateways(groupAttribute.GroupId) {
		err = b.EnsureAPIGateway(groupAttribute.GroupId, apiGateway)
		if err != nil {
			return
		}
	}

	return
}

func (b Builder) Qi(m chan string) error {
	c := config.GetConfig()
	app := app.NewApp(&c)

	if err := b.Build(c.CodePath, m); err != nil {
		return err
	}

	if err := b.Prepare(app, m); err != nil {
		return err
	}

	if err := b.Deploy(app, m); err != nil {
		return err
	}

	if err := b.Configuration(app, m); err != nil {
		return err
	}
	return nil
}
