package payment_processor

// PaymentProcessor interface to allow code flexibility when introducing other processors like paypal etc.
type PaymentProcessor interface {
	ProcessPayment(paymentRequest ProcessPaymentRequest) (*ProcessPaymentResponse, error)
}

type ProcessPaymentRequest struct {
	CardNumber string `json:"card_number" validate:"required,numeric,min=14,max=19"`
	ExpiryDate string `json:"expiry_date" validate:"required,str_date_gt"`
	Currency   string `json:"currency" validate:"required,len=3,iso4217"`
	Amount     int    `json:"amount" validate:"required"`
	Cvv        int    `json:"cvv" validate:"required,int_min_len=3,int_max_len=4"`
}

type ProcessPaymentResponse struct {
	Authorized     bool   `json:"authorized"`
	AuthorizedCode string `json:"authorized_code"`
}
