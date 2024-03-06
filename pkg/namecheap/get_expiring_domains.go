package namecheap

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"fmt"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

//TODO: to fix the element not interactable

func Login(config util.TgsConfig) {
	selenium.SetDebug(false)
	service, err := selenium.NewChromeDriverService(config.ChromeDriverPath, 4444)
	if err != nil {
		log.LogFatal(fmt.Sprintf("Error starting the ChromeDriver server: %v", err.Error()))
		service.Stop()
	}
	caps := selenium.Capabilities{"browserName": "chrome"}

	chromeCaps := chrome.Capabilities{
		Path: "",
		Args: []string{
			"--disable-gpu",
			//"--no-sandbox",
			//"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			//"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}
	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, "http://localhost:4444/wd/hub")
	if err != nil {
		fmt.Println("Failed to open session:", err)
		return
	}

	time.Sleep(5 * time.Second)
	wd.SwitchWindow("chrome")
	// Maximize the browser window
	if err := wd.MaximizeWindow(""); err != nil {
		fmt.Printf("Failed to maximize window: %v\n", err)
		return
	}

	//defer wd.Quit()

	if err := wd.Get("https://www.namecheap.com/myaccount/signout"); err != nil {
		fmt.Println("Failed to load page:", err)
		return
	}

	fmt.Println("test4")

	usernameField, err := wd.FindElement(selenium.ByName, "LoginUserName")
	if err != nil {
		fmt.Println("Failed to find username field:", err)
		return
	}

	size, err := usernameField.Size()

	fmt.Println("size:", size.Width, size.Height)
	time.Sleep(5 * time.Second)

	if err := usernameField.Click(); err != nil {
		fmt.Println("Failed to click username field:", err)
		return
	}

	fmt.Println("test6")

	if err := usernameField.SendKeys(config.Namecheap.NamecheapUsername); err != nil {
		fmt.Println("Failed to enter username:", err)
		return
	}

	fmt.Println("test7")

	passwordField, err := wd.FindElement(selenium.ByName, "LoginPassword")
	if err != nil {
		fmt.Println("Failed to find password field:", err)
		return
	}

	usernameField.IsDisplayed()
	if err := usernameField.Click(); err != nil {
		fmt.Println("Failed to click username field:", err)
		return
	}

	if err := passwordField.SendKeys(config.Namecheap.NamecheapPassword); err != nil {
		fmt.Println("Failed to enter password:", err)
		return
	}

	// Find the login button and click it
	loginSubmitButton, err := wd.FindElement(selenium.ByName, "ctl00$ctl00$ctl00$ctl00$base_content$web_base_content$home_content$page_content_left$ctl02$LoginButton")
	if err != nil {
		fmt.Println("Failed to find login submit button:", err)
		return
	}
	if err := loginSubmitButton.Click(); err != nil {
		fmt.Println("Failed to click login submit button:", err)
		return
	}

	fmt.Println("Successfully logged in to Namecheap.")
}

// func waitForElement(wd selenium.WebDriver, by, value string) (selenium.WebElement, error) {
// 	const timeout = 10 * time.Second
// 	const interval = 100 * time.Millisecond

// 	w, err := selenium.Wait(wd, timeout, interval)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return w.Until(wd.FindElement(by, value))
// }
