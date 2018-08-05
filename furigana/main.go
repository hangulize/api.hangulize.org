package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hangulize/hangulize/pronounce/furigana"
	"google.golang.org/appengine"
)

func handler(c *gin.Context) {
	word := c.Param("word")

	p := &furigana.P
	pronounced := p.Pronounce(word)

	switch c.NegotiateFormat(gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"pronouncer": "furigana",
			"word":       word,
			"pronounced": pronounced,
		})

	default:
		c.String(http.StatusOK, pronounced)
	}
}

func register(r gin.IRouter) {
	r.GET("/pronounced/furigana/:word", handler)
}

func init() {
	router := gin.New()
	register(router.Group("/v2"))
	http.Handle("/", router)
}

func main() {
	appengine.Main()
}
