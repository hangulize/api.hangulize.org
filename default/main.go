package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
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
	appengine.Main()
}
