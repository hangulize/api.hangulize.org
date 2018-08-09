package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func handler(c *gin.Context) {
	pID := c.Param("pronouncer")
	word := c.Param("word")

	p, ok := pronouncers[pID]

	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
	}

	pronounced := p.Pronounce(word)

	switch c.NegotiateFormat(gin.MIMEPlain, gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"pronouncer": "furigana",
			"word":       word,
			"pronounced": pronounced,
		})

	case gin.MIMEPlain:
		fallthrough

	default:
		c.String(http.StatusOK, pronounced)
	}
}

func register(r gin.IRouter) {
	r.GET("/pronounced/:pronouncer/:word", handler)
}

func init() {
	router := gin.New()
	register(router.Group("/v2"))
	http.Handle("/", router)
}

func main() {
	appengine.Main()
}
