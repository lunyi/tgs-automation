package create_domain

import (
	"cdnetwork/internal/auth"
	"cdnetwork/pkg/cdnetworks/create_domain/models"
	"fmt"
	"strings"
)

func CreateDomains(createDomainNames string, originSet string) {
	domainNames := strings.Split(createDomainNames, ",")

	for _, createDomainName := range domainNames {
		err := addCdnDomain(createDomainName, originSet)
		if err != nil {
			fmt.Printf("Error adding CDN domain %s: %v\n", createDomainName, err)
		}
	}
}

func addCdnDomain(domainName string, originSet string) error {
	originConfig := &models.AddCdnDomainRequestOriginConfig{}
	originConfig.SetOriginIps(originSet)

	request := &models.AddCdnDomainRequest{}
	request.SetDomainName(domainName)
	fmt.Println(domainName)
	request.SetContractId("40016213")
	request.SetItemId("30")
	request.SetOriginConfig(originConfig)
	request.SetVersion("1.0.0")
	request.SetAccelerateNoChina("true")

	config := auth.BasicConfig{
		Uri:    "/cdnw/api/domain",
		Method: "POST",
	}

	response := auth.Invoke(config, request.String())
	fmt.Printf("Response body for domain %s: %#v\n", domainName, response)
	return nil
}
