//go:build wireinject
// +build wireinject

package main

import (
	"tgs-automation/internal/util"

	"github.com/google/wire"
)

func InitExpiredDomainsService(cfg util.TgsConfig) (*ExpiredDomainsService, error) {
	wire.Build(
		newExpiredDomainsService,
	)
	return &ExpiredDomainsService{}, nil
}
