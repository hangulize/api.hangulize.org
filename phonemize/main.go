package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func handler(c *gin.Context) {
	pID := c.Param("phonemizer")
	word := c.Param("word")

	p, ok := phonemizers[pID]

	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
	}

	phonemized := p.Phonemize(word)

	switch c.NegotiateFormat(gin.MIMEPlain, gin.MIMEJSON) {

	case gin.MIMEJSON:
		c.JSON(http.StatusOK, gin.H{
			"phonemizer": "furigana",
			"word":       word,
			"phonemized": phonemized,
		})

	case gin.MIMEPlain:
		fallthrough

	default:
		c.String(http.StatusOK, phonemized)
	}
}

func register(r gin.IRouter) {
	r.GET("/phonemized/:phonemizer/:word", handler)
}

func init() {
	router := gin.New()
	register(router.Group("/v2"))
	http.Handle("/", router)
}

func main() {
	appengine.Main()
}
