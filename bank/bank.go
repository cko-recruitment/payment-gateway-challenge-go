// Package bank implements remote call functionality for payment approvals,
//
// For each payment request, it encapsulates (possibly mocked, in tests)
// an API call to an external authorization provider.
package bank

import(
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cko-recruitment/payment-gateway-challenge-go/common"
	"github.com/pkg/errors"
)

const applicationJson = "application/json"

var statuses = map[bool]string{true: "Authorized", false: "Unauthorized"}

// This interface defines a single call from http.Client.
// Consequently, every http.Client implements it.
type RemoteConnection interface {
	Post(url, contentType string, body io.Reader) (*http.Response, error)
}

// The primary interface of this module.
type Bank interface {
	RequestPayment(req *common.PaymentRequest) (*common.PaymentResponse, error)
}

type Options struct {
	URL string
	Client RemoteConnection
}

type requestToBank struct {
	CardNumber string `json:"card_number"`
	ExpiryDate string `json:"expiry_date"`
	Currency string `json:"currency"`
	Amount uint32 `json:"amount"`
	CVV string `json:"cvv"`
}

type responseFromBank struct {
	Authorized bool `json:"authorized"`
	AuthorizationCode string `json:"authorization_code"`
	id common.ID
}

// A concrete implementation of the interface Bank.
type bank struct {
	url string
	client RemoteConnection
	logger *log.Logger
}

func New(options *Options) (Bank, error) {
	return &bank{
		url: options.URL,
		client: options.Client,
		logger: common.MakeLogger("BANK"),
	}, nil
}

// RequestPayment performs a remote call to the payment authorizin authority.
// req   incoming payment request
// ret   (unmarshalled response from bank,amy error happened during marshalling and sending req/reading and unmarshalling response)
func (b *bank) RequestPayment(req *common.PaymentRequest) (*common.PaymentResponse, error) {
	bankResp, err := b.downstreamRequest(&requestToBank{
		CardNumber: req.CardNumber,
		ExpiryDate: fmt.Sprintf("%02d/%d", req.ExpiryMonth, req.ExpiryYear),
		Currency: req.Currency,
		Amount: req.Amount,
		CVV: req.CVV,
	})
	if err != nil {
		return nil, err
	}
	return &common.PaymentResponse{
		ID: bankResp.id,
		Status: statuses[bankResp.Authorized],
		CardNumber: req.CardNumber[len(req.CardNumber) - 4:],
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear: req.ExpiryYear,
		Currency: req.Currency,
		Amount: req.Amount,
	}, nil
}

func (b *bank) downstreamRequest(req *requestToBank) (*responseFromBank, error) {
	postBody, err := common.EncodeJSON(req)
	if err != nil {
		return nil, errors.Wrap(err, "While encoding")
	}
	// TODO: consider obuscating card number for logging purposes
	b.logger.Printf("Sending downstream request %v to %s\n", postBody, b.url)
	resp, err := b.client.Post(b.url, applicationJson, postBody)
	if err != nil {
		return nil, errors.Wrap(err, "Bank request error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Invalid response from the bank")
	}
	defer resp.Body.Close()
	var rv responseFromBank
	if err = json.NewDecoder(resp.Body).Decode(&rv); err != nil {
		return nil, errors.Wrap(err, "Malformed response from the bank")
	}
	rv.id, err = common.NextID()
	return &rv, err
}
