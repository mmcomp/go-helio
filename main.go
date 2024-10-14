package main

import (
	"fmt"

	helio "github.com/mmcomp/go-helio/helio"
)

func main() {

	currencyCode := "USDC"
	amount := 1.0
	walletId := "6601331a9bb9546495332845"
	publicApiKey := "Zoogyno_0an8Lanf67B8JScf67Hx53SKDierF71naX36e0JZ0EkMdlD6~n53wxme"
	secretApiKey := "sX3VZ5zJIaYHCF5oXGYtTZjk1qXRsfX4+bY/ADyAHvPlVwSHG9GqM4wjbKDMP/MCSnDh4ah6wXwcsdtbTYti5v42gFUI5KB/DXRcOeLewBFsYysas6hDvlGXeXPh3OG6"
	callBackUrl := "http://localhost:3000"

	helio := helio.Helio{}.Init(currencyCode, walletId, publicApiKey, secretApiKey, callBackUrl, amount)

	paylink, err := helio.GeneratePayLink()
	if err != nil {
		fmt.Printf("Error %s", err)
	}
	fmt.Println("paylinkId : ", paylink.Url)
}
