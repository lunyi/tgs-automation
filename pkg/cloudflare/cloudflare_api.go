package cloudflare

import "tgs-automation/internal/util"

//mockgen -destination pkg/cloudflare/mock_cloudflare.go -package=cloudflare tgs-automation/pkg/Cloudflare CloudflareApi
//go generate mockgen -destination ../mocks/mock_cloudflare.go -package=cloudflare tgs-automation/pkg/Cloudflare CloudflareApi
type CloudflareApi interface {
	CreateDNS(domain string) error
	DeleteDNS(domain string) error
	GetDnsInfo(domain string) error
	GetTargetNameServers(domain string) (string, error)
}

type CloudflareService struct {
	Config util.TgsConfig
}

func NewClloudflare(config util.TgsConfig) CloudflareApi {
	return &CloudflareService{
		Config: config,
	}
}
