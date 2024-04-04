package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Define a flag for the port
	port := flag.String("p", "8080", "port number")
	flag.Parse()

	// Set the Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize the Gin engine
	r := gin.New()

	// Define a route
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Start the server
	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
