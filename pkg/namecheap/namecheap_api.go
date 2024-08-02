package namecheap

import (
	"context"
	"tgs-automation/internal/util"
)

// mockgen -destination pkg/namecheap/mocks/mock_namecheap.go -package=namecheap tgs-automation/pkg/namecheap NamecheapAPI
// go generate mockgen -destination mocks/mock_namecheap.go -package=namecheap tgs-automation/pkg/namecheap NamecheapAPI
type NamecheapApi interface {
	CheckDomainAvailable(ctx context.Context, domain string) (bool, error)
	GetCouponCode(ctx context.Context) (string, error)
	GetDomainPrice(ctx context.Context, domain string) (*CheckDomainPriceResponse, error)
	CreateDomain(ctx context.Context, domainName string, promotionCode string) (string, error)
	GetBalance(ctx context.Context) (string, error)
	GetExpiredDomains() ([]FilteredDomain, error)
	UpdateNameServer(domain string, nameServers string) (*UpdateNameServerApiResponse, error)
}

type NamecheapService struct {
	Config util.TgsConfig
}

func New(config util.TgsConfig) NamecheapApi {
	return &NamecheapService{
		Config: config,
	}
}
