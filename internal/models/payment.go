package models

type Payment struct {
	Id                 string `json:"id" validate:"required,uuid4"`
	PaymentStatus      string `json:"payment_status" validate:"required,oneof= Authorized Declined"`
	CardNumberLastFour int    `json:"card_number_last_four" validate:"required,int_len=4"`
	ExpiryMonth        int    `json:"expiry_month" validate:"required,min=1,max=12"`
	ExpiryYear         int    `json:"expiry_year" validate:"required,int_len=4,month_year_gt"`
	Currency           string `json:"currency" validate:"required,len=3,iso4217"`
	Amount             int    `json:"amount" validate:"required"`
}

type PostPaymentRequest struct {
	CardNumber  string `json:"card_number" validate:"required,numeric,min=14,max=19"`
	ExpiryMonth int    `json:"expiry_month" validate:"required,min=1,max=12"`
	ExpiryYear  int    `json:"expiry_year" validate:"required,int_len=4,month_year_gt"`
	Currency    string `json:"currency" validate:"required,len=3,iso4217"`
	Amount      int    `json:"amount" validate:"required"`
	Cvv         int    `json:"cvv" validate:"required,int_min_len=3,int_max_len=4"`
}

type PostPaymentResponse struct {
	Id                 string `json:"id" validate:"required,uuid4"`
	PaymentStatus      string `json:"payment_status" validate:"required,oneof= Authorized Declined"`
	CardNumberLastFour int    `json:"card_number_last_four" validate:"required,int_len=4"`
	ExpiryMonth        int    `json:"expiry_month" validate:"required,min=1,max=12"`
	ExpiryYear         int    `json:"expiry_year" validate:"required,int_len=4,month_year_gt"`
	Currency           string `json:"currency" validate:"required,len=3,iso4217"`
	Amount             int    `json:"amount" validate:"required"`
}

type GetPaymentResponse struct {
	Id                 string `json:"id" validate:"required,uuid4"`
	PaymentStatus      string `json:"payment_status" validate:"required,oneof= Authorized Declined"`
	CardNumberLastFour int    `json:"card_number_last_four" validate:"required,int_len=4"`
	ExpiryMonth        int    `json:"expiry_month" validate:"required,min=1,max=12"`
	ExpiryYear         int    `json:"expiry_year" validate:"required,int_len=4,month_year_gt"`
	Currency           string `json:"currency" validate:"required,len=3,iso4217"`
	Amount             int    `json:"amount" validate:"required"`
}
