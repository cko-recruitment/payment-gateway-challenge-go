package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/models"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/repository"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/validator"
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

	payments := NewPaymentsHandler(ps)

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
	payments := NewPaymentsHandler(ps)

	r := chi.NewRouter()
	r.Post("/api/payments", payments.PostHandler())

	httpServer := &http.Server{
		Addr:    ":8091",
		Handler: r,
	}

	go func() error {
		return httpServer.ListenAndServe()
	}()

	validator.NewValidator()

	t.Run("Post Payment will return 400 as CardNumber contains alphabets", func(t *testing.T) {
		// Create a new HTTP request for testing
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
