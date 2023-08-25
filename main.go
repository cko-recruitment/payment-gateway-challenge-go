package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fmt.Printf("version %s, commit %s, built at %s", version, commit, date)

	var mode string
	flag.StringVar(&mode, "mode", "debug", "Set Gin mode")
	flag.Parse()

	gin.SetMode(mode)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
