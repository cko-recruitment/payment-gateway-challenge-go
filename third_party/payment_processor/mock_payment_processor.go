package payment_processor

type MockPaymentProcessor struct {
	ProcessPaymentFunc func(req ProcessPaymentRequest) (*ProcessPaymentResponse, error) // to be passed by tests to mock different responses
}

func (m *MockPaymentProcessor) ProcessPayment(paymentRequest ProcessPaymentRequest) (*ProcessPaymentResponse, error) {
	if m.ProcessPaymentFunc != nil {
		return m.ProcessPaymentFunc(paymentRequest)
	}
	return nil, nil
}
