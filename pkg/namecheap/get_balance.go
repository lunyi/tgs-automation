package namecheap

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (n *NamecheapService) GetBalance(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s&ApiUser=%s&ApiKey=%s&UserName=%s&Command=namecheap.users.getBalances&ClientIp=%s",
		n.Config.Namecheap.NamecheapBaseUrl,
		n.Config.Namecheap.NamecheapUsername,
		n.Config.Namecheap.NamecheapApiKey,
		n.Config.Namecheap.NamecheapUsername,
		n.Config.Namecheap.NamecheapClientIp)

	resp, err := http.Post(url, "application/xml", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
