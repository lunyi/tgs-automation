package namecheap

import (
	"context"

	"github.com/gocolly/colly"
)

func (n *NamecheapService) GetCouponCode(ctx context.Context) (string, error) {
	c := colly.NewCollector()

	var couponCode string
	c.OnHTML("button", func(e *colly.HTMLElement) {
		if e.Attr("class") == "button" {
			couponCode = e.Text
		}
	})

	err := c.Visit("https://www.namecheap.com/promos/coupons/")
	if err != nil {
		return "", err
	}

	return couponCode, nil
}
