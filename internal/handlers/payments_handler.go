package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/validator"
	paymentProcessor "github.com/cko-recruitment/payment-gateway-challenge-go/third_party/payment_processor"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PaymentsHandler struct {
	storage          *repository.PaymentsRepository
	paymentProcessor paymentProcessor.PaymentProcessor
}

func NewPaymentsHandler(storage *repository.PaymentsRepository, paymentProcessor paymentProcessor.PaymentProcessor) *PaymentsHandler {
	return &PaymentsHandler{
		storage:          storage,
		paymentProcessor: paymentProcessor,
	}
}

// GetHandler returns an http.HandlerFunc that handles HTTP GET requests.
// It retrieves a payment record by its ID from the storage.
// The ID is expected to be part of the URL.
func (h *PaymentsHandler) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		payment := h.storage.GetPayment(id)

		if payment != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(payment); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func (ph *PaymentsHandler) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Decode post payment request body
		var paymentRequest models.PostPaymentRequest
		if err := json.NewDecoder(r.Body).Decode(&paymentRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// 2. Validate post payment request body
		if err := validatePaymentRequest(paymentRequest); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// 3. Process payment
		payment, err := processPayment(paymentRequest, ph.paymentProcessor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		// 4. Save payment
		ph.storage.AddPayment(*payment)

		// 5. Convert payment model to payment response json
		if err := json.NewEncoder(w).Encode(payment); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
	}
}

func validatePaymentRequest(pr models.PostPaymentRequest) error {
	validator := validator.GetValidator()
	return validator.ValidateStruct(&pr)
}

func processPayment(pr models.PostPaymentRequest, pp paymentProcessor.PaymentProcessor) (*models.PostPaymentResponse, error) {
	// 1. Create process payment request
	processPaymentReq := paymentProcessor.ProcessPaymentRequest{
		CardNumber: pr.CardNumber,
		ExpiryDate: time.Date(pr.ExpiryYear, time.Month(pr.ExpiryMonth), 1, 0, 0, 0, 0, time.UTC).Format("01/2006"),
		Currency:   pr.Currency,
		Amount:     pr.Amount,
		Cvv:        pr.Cvv,
	}

	// 2. Process payment
	processPaymentResp, err := pp.ProcessPayment(processPaymentReq)
	if err != nil {
		return nil, err
	}

	// 3. Convert process payment request to payment model
	cardNumberLastFour, err := strconv.Atoi(pr.CardNumber[len(pr.CardNumber)-4:])
	if err != nil {
		return nil, fmt.Errorf("error while converting cardNumber to int ")
	}

	paymentStatus := "Declined"
	if processPaymentResp.Authorized {
		paymentStatus = "Authorized"
	}

	return &models.PostPaymentResponse{
		Id:                 uuid.New().String(),
		CardNumberLastFour: cardNumberLastFour,
		PaymentStatus:      paymentStatus,
		ExpiryMonth:        pr.ExpiryMonth,
		ExpiryYear:         pr.ExpiryYear,
		Currency:           pr.Currency,
		Amount:             pr.Amount,
	}, nil

}
