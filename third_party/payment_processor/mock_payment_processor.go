package payment_processor

type MockPaymentProcessor struct {
	ProcessPaymentFunc func(req ProcessPaymentRequest) (*ProcessPaymentResponse, error)
}

func (m *MockPaymentProcessor) ProcessPayment(paymentRequest ProcessPaymentRequest) (*ProcessPaymentResponse, error) {
	if m.ProcessPaymentFunc != nil {
		return m.ProcessPaymentFunc(paymentRequest)
	}
	return nil, nil
}
