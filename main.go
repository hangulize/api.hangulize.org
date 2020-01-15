package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	router := gin.New()
	http.Handle("/", router)

	// Allow all origins in CORS.
	router.Use(cors.Default())

	// Serve openapi.yaml.
	router.StaticFile("/openapi.yaml", "openapi.yaml")

	// Route the API.
	Register(router.Group("/v2"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
