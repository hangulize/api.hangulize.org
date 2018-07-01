package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/hangulize/api.hangulize.org/v2"
)

func init() {
	router := gin.New()
	http.Handle("/", router)

	// Allow all origins in CORS.
	router.Use(cors.Default())

	// Serve openapi.yaml.
	router.StaticFile("/openapi.yaml", "openapi.yaml")

	// Route v2 API.
	v2.Register(router.Group("/v2"))
}

func main() {
	appengine.Main()
}
