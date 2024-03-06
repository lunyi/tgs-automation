package disable_domain

import (
	"cdnetwork/internal/auth"
	"cdnetwork/pkg/cdnetworks/disable_domain/models"
	"fmt"
)

func DisableDomain(domainid string) {
	request := models.DisableSingleDomainServiceRequest{}
	fmt.Println("disable domain id is ", domainid)
	apiuri := fmt.Sprintf("/api/domain/%s/disable", domainid)
	fmt.Println(apiuri)
	config := auth.BasicConfig{
		Uri:    apiuri,
		Method: "PUT",
	}
	fmt.Printf(config.Uri)
	response := auth.Invoke(config, request.String())
	fmt.Printf("response body is %#v\n", response)
}
