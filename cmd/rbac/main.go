package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	"github.com/omc-college/management-system/pkg/config"
	"github.com/omc-college/management-system/pkg/pubsub"
	"github.com/omc-college/management-system/pkg/rbac/api/routers"
	"github.com/omc-college/management-system/pkg/rbac/repository/postgres"
	"github.com/omc-college/management-system/pkg/rbac/service"
)

func main() {
	var serviceConfig Config
	var err error

	configPath := flag.StringP("config", "c", "./rbac-service-example-config.yaml", "path to service config")

	flag.Parse()

	err = config.Load(&serviceConfig, *configPath)
	if err != nil {
		logrus.Fatalf("%s", err)
	}
	// Open DB
	repository, err := postgres.NewRolesRepository(serviceConfig.RepositoryConfig)
	if err != nil {
		logrus.Fatalf("opening DB error")
	}

	defer repository.DB.Close()

	client, err := pubsub.NewQueueGroupClient(serviceConfig.PubSubConfig)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	rolesService := service.NewRolesService(repository, client)

	// Start server
	logrus.Fatal(http.ListenAndServe(":8000", routers.NewCrudRouter(rolesService)))
}
