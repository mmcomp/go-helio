package paylink

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

func (paylink Paylink) CreatePaylink(publicApiKey, secretApiKey, webhookUrl string) (string, error) {
	requestBody, err := json.Marshal(paylink)
	if err != nil {
		fmt.Printf("Error in marshalling: %s", err)
		return "", err
	}
	c := http.Client{Timeout: time.Duration(100) * time.Second}
	req, err := http.NewRequest("POST", "https://api.hel.io/v1/paylink/create/api-key?apiKey="+publicApiKey, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		fmt.Printf("Error in creating request: %s", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secretApiKey)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error in request: %s", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error in reading: %s", err)
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Error in response[%d]: %s", resp.StatusCode, body)
		return "", errors.New(resp.Status + ":" + string(body))
	}

	var dat PaylinkResponse

	if err = json.Unmarshal(body, &dat); err != nil {
		fmt.Printf("Error in unmarshal: %s", err)
		return "", err
	}

	err = paylink.AddPaylinkWebhook(publicApiKey, secretApiKey, dat.Id, webhookUrl)
	if err != nil {
		fmt.Printf("Error in adding webhook: %s", err)
		return "", err
	}

	return "https://app.hel.io/pay/" + dat.Id, nil
}

func (paylink Paylink) AddPaylinkWebhook(publicApiKey, secretApiKey, paylinkId, webhookUrl string) error {
	c := http.Client{Timeout: time.Duration(100) * time.Second}
	req, err := http.NewRequest("POST", "https://api.hel.io/v1/webhook/paylink/transaction?apiKey="+publicApiKey, bytes.NewBuffer([]byte(`{"paylinkId":"`+paylinkId+`","targetUrl":"`+webhookUrl+`","events": ["CREATED"]}`)))
	if err != nil {
		fmt.Printf("Error in creating request: %s", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secretApiKey)
	resp, err := c.Do(req)
	if err != nil {
		fmt.Printf("Error in request: %s", err)
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error in reading: %s", err)
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Error in response[%d]: %s", resp.StatusCode, body)
		return errors.New(resp.Status + ":" + string(body))
	}

	return nil
}
