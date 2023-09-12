package main

import (
	"flag"
	"fmt"
	"net/http"

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

//	@host		localhost:8080
//	@BasePath	/

// @securityDefinitions.basic	BasicAuth
func main() {
	fmt.Printf("version %s, commit %s, built at %s", version, commit, date)

	var mode string
	flag.StringVar(&mode, "mode", "debug", "Set Gin mode")
	flag.Parse()

	gin.SetMode(mode)
	docs.SwaggerInfo.Version = version

	r := gin.Default()
	r.GET("/ping", Ping)
	r.GET("/swagger/*any", gs.WrapHandler(sf.Handler))
	r.Run(":8080")
}

// PingExample godoc
// @Summary Ping example
// @Schemes
// @Description do ping

// @Produce json
// @Success 200 {object} Pong
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, Pong{Message: "pong"})
}

type Pong struct {
	Message string `json:"message"`
}
