package currency

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Currency struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Decimals        int      `json:"decimals"`
	Order           int      `json:"order"`
	MintAddress     string   `json:"mintAddress"`
	CoinMarketCapId int      `json:"coinMarketCapId"`
	Symbol          string   `json:"symbol"`
	SymbolPrefix    string   `json:"symbolPrefix"`
	Type            string   `json:"type"`
	IconUrl         string   `json:"iconUrl"`
	Features        []string `json:"features"`
	Blockchain      struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Engine struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"engine"`
	} `json:"blockchain"`
}

func (currencyStr Currency) Load(symbol string) (*Currency, error) {
	c := http.Client{Timeout: time.Duration(100) * time.Second}
	resp, err := c.Get("https://api.hel.io/v1/currency/all")
	if err != nil {
		fmt.Printf("Error in request: %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error in body reading: %s", err)
		return nil, err
	}

	var dat []Currency

	if err = json.Unmarshal(body, &dat); err != nil {
		fmt.Printf("Error in unmarshal: %s", err)
		return nil, err
	}

	for i, v := range dat {
		if v.Symbol == symbol {
			return &dat[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (currencyStr Currency) GetAmount(amount float64) string {
	decimals := currencyStr.Decimals
	return strconv.FormatFloat(amount*math.Pow(10, float64(decimals)), 'f', -1, 64)
}
