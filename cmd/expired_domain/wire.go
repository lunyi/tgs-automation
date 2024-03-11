//go:build wireinject
// +build wireinject

package main

import (
	"cdnetwork/internal/util"
	"cdnetwork/pkg/googlesheet"
	"cdnetwork/pkg/namecheap"
	"cdnetwork/pkg/postgresql"

	"github.com/google/wire"
)

func InitExpiredDomainsService(cfg util.TgsConfig) (*ExpiredDomainsService, error) {
	wire.Bind(expiredDomainSet, newExpiredDomainsService)
	return &ExpiredDomainsService{}, nil
}
