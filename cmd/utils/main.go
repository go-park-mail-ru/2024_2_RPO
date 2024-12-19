package main

// Программа CLI для обслуживания сервиса.
// Возможности:
// * Заполнить базу данных для нагрузочного тестирования;
// * Заменить индексы elasticsearch-а

import (
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/routines"
	"RPO_back/internal/pkg/utils/misc"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/olivere/elastic/v7"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("environment configuration is invalid: %v", err)
		return
	}

	postgresDB, err := misc.ConnectToPgx(config.CurrentConfig.Poll.PostgresPoolSize)
	if err != nil {
		log.Fatal("error connecting to PostgreSQL: ", err)
		return
	}
	defer postgresDB.Close()

	elasticClient, err := elastic.NewClient(elastic.SetURL("http://elastic:9200"), elastic.SetSniff(false))
	if err != nil {
		log.Error("error connecting to elasticsearch: " + err.Error())
		return
	}

	switch os.Args[1] {
	case "elastic-migrator":
		routines.ElasticMigrate(elasticClient, postgresDB)
	case "fill-db":
		routines.FillDB(postgresDB)
	}
}
