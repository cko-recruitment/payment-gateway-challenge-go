package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/cko-recruitment/payment-gateway-challenge-go/bank"
	"github.com/cko-recruitment/payment-gateway-challenge-go/payment"
	"github.com/cko-recruitment/payment-gateway-challenge-go/storage"


	"github.com/cko-recruitment/payment-gateway-challenge-go/docs"
	"github.com/gin-gonic/gin"
	sf "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)



//	@title			Payment Gateway Challenge Go
//	@description	Interview challenge for building a Payment Gateway - Go version

//	@host		localhost:8090
//	@BasePath	/

//	@securityDefinitions.basic	BasicAuth
func main() {
	fmt.Printf("version %s, commit %s, built at %s", version, commit, date)

	var (
		mode string
		bankUrl string
	)
	flag.StringVar(&mode, "mode", "debug", "Set Gin mode")
	flag.StringVar(&bankUrl, "bank-url", "http://localhost:8080/payments", "Bank gateway URL")
	flag.Parse()

	gin.SetMode(mode)
	docs.SwaggerInfo.Version = version

	storage, err := storage.New(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Storage error: %s\n", err.Error())
		os.Exit(1)
	}

	bank, err := bank.New(&bank.Options{URL: bankUrl, Client: &http.Client{}})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't initialize connectivity to bank: %s\n", err.Error())
		os.Exit(2)
	}

	paymentProcessor, err := payment.New(payment.Options{Bank: bank, Storage: storage})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Payment processor error: %s\n", err.Error())
		os.Exit(3)
	}

	eng := gin.Default()
	eng.LoadHTMLFiles("html/index.html")
	eng.GET("/ping", Ping)
	eng.POST("/pay", paymentProcessor.MakePayment)
	eng.GET("/recall", storage.Recall)
	eng.GET("/index.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	eng.POST("/make_payment", paymentRedirect(eng))
	eng.GET("/swagger/*any", gs.WrapHandler(sf.Handler))

	eng.Run(":8090")
}

// PingExample godoc
//	@Summary	Ping example
//	@Schemes
//	@Description	do ping

//	@Produce	json
//	@Success	200	{object}	Pong
//	@Router		/ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, Pong{Message: "pong"})
}

type Pong struct {
	Message string `json:"message"`
}
