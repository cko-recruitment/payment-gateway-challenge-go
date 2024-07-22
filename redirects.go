package main

import (
	"net/http"
	"strconv"

	"github.com/cko-recruitment/payment-gateway-challenge-go/common"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)


// paymentRedirect godoc
// @Summary Simple endpoint for the html/index.html example
// @Schemes
// @Description Compiles form values into JSON data suitable to post to /pay/

//	@Param	card_number	query	string	true	"Card number, digits only"
//	@Param	expiry_year	query	uint32	true	"Expiry year"
//	@Param	expiry_month	query	uint32	true	"month, from January=1 through December=12"
//	@Param	cvv		query	string	true	"CVV/CSC"
//	@Param	cyrrency	query	string	true	"ISO 4217 letter code"
//	@Param	amount		query	string	true	"amount, in currency cents"
//	@Produce	json
//	@Success	200	{object}	common.PaymentResponse
//	@Failure	400	{object}	error
//	@Router		/make_payment [post]
func paymentRedirect(eng *gin.Engine) func(*gin.Context) {
	return func(c *gin.Context) {
		value := c.PostForm("expiry_month")
		month, err := strconv.Atoi(value)
		if err != nil {
			common.Refuse(c, errors.Wrapf(err, "Invalid expiry month %s", value))
			return
		}
		value = c.PostForm("expiry_year")
		year, err := strconv.Atoi(value)
		if err != nil {
			common.Refuse(c, errors.Wrapf(err, "Invalid expiry year %s", value))
			return
		}
		value = c.PostForm("amount")
		amount, err := strconv.ParseFloat(value, 64)
		if err != nil {
			common.Refuse(c, errors.Wrapf(err, "Invalid amount %s", value))
			return
		}
		payReq := common.PaymentRequest{
			CardNumber: c.PostForm("card_number"),
			ExpiryMonth: month,
			ExpiryYear: year,
			Currency: c.PostForm("currency"),
			Amount: uint32(amount * 100),
			CVV: c.PostForm("cvv"),
		}

		postBody, err := common.EncodeJSON(payReq)
		if err != nil {
			common.Refuse(c, err)
			return
		}
		c.Request, err = http.NewRequest(http.MethodPost, "/pay", postBody)
		if err != nil {
			common.Refuse(c, err)
			return
		}

		eng.HandleContext(c)
	}
}
