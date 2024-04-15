package create_certificate

import (
	"fmt"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"time"

	"github.com/tebeka/selenium"
)

type CreateCDNCertServicer interface {
	CreateCertificate(domain string) error
	login()
	createCert(domain string)
}

type CreateCDNCertService struct {
	wd     selenium.WebDriver
	svc    *selenium.Service
	config util.TgsConfig
}

func CreateCertificateByDomain(domain string) {
	config := util.GetConfig()
	cdnService, err := NewCreateCDNCertService(config)
	if err != nil {
		fmt.Printf("Failed to initialize CDN service: %v\n", err)
		return
	}
	err = cdnService.CreateCertificate(domain)
	fmt.Printf("Failed to create certificate in CDN: %v\n", err)
}

func NewCreateCDNCertService(config util.TgsConfig) (*CreateCDNCertService, error) {
	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(config.ChromeDriverPath, 4444)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error starting the ChromeDriver server: %v", err.Error()))
		return nil, fmt.Errorf("error starting the ChromeDriver server: %v", err)
	}
	capabilities := selenium.Capabilities{
		"browserName": "chrome",
	}

	wd, err := selenium.NewRemote(capabilities, "http://localhost:4444/wd/hub")
	if err != nil {
		log.LogFatal(fmt.Sprintf("Failed to open session: %v\n", err))
		return nil, fmt.Errorf("failed to open session: %v", err)
	}
	return &CreateCDNCertService{
		wd:     wd,
		svc:    service,
		config: config,
	}, nil
}

func (cm *CreateCDNCertService) CreateCertificate(domain string) error {
	// defer cm.svc.Stop()
	// defer cm.wd.Quit()

	if err := cm.login(); err != nil {
		return err
	}

	if err := cm.createCertificate(domain); err != nil {
		return err
	}
	return nil
}

func (cm *CreateCDNCertService) login() error {
	if err := goToURL(cm.wd, cm.config.tgs-automation.CdnLoginUrl); err != nil {
		return err
	}

	if err := sendKeysToElement(cm.wd, selenium.ByID, "js-accountinput", cm.config.tgs-automation.CdnUserName); err != nil {
		return err
	}

	if err := sendKeysToElement(cm.wd, selenium.ByID, "js-pwdinput", cm.config.tgs-automation.CdnPassword); err != nil {
		return err
	}

	if err := clickElement(cm.wd, selenium.ByXPATH, "//button[@type='submit']"); err != nil {
		return err
	}

	time.Sleep(1 * time.Second)
	return nil
}

func (cm *CreateCDNCertService) createCertificate(domain string) error {
	if err := goToURL(cm.wd, cm.config.tgs-automation.CdnCertificateCreationUrl); err != nil {
		return err
	}

	if err := clickElement(cm.wd, selenium.ByCSSSelector, ".el-checkbox__input"); err != nil {
		return err
	}

	if err := sendKeysToElement(cm.wd, selenium.ByCSSSelector, "input.el-input__inner", domain); err != nil {
		return err
	}

	if err := clickElement(cm.wd, selenium.ByCSSSelector, "button.el-button.submit-btn.el-button--success"); err != nil {
		return err
	}

	return nil
}

func clickElement(wd selenium.WebDriver, by, selector string) error {
	ele, err := findElement(wd, by, selector)
	if err != nil {
		return err
	}
	if err := ele.Click(); err != nil {
		return err
	}
	return nil
}

func sendKeysToElement(wd selenium.WebDriver, by, selector, keys string) error {
	ele, err := findElement(wd, by, selector)
	if err != nil {
		return err
	}
	if err := ele.SendKeys(keys); err != nil {
		return err
	}
	return nil
}

func findElement(wd selenium.WebDriver, by string, selector string) (selenium.WebElement, error) {
	element, err := wd.FindElement(by, selector)
	if err != nil {
		log.LogFatal(fmt.Sprintf("failed to find element %s: %v", selector, err))
		return nil, err
	}
	return element, nil
}

func goToURL(wd selenium.WebDriver, url string) error {
	if err := wd.Get(url); err != nil {
		fmt.Printf("Failed to load page (%s): %v\n", url, err)
		log.LogFatal(fmt.Sprintf("Failed to load page (%s): %v\n", url, err))
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}
