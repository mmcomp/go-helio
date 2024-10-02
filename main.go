package main

import (
	"fmt"

	currencyPackage "github.com/mmcomp/go-helio/helio/currency"
	paylinkPackage "github.com/mmcomp/go-helio/helio/paylink"
)

func main() {

	currencyCode := "USDC"
	amount := 1.0
	walletId := "6601331a9bb9546495332845"
	publicApiKey := "Zoogyno_0an8Lanf67B8JScf67Hx53SKDierF71naX36e0JZ0EkMdlD6~n53wxme"
	secretApiKey := "sX3VZ5zJIaYHCF5oXGYtTZjk1qXRsfX4+bY/ADyAHvPlVwSHG9GqM4wjbKDMP/MCSnDh4ah6wXwcsdtbTYti5v42gFUI5KB/DXRcOeLewBFsYysas6hDvlGXeXPh3OG6"

	currency, err := currencyPackage.Currency{}.Load(currencyCode)
	if err != nil {
		fmt.Printf("Error %s", err)
	}
	convertedAmount := currency.GetAmount(amount)

	paylink := paylinkPackage.Paylink{
		Template:        "OTHER",
		Name:            "Space Host",
		Price:           convertedAmount,
		PricingCurrency: currency.ID,
		Recipients: []struct {
			WalletId string `json:"walletId"`
			Currency string `json:"currencyId"`
		}{
			{WalletId: walletId, Currency: currency.ID},
		},
	}

	paylinkUrl, err := paylink.CreatePaylink(publicApiKey, secretApiKey, "http://localhost:3000")
	if err != nil {
		fmt.Printf("Error %s", err)
	}
	fmt.Println("paylinkId : ", paylinkUrl)
}
