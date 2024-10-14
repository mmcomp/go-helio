package paylink

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type Paylink struct {
	Template        string `json:"template"`
	Name            string `json:"name"`
	Price           string `json:"price"`
	PricingCurrency string `json:"pricingCurrency"`
	Features        struct {
	} `json:"features"`
	Recipients []struct {
		WalletId string `json:"walletId"`
		Currency string `json:"currencyId"`
	} `json:"recipients"`
}

type PaylinkResponse struct {
	Id string `json:"id"`
}

type PayLink struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func (paylink Paylink) CreatePaylink(publicApiKey, secretApiKey, webhookUrl string) (PayLink, error) {
	requestBody, err := json.Marshal(paylink)
	if err != nil {
		return PayLink{}, err
	}
	c := http.Client{Timeout: time.Duration(100) * time.Second}
	req, err := http.NewRequest("POST", "https://api.hel.io/v1/paylink/create/api-key?apiKey="+publicApiKey, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return PayLink{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secretApiKey)
	resp, err := c.Do(req)
	if err != nil {
		return PayLink{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PayLink{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return PayLink{}, errors.New(resp.Status + ":" + string(body))
	}

	var dat PaylinkResponse

	if err = json.Unmarshal(body, &dat); err != nil {
		return PayLink{}, err
	}

	err = paylink.AddPaylinkWebhook(publicApiKey, secretApiKey, dat.Id, webhookUrl)
	if err != nil {
		return PayLink{}, err
	}

	return PayLink{Id: dat.Id, Url: "https://app.hel.io/pay/" + dat.Id}, nil
}

func (paylink Paylink) AddPaylinkWebhook(publicApiKey, secretApiKey, paylinkId, webhookUrl string) error {
	c := http.Client{Timeout: time.Duration(100) * time.Second}
	req, err := http.NewRequest("POST", "https://api.hel.io/v1/webhook/paylink/transaction?apiKey="+publicApiKey, bytes.NewBuffer([]byte(`{"paylinkId":"`+paylinkId+`","targetUrl":"`+webhookUrl+`","events": ["CREATED"]}`)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secretApiKey)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status + ":" + string(body))
	}

	return nil
}
