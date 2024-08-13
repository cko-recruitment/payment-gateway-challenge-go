package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cko-recruitment/payment-gateway-challenge-go/docs"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/api"
	"github.com/cko-recruitment/payment-gateway-challenge-go/internal/validator"
	paymentProcessor "github.com/cko-recruitment/payment-gateway-challenge-go/third_party/payment_processor"
)

var (
	env = os.Getenv("ENV")
)

//	@title			Payment Gateway Challenge Go
//	@description	Interview challenge for building a Payment Gateway - Go version

//	@host		localhost:8090
//	@BasePath	/

// @securityDefinitions.basic	BasicAuth
func main() {
	fmt.Printf("environment: %s\n", env)
	docs.SwaggerInfo.Version = env

	err := run()
	if err != nil {
		fmt.Printf("fatal API error: %v\n", err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// graceful shutdown
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Printf("sigterm/interrupt signal\n")
		cancel()
	}()

	defer func() {
		// recover after panic
		if x := recover(); x != nil {
			fmt.Printf("run time panic:\n%v\n", x)
			panic(x)
		}
	}()

	// Initialize the Validator
	validator.NewValidator()

	// Initialize the PaymentProcessor
	paymentProcessor.NewBankPaymentProcessor(os.Getenv("BANK_PAYMENT_PROCESSOR_BASE_URL"))
	api := api.New()
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if err := api.Run(ctx, port); err != nil {
		return err
	}

	return nil
}
