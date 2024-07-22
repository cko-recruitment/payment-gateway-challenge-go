package bank_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/bank"
	"github.com/cko-recruitment/payment-gateway-challenge-go/common"

	"github.com/stretchr/testify/assert"
)

const expectedContentType = "application/json"

type singleRequestCase struct {
	mockedPaymentRequest common.PaymentRequest
	expectedRequestToBank string

	mockedResponseFromBank *http.Response
	mockedError error

	expectedPaymentResponse *common.PaymentResponse
	expectedError string
}

type mockConnection struct {
	assert *assert.Assertions
	url string
	cases []singleRequestCase
	i int
}

func (mc *mockConnection) isEmpty() bool {
	return mc.i == len(mc.cases)
}

func (mc *mockConnection) nextCase() *singleRequestCase {
	return &mc.cases[mc.i]
}

func (mc *mockConnection) runCase(b bank.Bank) {
	currentCase := mc.cases[mc.i]
	message := fmt.Sprintf("case number %d", mc.i)
	resp, err := b.RequestPayment(&currentCase.mockedPaymentRequest)
	if resp != nil {
		mc.assert.NotNil(currentCase.expectedPaymentResponse, message)
		mc.assert.Equal(currentCase.expectedPaymentResponse.Status     , resp.Status     , message)
		mc.assert.Equal(currentCase.expectedPaymentResponse.CardNumber , resp.CardNumber , message)
		mc.assert.Equal(currentCase.expectedPaymentResponse.ExpiryMonth, resp.ExpiryMonth, message)
		mc.assert.Equal(currentCase.expectedPaymentResponse.ExpiryYear , resp.ExpiryYear , message)
		mc.assert.Equal(currentCase.expectedPaymentResponse.Currency   , resp.Currency   , message)
		mc.assert.Equal(currentCase.expectedPaymentResponse.Amount     , resp.Amount     , message)
	} else {
		mc.assert.Nil(currentCase.expectedPaymentResponse, message)
	}
	if currentCase.expectedError != "" {
		mc.assert.EqualError(err, currentCase.expectedError, message)
	} else {
		mc.assert.Nil(err, message)
	}
	mc.i++
}

func (mc *mockConnection) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	msg := common.TestMessage(mc.i)
	bodyContent, err := common.ReaderContent(body)
	mc.assert.Nil(err, msg)
	mc.assert.Falsef(mc.isEmpty(), "Unexpected request to %s, Content-Type = %s, body: %s", url, contentType, bodyContent)
	if mc.isEmpty() {
		return nil, fmt.Errorf("One too many requests")
	}
	nextCase := mc.nextCase()
	mc.assert.Equal(mc.url, url, msg)
	mc.assert.Equal(expectedContentType, contentType, msg)
	mc.assert.Equal(nextCase.expectedRequestToBank, bodyContent, msg)
	return nextCase.mockedResponseFromBank, nextCase.mockedError
}

func makeBody(mockedContent string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(mockedContent))
}

var cases = []singleRequestCase{
	singleRequestCase{
		mockedPaymentRequest: common.PaymentRequest{
			// it's not validated anymore so needs not to be a plausible card no
			CardNumber: "123456789",
			ExpiryMonth: 12,
			ExpiryYear: 2025,
			Currency: "GBP",
			Amount: 100,
			CVV: "123",
		},
		expectedRequestToBank: `{"card_number":"123456789","expiry_date":"12/2025","currency":"GBP","amount":100,"cvv":"123"}` + "\n",
		mockedResponseFromBank: &http.Response{
			Status: "200 OK",
			StatusCode: 200,
			Proto: "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: map[string][]string{"Content-Type": []string{expectedContentType}},
			Body: makeBody(`{"authorized":true,"authorization_code":"Hello world!"}`),
			ContentLength: 57,
		},
		mockedError: nil,
		expectedPaymentResponse: &common.PaymentResponse{
			// for simplicity, IDs are not mocked here (they could be, it would only require to 
			Status: "Authorized",
			CardNumber: "6789",
			ExpiryMonth: 12,
			ExpiryYear: 2025,
			Currency: "GBP",
			Amount: 100,
		},
		expectedError: "",
	},
	singleRequestCase{
		mockedPaymentRequest: common.PaymentRequest{
			// it's not validated anymore so needs not to be a plausible card no
			CardNumber: "987654321",
			ExpiryMonth: 12,
			ExpiryYear: 2025,
			Currency: "USD",
			Amount: 200,
			CVV: "123",
		},
		expectedRequestToBank: `{"card_number":"987654321","expiry_date":"12/2025","currency":"USD","amount":200,"cvv":"123"}` + "\n",
		mockedResponseFromBank: &http.Response{
			Status: "200 OK",
			StatusCode: 200,
			Proto: "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: map[string][]string{"Content-Type": []string{expectedContentType}},
			Body: makeBody(`{"authorized":false,"reason":"the authorizer was having lunch"}`),
			ContentLength: 71,
		},
		mockedError: nil,
		expectedPaymentResponse: &common.PaymentResponse{
			Status: "Unauthorized",
			CardNumber: "4321",
			ExpiryMonth: 12,
			ExpiryYear: 2025,
			Currency: "USD",
			Amount: 200,
		},
		expectedError: "",
	},
	singleRequestCase{
		mockedPaymentRequest: common.PaymentRequest{
			CardNumber: "Ceci n'est pas un numero de carte",
			ExpiryMonth: 12,
			ExpiryYear: 2025,
			Currency: "BEF",
			Amount: 300,
			CVV: "123",
		},
		expectedRequestToBank: `{"card_number":"Ceci n'est pas un numero de carte","expiry_date":"12/2025","currency":"BEF","amount":300,"cvv":"123"}` + "\n",
		mockedResponseFromBank: &http.Response{
			Status: "400 OK",
			StatusCode: 400,
			Proto: "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: map[string][]string{"Content-Type": []string{expectedContentType}},
			Body: makeBody(`{"error":"card number is too surrealiste"}`),
			ContentLength: 42,
		},
		mockedError: nil,
		expectedPaymentResponse: nil,
		expectedError: "Invalid response from the bank",
	},
	singleRequestCase{
		mockedPaymentRequest: common.PaymentRequest{
			CardNumber: "Au contraire, cet ici est le numero d'une carte",
			ExpiryMonth: 12,
			ExpiryYear: 2025,
			Currency: "BEF",
			Amount: 300,
			CVV: "123",
		},
		expectedRequestToBank: `{"card_number":"Au contraire, cet ici est le numero d'une carte","expiry_date":"12/2025","currency":"BEF","amount":300,"cvv":"123"}` + "\n",
		mockedResponseFromBank: &http.Response{
			Status: "400 OK",
			StatusCode: 400,
			Proto: "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: map[string][]string{"Content-Type": []string{expectedContentType}},
			Body: makeBody(`{"error":"card number is too surrealiste"}`),
			ContentLength: 42,
		},
		mockedError: fmt.Errorf("CrowdStrike is in effect in this bank."),
		expectedPaymentResponse: nil,
		expectedError: "Bank request error: CrowdStrike is in effect in this bank.",
	},
}

func TestRequestPayment(t *testing.T) {
	assert := assert.New(t)
	client := &mockConnection{assert, "http://globalhost", cases, 0}
	bank, _ := bank.New(&bank.Options{URL: "http://globalhost", Client: client})
	for !client.isEmpty() {
		client.runCase(bank)
	}
}
