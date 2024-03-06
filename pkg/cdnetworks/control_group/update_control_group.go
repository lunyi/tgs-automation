package control_group

import (
	"cdnetwork/internal/auth" // Import your auth package
	"cdnetwork/pkg/cdnetworks/control_group/models"
	"fmt"
	"strings"
)

// ControlGroupManager represents a manager for control groups
type ControlGroupManager struct {
	ControlGroupName string
	LoginNames       []*models.LoginName
	AuthConfig       auth.BasicConfig
}

// NewControlGroupManager creates a new instance of ControlGroupManager
func NewControlGroupManager(controlGroupName, uri, method string) *ControlGroupManager {
	config := auth.BasicConfig{
		Uri:    uri,
		Method: method,
	}

	return &ControlGroupManager{
		ControlGroupName: controlGroupName,
		AuthConfig:       config,
	}
}

// SetLoginNameList sets the login name list
func (c *ControlGroupManager) setLoginNameList(loginNames []*models.LoginName) {
	c.LoginNames = loginNames
}

// UpdateControlGroup updates the control group for a given domain
func (c *ControlGroupManager) updateControlGroup(domain string) {
	request := models.EditControlGroupRequest{}
	request.SetDomainList([]*string{&domain})
	request.SetAccountList(c.LoginNames)
	request.SetIsAdd(true)
	request.SetControlGroupName(c.ControlGroupName)

	response := auth.Invoke(c.AuthConfig, request.String())
	fmt.Printf("Response body for domain %s is %#v\n", domain, response)
}

func UpdateControlGroup(domainName string) {
	// 使用逗号分割多个域名
	domainNames := strings.Split(domainName, ",")

	loginNameList := []*models.LoginName{}
	loginName1 := models.LoginName{}
	loginName1.SetLoginName("raypro")
	loginNameList = append(loginNameList, &loginName1)

	controlGroupName := "Raypro"

	manager := NewControlGroupManager(
		controlGroupName,
		"/user/control-groups/fdb99358-9e43-4197-95c7-7c2213a71a6f",
		"PUT",
	)
	manager.setLoginNameList(loginNameList)

	for _, createDomainName := range domainNames {
		manager.updateControlGroup(createDomainName)
	}
}
