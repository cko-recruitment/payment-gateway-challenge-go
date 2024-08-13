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

var instance *BankPaymentProcessor

// NewBankPaymentProcessor initalizes a bank payment processor
func NewBankPaymentProcessor(baseURL string) *BankPaymentProcessor {
	instance = &BankPaymentProcessor{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10, // Set a timeout for requests
		},
	}
	return instance
}

// GetBankPaymentProcessor returns an instance of bankPaymentProcessor
func GetBankPaymentProcessor() *BankPaymentProcessor {
	return instance
}

// ProcessPayment makes an api call to the paymentProcessor server to process the payment and returns the processor's response
func (bp *BankPaymentProcessor) ProcessPayment(paymentRequest ProcessPaymentRequest) (*ProcessPaymentResponse, error) {
	// 1. Encode payment request
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(paymentRequest); err != nil {
		return nil, fmt.Errorf("error encoding data to JSON: %v", err)
	}

	// 2. Make the API call to paymentProcessor service
	requestPath := fmt.Sprintf("%s%s", bp.BaseURL, "/payments")
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

	// 3. Decode the response body
	var paymentResp ProcessPaymentResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("error decoding response body: %v", err)
	}

	return &paymentResp, nil
}
