package delete_domain

import (
	"cdnetwork/internal/auth"
	"cdnetwork/pkg/cdnetworks/delete_domain/models"
	"cdnetwork/pkg/cdnetworks/get_domain_id"
	"encoding/xml"
	"fmt"
	"strconv"
)

type Response struct {
	XMLName xml.Name `xml:"response"`
	Code    string   `xml:"code"`
	Message string   `xml:"message"`
}

func DeleteDomain(domainName string) {

	fmt.Printf("call delete domain %#v\n", domainName)
	_domainId := get_domain_id.GetDomainId(domainName)

	fmt.Printf("response delete domain is %#v\n", _domainId)

	domainId := strconv.Itoa(_domainId)
	request := &models.DeleteApiDomainServiceRequest{}
	request.SetDomainName(domainName)
	request.SetDomainId(domainId)

	apiuri := fmt.Sprintf("/api/domain/%s", domainId)
	config := auth.BasicConfig{
		Uri:    apiuri,
		Method: "DELETE"}
	xml_response := auth.Invoke(config, request.String())

	var response Response
	err := xml.Unmarshal([]byte(xml_response), &response)
	if err != nil {
		fmt.Println("Error parsing XML:", err)
		return
	}

	fmt.Printf("response: %#v\n", response.Message)
}
