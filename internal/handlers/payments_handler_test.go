package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/validator"
	paymentProcessor "github.com/cko-recruitment/payment-gateway-challenge-go/third_party/payment_processor"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetPaymentHandler(t *testing.T) {
	payment := models.PostPaymentResponse{
		Id:                 "test-id",
		PaymentStatus:      "test-successful-status",
		CardNumberLastFour: 1234,
		ExpiryMonth:        10,
		ExpiryYear:         2035,
		Currency:           "GBP",
		Amount:             100,
	}
	ps := repository.NewPaymentsRepository()
	ps.AddPayment(payment)

	payments := NewPaymentsHandler(ps, &paymentProcessor.MockPaymentProcessor{})

	r := chi.NewRouter()
	r.Get("/api/payments/{id}", payments.GetHandler())

	httpServer := &http.Server{
		Addr:    ":8091",
		Handler: r,
	}

	go func() error {
		return httpServer.ListenAndServe()
	}()

	t.Run("PaymentFound", func(t *testing.T) {
		// Create a new HTTP request for testing
		req, _ := http.NewRequest("GET", "/api/payments/test-id", nil)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})
	t.Run("PaymentNotFound", func(t *testing.T) {
		// Create a new HTTP request for testing with a non-existing payment ID
		req, _ := http.NewRequest("GET", "/api/payments/NonExistingID", nil)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the HTTP status code in the response

		assert.Equal(t, w.Code, http.StatusNotFound)
	})
}

func TestPostPaymentHandler(t *testing.T) {
	ps := repository.NewPaymentsRepository()
	r := chi.NewRouter()
	httpServer := &http.Server{
		Addr:    ":8091",
		Handler: r,
	}

	go func() error {
		return httpServer.ListenAndServe()
	}()

	validator.NewValidator()
	payments := NewPaymentsHandler(ps, &paymentProcessor.MockPaymentProcessor{})

	t.Run("Post Payment Successfully", func(t *testing.T) {
		var buf bytes.Buffer
		paymentRequestBody := &models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 10,
			ExpiryYear:  2035,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		}
		json.NewEncoder(&buf).Encode(paymentRequestBody)

		mockPaymentProcessor := &paymentProcessor.MockPaymentProcessor{
			ProcessPaymentFunc: func(req paymentProcessor.ProcessPaymentRequest) (*paymentProcessor.ProcessPaymentResponse, error) {
				return &paymentProcessor.ProcessPaymentResponse{
					Authorized: true,
				}, nil
			},
		}

		payments = NewPaymentsHandler(ps, mockPaymentProcessor)
		r.Post("/api/payments", payments.PostHandler())
		req, _ := http.NewRequest("POST", "/api/payments", &buf)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		resp := w.Result()

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check that specific keys and values are present
		var responseBody map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}

		if _, ok := responseBody["id"]; !ok {
			t.Error("Expected response body to contain an id")
		}

		if status, ok := responseBody["payment_status"]; !ok || status != "Authorized" {
			t.Errorf("Expected payment_status 'Authorized', got '%v'", status)
		}

		if amount, ok := responseBody["amount"].(float64); !ok || amount != float64(paymentRequestBody.Amount) {
			t.Errorf("Expected amount %v, got %v", paymentRequestBody.Amount, amount)
		}

		if currency, ok := responseBody["currency"]; !ok || currency != paymentRequestBody.Currency {
			t.Errorf("Expected currency %v, got '%v'", paymentRequestBody.Currency, currency)
		}

		if cardNumberLastFour, ok := responseBody["card_number_last_four"].(float64); !ok || cardNumberLastFour != 8112 {
			t.Errorf("Expected cardNumberLastFour 8112, got '%v'", cardNumberLastFour)
		}

		if expiryMonth, ok := responseBody["expiry_month"].(float64); !ok || expiryMonth != float64(paymentRequestBody.ExpiryMonth) {
			t.Errorf("Expected expiryMonth %v, got '%v'", paymentRequestBody.ExpiryMonth, expiryMonth)
		}

		if expiryYear, ok := responseBody["expiry_year"].(float64); !ok || expiryYear != float64(paymentRequestBody.ExpiryYear) {
			t.Errorf("Expected expiryYear %v, got '%v'", paymentRequestBody.ExpiryYear, expiryYear)
		}

		// Check that the payment was saved in the storage
		id, _ := responseBody["id"].(string)
		assert.NotNil(t, ps.GetPayment(id))

	})
	t.Run("Post Payment will fail as paymentProcessor returned a 400", func(t *testing.T) {
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 10,
			ExpiryYear:  2035,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})
		mockPaymentProcessor := &paymentProcessor.MockPaymentProcessor{
			ProcessPaymentFunc: func(req paymentProcessor.ProcessPaymentRequest) (*paymentProcessor.ProcessPaymentResponse, error) {
				return nil, fmt.Errorf("400 Bad Request")
			},
		}

		payments = NewPaymentsHandler(ps, mockPaymentProcessor)
		r.Post("/api/payments", payments.PostHandler())
		req, _ := http.NewRequest("POST", "/api/payments", &buf)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	})
	t.Run("Post Payment will return 400 as CardNumber contains alphabets", func(t *testing.T) {
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "AA22405343248112",
			ExpiryMonth: 10,
			ExpiryYear:  2035,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as CardNumber length is less than 14 digits", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2224053432481",
			ExpiryMonth: 10,
			ExpiryYear:  2035,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as CardNumber length is greater than 19 digits", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "22240534324811234567",
			ExpiryMonth: 10,
			ExpiryYear:  2035,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as ExpiryMonth is less than 1", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 0,
			ExpiryYear:  2035,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as ExpiryMonth is greater than 12", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 13,
			ExpiryYear:  2035,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("Post Payment will return 400 as ExpiryYear as number of digits is greater than 4", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 1,
			ExpiryYear:  20355,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as ExpiryYear as number of digits is less than 4", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 1,
			ExpiryYear:  203,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as ExpiryYear and ExpiryMonth are not in the future", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 1,
			ExpiryYear:  2024,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as Currency length is less than 3", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 1,
			ExpiryYear:  2040,
			Currency:    "GB",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as Currency length is greater than 3", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 1,
			ExpiryYear:  2040,
			Currency:    "GBPP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as Currency is not a correct ISO currency", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 1,
			ExpiryYear:  2040,
			Currency:    "EPP",
			Amount:      100,
			Cvv:         123,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as CVV number of digits are less than 3", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 1,
			ExpiryYear:  2040,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         12,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
	t.Run("Post Payment will return 400 as CVV number of digits are greater than 4", func(t *testing.T) {
		// Create a new HTTP request for testing
		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(&models.PostPaymentRequest{
			CardNumber:  "2222405343248112",
			ExpiryMonth: 1,
			ExpiryYear:  2040,
			Currency:    "GBP",
			Amount:      100,
			Cvv:         12345,
		})

		req, _ := http.NewRequest("POST", "/api/payments", &buf)

		// Create a new HTTP request recorder for recording the response
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Check the body is not nil
		assert.NotNil(t, w.Body)

		// Check the HTTP status code in the response
		if status := w.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
}
