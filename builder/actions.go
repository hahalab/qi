package builder

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"encoding/base64"
	"github.com/hahalab/qi/archive"
	"github.com/hahalab/qi/config"
	"github.com/hahalab/qi/app"
	"github.com/hahalab/qi/aliyun/entity"
	"fmt"
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

	role, err := b.EnsureRole(app.Role())
	if err != nil {
		return
	}
	err = b.EnsureLogStore(app.ProjectName(), app.StoreName())
	if err != nil {
		return
	}
	service := app.Service()
	service.Role = role.Arn
	err = b.EnsureService(service)
	if err != nil {
		panic(err)
	}
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

func (b Builder) Configuration(app *app.App, m chan string) (groupAttribute entity.APIGroupAttribute, err error) {
	m <- "Configuration"

	role, err := b.EnsureRole(app.Role())
	if err != nil {
		return
	}
	groupAttribute, err = b.EnsureAPIGroup(app.APIGroup())
	if err != nil {
		return
	}

	for _, apiGateway := range app.APIGateways(groupAttribute.GroupId, role.Arn) {
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

	groupAttribute, err := b.Configuration(app, m)
	if err != nil {
		return err
	}

	m <- "Done"
	fmt.Printf("部署成功！请访问如下子域名:\n%s\n", groupAttribute.SubDomain)
	return nil
}
