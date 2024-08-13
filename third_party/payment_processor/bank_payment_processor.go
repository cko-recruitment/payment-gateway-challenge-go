package payment_processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BankPaymentProcessor struct {
	BaseURL    string
	HTTPClient *http.Client
}

// type ProcessPaymentRequest struct {
// 	CardNumber string `json:"card_number" validate:"required,numeric,min=14,max=19"`
// 	ExpiryDate string `json:"expiry_date" validate:"required,str_date_gt"`
// 	Currency   string `json:"currency" validate:"required,len=3,iso4217"`
// 	Amount     int    `json:"amount" validate:"required"`
// 	Cvv        int    `json:"cvv" validate:"required,int_min_len=3,int_max_len=4"`
// }

// type ProcessPaymentResponse struct {
// 	Authorized     bool   `json:"authorized"`
// 	AuthorizedCode string `json:"authorized_code"`
// }

var instance *BankPaymentProcessor

func NewBankPaymentProcessor(baseURL string) *BankPaymentProcessor {
	instance = &BankPaymentProcessor{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10, // Set a timeout for requests
		},
	}
	return instance
}

func GetBankPaymentProcessor() *BankPaymentProcessor {
	return instance
}

func (bp *BankPaymentProcessor) ProcessPayment(paymentRequest ProcessPaymentRequest) (*ProcessPaymentResponse, error) {
	requestPath := fmt.Sprintf("%s%s", bp.BaseURL, "/payments")

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(paymentRequest); err != nil {
		return nil, fmt.Errorf("error encoding data to JSON: %v", err)
	}

	req, err := http.NewRequest("POST", requestPath, &buf)
	if err != nil {
		return nil, fmt.Errorf("error creating POST request: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := bp.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pocessing payment endpoint failed: %s", resp.Status)
	}

	var paymentResp ProcessPaymentResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return &paymentResp, nil
}
