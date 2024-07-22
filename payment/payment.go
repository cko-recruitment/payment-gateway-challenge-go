// Package payment is the primary module of this whole project.
//
// Package payment provides the handling of the /pay endpoint.
package payment

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/bank"
	"github.com/cko-recruitment/payment-gateway-challenge-go/common"
	"github.com/cko-recruitment/payment-gateway-challenge-go/storage"

	"github.com/gin-gonic/gin"
	v "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Options struct {
	bank.Bank
	storage.Storage
}

// The primary interface of this module, its sole method is the one backing the /pay endpoint.
type PaymentProcessor interface {
	MakePayment(c *gin.Context)	// pay for gin in this context
}

// A concrete implementation of the interface PaymentProcessor.
type paymentProcessor struct {
	bank bank.Bank
	storage storage.Storage
	validate *v.Validate
	client *http.Client
	logger *log.Logger
}

func currentMonth() (year int, month int) {
	// Unsure about the rules banks might have about timezones concerning card expiration.
	// Let's just assume UTC is fine.
	year, m, _ := time.Now().UTC().Date()
	month = int(m)
	return
}

func validate() (*v.Validate, error) {
	rv := v.New()

	// validator v10 provides a built-in type "numeric" but it can't be used here, as it's different from "alpha" and "alphanumeric".
	numericRegexp, err := regexp.Compile("^\\d*$")
	if err != nil {
		return nil, errors.Wrap(err, "Unexpected error at service initialization")
	}
	rv.RegisterValidation(
		"numeric_only",
		func (fl v.FieldLevel) bool {
			s := fl.Field().String()
			return numericRegexp.MatchString(s)
		})

	// This can be a config value, hardcoded here for simpliciy
	knownCurrencies := map[string]bool{"USD": true, "GBP": true, "EUR": true}
	rv.RegisterValidation(
		"known_currency",
		func (fl v.FieldLevel) bool {
			s := fl.Field().String()
			return knownCurrencies[s]
		})
	
	rv.RegisterStructValidation(
		func(sl v.StructLevel) {
			req := sl.Current().Interface().(common.PaymentRequest)
			year, month := currentMonth()
			if req.ExpiryYear < year {
				sl.ReportError(req.ExpiryYear, "expiry_year", "ExpiryYear", "expiryyear", "")
			} else if req.ExpiryYear == year && req.ExpiryMonth < int(month) {
				sl.ReportError(req.ExpiryMonth, "expiry_month/expiry_year", "ExpiryMonth/ExpiryYear", "expirydate", "")
			}
		}, common.PaymentRequest{})
	
	return rv, nil
}

// New returns a fresh PaymentProcessor.
//	options	 processor's dependencies, Bank and Storage.
//	returns  (a fresh PaymentProcessor, any error occured)
func New(options Options) (PaymentProcessor, error) {
	validate, err := validate()
	if err != nil {
		return nil, err
	}

	return &paymentProcessor{
		bank: options.Bank,
		storage: options.Storage,
		validate: validate,
		client: &http.Client{},
		logger: common.MakeLogger("PAY"),
	}, nil
}

// MakePayment godoc
//	@Summary	The /pay endpoint
//	@Schemes
//	@Description	Request payment authoirization from bank and return the result.

//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	common.PaymentResponse
//	@Failure	400	{object}	error
//	@Router		/pay [post]
func (pp *paymentProcessor) MakePayment(c *gin.Context) {
	pp.logger.Printf("Incoming payment request")
	var req common.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pp.logger.Printf("Malformed payment request: %s", err.Error())
		common.Refuse(c, err)
		return
	}
	if err := pp.validate.Struct(&req); err != nil {
		pp.logger.Printf("Invalid payment request: %s", err.Error())
		common.Refuse(c, err)
		return
	}
	resp, err := pp.bank.RequestPayment(&req)
	if err != nil {
		pp.logger.Printf("Invalid response from bank: %s", err.Error())
		common.Refuse(c, err)
		return
	}
	err = pp.storage.RecordPayment(resp)
	if err != nil {
		pp.logger.Printf("Storage shortage: %s", err.Error())
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	pp.logger.Printf("Processed payment %v for card %s, result: %s", resp.ID, resp.CardNumber, resp.Status)
	c.JSON(http.StatusOK, resp)
}
