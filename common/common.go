// Package common defines project-wide utilities, common types,
// and some testing util fnuctions.
package common

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ID = uuid.UUID

var lastID ID

// NextID generates an id for a successful authorization request to the bank
func NextID() (ID, error) {
	// Could be using the response bytes, but the plain time-based is fine, too.
	var err error
	lastID, err = uuid.NewV7()
	return lastID, err
}

// ParseID converts strings into payment ids
func ParseID(id string) (ID, error) {
	return uuid.Parse(id)
}

// PaymentRequest structure is expected to be found in the data of a valid payment request.
type PaymentRequest struct {
	CardNumber string `json:"card_number" binding:"required" validate:"min=14,max=19,numeric_only"`
	// NB: the following definition seemed more natural but one of the testing payment reqs proivided contains an invalid card number
	// CardNumber string `json:"card_number" binding:"required" validate:"min=14,max=19,credit_card,excludesrune= "`
	ExpiryMonth int `json:"expiry_month" binding:"required" validate:"min=1,max=12"`
	// To not be used after 9999 AD, time travel to 2023 or earlier is not recommended either.
	ExpiryYear int `json:"expiry_year" binding:"required" validate:"min=2024,max=9999"`
	Currency string `json:"currency" binding:"required" validate:"iso4217,known_currency"`
	// The maximum transaction value is intentionally capped ad ~43M. Don't use fast card payments for larger amounts.
	Amount uint32 `json:"amount" binding:"required" validate:"gt=0"`
	CVV string `json:"cvv" binding:"required" validate:"min=3,max=4,numeric_only"`
}

// PaymentResponse structure is sent out in response to a valid payment request.
type PaymentResponse struct {
	ID ID `json:"id"`
	Status string `json:"status"`
	CardNumber string `json:"card_number"`
	ExpiryMonth int `json:"expiry_month"`
	ExpiryYear int `json:"expiry_year"`
	Currency string `json:"currency"`
	Amount uint32 `json:"amount"`
}

// Refuse sends a 400 response to invalid requests.
func Refuse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

// EncodeJSON marshals an object and returns the result as a Reader.
func EncodeJSON(obj any) (io.Reader, error) {
	var encodingArena bytes.Buffer
	err := json.NewEncoder(&encodingArena).Encode(obj)
	return &encodingArena, err
}

// Produces package-specific logger, with a standard package-specific prefix.
func MakeLogger(prefix string) *log.Logger {
	return log.New(os.Stdout, "[" + prefix + "] ", log.Default().Flags())
}
