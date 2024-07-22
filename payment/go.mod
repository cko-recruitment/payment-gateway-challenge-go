module github.com/cko-recruitment/payment-gateway-challenge-go/payment

go 1.22.5

require (
	github.com/cko-recruitment/payment-gateway-challenge-go/common v0.0.0
	github.com/cko-recruitment/payment-gateway-challenge-go/storage v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.14.0
	github.com/google/uuid v1.6.0
	github.com/pkg/errors v0.9.1
)

require github.com/stretchr/testify v1.9.0

replace (
	github.com/cko-recruitment/payment-gateway-challenge-go/common => ../common
	github.com/cko-recruitment/payment-gateway-challenge-go/storage => ../storage
)
