package main

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/postgresql"
)

func getPromotionTypes(config util.TgsConfig) error {

	app1 := postgresql.NewPromotionTypesInterface(config.Postgresql)
	defer app1.Close()
	types, err := app1.GetPromotionTypes()
	if err != nil {
		log.LogFatal(err.Error())
	}

	log.LogInfo(fmt.Sprintf("Promotion types: %v", types))

	return nil

}
