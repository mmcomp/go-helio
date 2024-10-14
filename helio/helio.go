package helio

import (
	currencyPackage "github.com/mmcomp/go-helio/helio/currency"
	paylinkPackage "github.com/mmcomp/go-helio/helio/paylink"
)

type Helio struct {
	CurrencyCode string
	Amount       float64
	WalletId     string
	PublicApiKey string
	SecretApiKey string
	CallbackUrl  string
}

func (helio Helio) Init(
	currencyCode,
	walletId,
	publicApiKey,
	secretApiKey,
	callbackUrl string,
	amount float64) Helio {
	helio.CurrencyCode = currencyCode
	helio.WalletId = walletId
	helio.PublicApiKey = publicApiKey
	helio.SecretApiKey = secretApiKey
	helio.Amount = amount
	helio.CallbackUrl = callbackUrl
	return helio
}

func (helio Helio) ConvertAmount() (string, *currencyPackage.Currency, error) {
	currency, err := currencyPackage.Currency{}.Load(helio.CurrencyCode)
	if err != nil {
		return "", nil, err
	}
	return currency.GetAmount(helio.Amount), currency, nil
}

func (helio Helio) GeneratePayLink() (paylinkPackage.PayLink, error) {
	convertedAmount, currency, err := helio.ConvertAmount()
	if err != nil {
		return paylinkPackage.PayLink{}, err
	}
	paylink := paylinkPackage.Paylink{
		Template:        "OTHER",
		Name:            "Space Host",
		Price:           convertedAmount,
		PricingCurrency: currency.ID,
		Recipients: []struct {
			WalletId string `json:"walletId"`
			Currency string `json:"currencyId"`
		}{
			{WalletId: helio.WalletId, Currency: currency.ID},
		},
	}

	return paylink.CreatePaylink(helio.PublicApiKey, helio.SecretApiKey, helio.CallbackUrl)
}
