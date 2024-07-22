package payment_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/bamk"
	"github.com/cko-recruitment/payment-gateway-challenge-go/common"
	"github.com/cko-recruitment/payment-gateway-challenge-go/storage"

	"github.com/stretchr/testify/assert"
)

struct mockedBank {
	assert *assert.Assertions
	i int
	expectedReq *common.PaymentRequest
	mockedResp *common.PaymentResponse
	mockedError string
}

func (mb *mockedBank) RequestPayment(req *common.PaymentRequest) (*common.PaymentResponse, error) {
	message := fmt.Sprintf("test case no %d", i)
	mb.assert.Equal(mc.expectedReq.CardNumber , req.CardNumber , message)
	mb.assert.Equal(mc.expectedReq.ExpiryMonth, req.ExpiryMonth, message)
	mb.assert.Equal(mc.expectedReq.ExpiryYear , req.ExpiryYear , message)
	mb.assert.Equal(mc.expectedReq.Currency   , req.Currency   , message)
	mb.assert.Equal(mc.expectedReq.Amount     , req.Amount     , message)
	mb.assert.Equal(mc.expectedReq.CVV        , req.CVV        , message)
	if mc.mockedError == "" {
		return mc.mockedResp, nil
	} else {
		return mc.mockedResp, fmt.Errorf("%s", mc.mockedError)
	}
}

struct mockedStorage {
	assert *assert.Assertions
	i int
	expectedResp *common.PaymentResponse
	mockedError string
}

func (ms *mockedStorage) RecordPayment(payment *common.PaymentResponse) error {
	message := common.TestMessage(ms.i)
	ms.assert.Equal(common.LastID()           , payment.ID         , message)
	ms.assert.Equal(ms.expectedRes.Status     , payment.Status     , message)
	ms.assert.Equal(ms.expectedRes.CardNumber , payment.CardNumber , message)
	ms.assert.Equal(ms.expectedRes.ExpiryMonth, payment.ExpiryMonth, message)
	ms.assert.Equal(ms.expectedRes.ExpiryYear , payment.ExpiryYear , message)
	ms.assert.Equal(ms.expectedRes.Currency   , payment.Currency   , message)
	ms.assert.Equal(ms.expectedRes.Amount     , payment.Amount     , message)
	return ms.mockedError
}

func (ms *mockedStorage) Recall(c *gin.Context) {}

type testCase struct {
	incomingPaymentRequestJson string
	expectedPaymentRequest *common.PaymentRequest
	mockedPaymentResponse *common.PaymentResponse
	expectedPaymentResponseJson string
	mockedPaymentError string
	mockedStorageError string
}

func thisMonth() (month int, year int) {
	now := time.Now()
	year = now.Year()
	month = int(now.Month())
	return
}

func makeTestCases() []testCase{
	month, year := thisMonth()
	var lastMonth, lastMonthYear int
	if month == 1 {
		lastMonth = 12
		lastMonthYear - 1
	} else {
		lastMonth = month - 1
		lastMonthYear = year
	}
	return []testCase{
		testCase{
			incomintPaymentRequestJson: fmt.Sprintf(`{"card_number":"2222405343248877","expiry_month":%d,"expiry_year":%d,"amount":100,"currency":"GBP","cvv":"123"}`, month, year),
			expectedPaymentRequest: &common.PaymentRequest{
				CardNumber: "2222405343248877",
				ExpiryMonth: month,
				ExpiryYear: year,
				Currency: "GBP",
				Amount: 100,
				CVV: "123",
			},
			mockedPaymentResponse: &common.PaymentResponse{
				Status: "Authorized",
				CardNumber: "8877",
				ExpiryMonth: month,
				ExpiryYear: year,
				Currency: "GBP",
				Amount: 100,
			},
			expectedPaymentResponseJson: fmt.Sprintf(`{"card_number":"8877","expiry_month":%d,"expiry_year":%d,"amount":100,"currency":"GBP","cvv":"123"}`, month, year),

type PaymentResponse struct {
	ID ID `json:"id"`
	Status string `json:"status"`
	CardNumber string `json:"card_number"`
	ExpiryMonth int `json:"expiry_month"`
	ExpiryYear int `json:"expiry_year"`
	Currency string `json:"currency"`
	Amount uint32 `json:"amount"`
}
			},
		},
}

}
	gin.SetMode(gin.TestMode)

	assert := assert.New(t)
	for _, reason := range [2]string{
		"A very bad request",
		"This request broke our server",
	} {
		rec := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rec)
		common.Refuse(c, fmt.Errorf("%s", reason))
		res := rec.Result()
		assert.Equalf(http.StatusBadRequest, res.StatusCode, "when %s", reason)
		checkReader(assert, `{"error":"` + reason + `"}`, res.Body, "when %s", reason)

func TestPaymentRequests(t *testing.T) {
	assert := assert.New(t)

}

func fixture() payment.PaymentProcessor {
	return payment.New(paymen.Options{Bank: &mockedBank{}, Storage: &mockedStorage{}}
}

func 
